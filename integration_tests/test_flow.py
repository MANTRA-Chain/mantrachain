import os
import shutil
from itertools import takewhile

import web3
from eth_account import Account
from pystarport import ports

from .network import CosmosCLI
from .utils import (
    ADDRS,
    CHAIN_ID,
    DEFAULT_DENOM,
    DEFAULT_FEE,
    WEI_PER_UOM,
    assert_balance,
    eth_to_bech32,
    find_log_event_attrs,
    send_transaction,
)


def connect_cli(tmp_path, rpc, chain_binary="mantrachaind"):
    if tmp_path.exists():
        shutil.rmtree(tmp_path)
    cli = CosmosCLI(tmp_path, rpc, chain_binary)
    return cli


def connect_w3(rpc):
    w3 = web3.Web3(web3.providers.HTTPProvider(rpc))
    assert w3.eth.chain_id == 5887
    return w3


def get_fee(events):
    attrs = find_log_event_attrs(events, "tx", lambda attrs: "fee" in attrs)
    return int("".join(takewhile(lambda s: s.isdigit() or s == ".", attrs["fee"])))


def fund_recover(rpc, evm_rpc, tmp_path):
    """
    transfer fund from community to recover cosmos addr
    """
    community = "community"
    addr_community = eth_to_bech32(ADDRS[community])
    cli = connect_cli(tmp_path, rpc)
    w3 = connect_w3(evm_rpc)
    assert (
        cli.create_account(
            community,
            mnemonic=os.getenv("COMMUNITY_MNEMONIC"),
            home=tmp_path,
        )["address"]
        == addr_community
    )
    balance_community = assert_balance(cli, w3, addr_community)
    addr_recover = "mantra1h5tsd8wjefus259xmff367ltg0rpf54a9ktpza"
    balance_recover = assert_balance(cli, w3, addr_recover)
    amt = 4000
    if balance_recover >= amt:
        return
    tx = cli.transfer(
        addr_community,
        addr_recover,
        f"{amt}{DEFAULT_DENOM}",
        generate_only=True,
        chain_id=CHAIN_ID,
    )
    tx_json = cli.sign_tx_json(
        tx, addr_community, home=tmp_path, node=rpc, chain_id=CHAIN_ID
    )
    rsp = cli.broadcast_tx_json(tx_json, home=tmp_path)
    assert rsp["code"] == 0, rsp["raw_log"]
    fee = get_fee(rsp["events"])
    assert fee == DEFAULT_FEE
    assert assert_balance(cli, w3, addr_community) == balance_community - amt - fee
    assert assert_balance(cli, w3, addr_recover) == balance_recover + amt


def run_flow(rpc, evm_rpc, tmp_path):
    community = "community"
    recover = "recover"
    amt = 4000
    addr_recover = "mantra1h5tsd8wjefus259xmff367ltg0rpf54a9ktpza"

    # recover cosmos addr outside from node
    cli = connect_cli(tmp_path, rpc)
    assert (
        cli.create_account(
            recover,
            mnemonic=os.getenv("RECOVER_MNEMONIC"),
            coin_type=118,
            key_type="secp256k1",
            home=tmp_path,
        )["address"]
        == addr_recover
    )
    w3 = connect_w3(evm_rpc)
    assert assert_balance(cli, w3, recover) >= amt

    # fund test1 from all recover's balance via cosmos tx
    acc_test1 = Account.from_mnemonic(os.getenv("TESTER1_MNEMONIC"))
    addr_test1 = eth_to_bech32(acc_test1.address)
    amt2 = amt - DEFAULT_FEE
    tx = cli.transfer(
        addr_recover,
        addr_test1,
        f"{amt2}{DEFAULT_DENOM}",
        generate_only=True,
        chain_id=CHAIN_ID,
    )
    tx_json = cli.sign_tx_json(
        tx, addr_recover, home=tmp_path, node=rpc, chain_id=CHAIN_ID
    )
    rsp = cli.broadcast_tx_json(tx_json, home=tmp_path)
    assert rsp["code"] == 0, rsp["raw_log"]
    fee = get_fee(rsp["events"])
    assert fee == DEFAULT_FEE
    assert assert_balance(cli, w3, addr_test1) == amt2
    assert assert_balance(cli, w3, addr_recover) == amt2 - fee

    # send 1 wei from test1 to test2 for tolerance check
    acc_test2 = Account.from_mnemonic(os.getenv("TESTER2_MNEMONIC"))
    addr_test2 = eth_to_bech32(acc_test2.address)
    value = 1
    gas_price = 11250000000
    gas = 21000
    evm_tx = {
        "to": acc_test2.address,
        "value": value,
        "gas": gas,
        "gasPrice": gas_price,
        "nonce": w3.eth.get_transaction_count(acc_test1.address),
    }
    receipt = send_transaction(w3, evm_tx, acc_test1.key)
    assert receipt.status == 1
    evm_fee = receipt.gasUsed * receipt.effectiveGasPrice
    assert w3.eth.get_balance(acc_test2.address) == value
    assert assert_balance(cli, w3, addr_test2) == value // WEI_PER_UOM
    amt2 = (amt2 * WEI_PER_UOM - evm_fee) // WEI_PER_UOM
    assert assert_balance(cli, w3, addr_test1) == amt2

    # send 10^12 wei from test1 to test2 for tolerance check
    value = 10**12
    evm_tx["value"] = value
    evm_tx["nonce"] = w3.eth.get_transaction_count(acc_test1.address)
    receipt = send_transaction(w3, evm_tx, acc_test1.key)
    assert receipt.status == 1
    evm_fee = receipt.gasUsed * receipt.effectiveGasPrice
    assert w3.eth.get_balance(acc_test2.address) == value + 1
    assert assert_balance(cli, w3, addr_test2) == value // WEI_PER_UOM
    amt2 = (amt2 * WEI_PER_UOM - evm_fee) // WEI_PER_UOM
    assert assert_balance(cli, w3, addr_test1) == amt2

    # recycle test1's balance back to community
    balance_community_evm = w3.eth.get_balance(ADDRS[community])
    value = w3.eth.get_balance(acc_test1.address) - gas * gas_price
    evm_tx["to"] = ADDRS[community]
    evm_tx["value"] = value
    evm_tx["nonce"] = w3.eth.get_transaction_count(acc_test1.address)
    receipt = send_transaction(w3, evm_tx, acc_test1.key)
    assert receipt.status == 1
    assert w3.eth.get_balance(acc_test1.address) == 0
    assert w3.eth.get_balance(ADDRS[community]) == balance_community_evm + value
    assert assert_balance(cli, w3, eth_to_bech32(ADDRS[community])) > 0


def test_flow(mantra, tmp_path):
    rpc = os.getenv("RPC")
    evm_rpc = os.getenv("EVM_RPC")
    if not rpc or not evm_rpc:
        port = mantra.base_port(0)
        rpc = f"http://127.0.0.1:{ports.rpc_port(port)}"
        evm_rpc = f"http://127.0.0.1:{ports.evmrpc_port(port)}"
        fund_recover(rpc, evm_rpc, tmp_path)
    run_flow(rpc, evm_rpc, tmp_path)
