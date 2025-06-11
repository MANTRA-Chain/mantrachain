from itertools import takewhile
from concurrent.futures import ThreadPoolExecutor, as_completed
from eth_bloom import BloomFilter
from eth_utils import abi, big_endian_to_int
from hexbytes import HexBytes
import web3
import pytest

from .utils import (
    ADDRS,
    CONTRACTS,
    DEFAULT_DENOM,
    KEYS,
    Greeter,
    RevertTestContract,
    deploy_contract,
    find_log_event_attrs,
    send_transaction,
)


def test_simple(mantra):
    """
    check number of validators
    """
    cli = mantra.cosmos_cli()
    assert len(cli.validators()) == 2
    # check vesting account
    addr = cli.address("reserve")
    account = cli.account(addr)["account"]
    assert account["type"] == "/cosmos.vesting.v1beta1.DelayedVestingAccount"
    assert account["value"]["base_vesting_account"]["original_vesting"] == [
        {"denom": DEFAULT_DENOM, "amount": "10000000000000000000000"}
    ]


def test_transfer(mantra):
    """
    check simple transfer tx success
    """
    cli = mantra.cosmos_cli()
    addr_a = cli.address("community")
    addr_b = cli.address("reserve")
    balance_a = cli.balance(addr_a)
    balance_b = cli.balance(addr_b)
    amt = 1
    rsp = cli.transfer(addr_a, addr_b, f"{amt}{DEFAULT_DENOM}")
    assert rsp["code"] == 0, rsp["raw_log"]
    res = find_log_event_attrs(rsp["events"], "tx", lambda attrs: "fee" in attrs)
    fee = int("".join(takewhile(lambda s: s.isdigit() or s == ".", res["fee"])))
    assert cli.balance(addr_a) == balance_a - amt - fee
    assert cli.balance(addr_b) == balance_b + amt


def test_basic(mantra):
    assert mantra.w3.eth.chain_id == 5887


def test_send_transaction(mantra):
    w3 = mantra.w3
    txhash = w3.eth.send_transaction(
        {
            "from": ADDRS["validator"],
            "to": ADDRS["community"],
            "value": 1000,
        }
    )
    receipt = w3.eth.wait_for_transaction_receipt(txhash)
    assert receipt.status == 1
    assert receipt.gasUsed == 21000


def test_events(mantra):
    w3 = mantra.w3
    erc20 = deploy_contract(
        w3,
        CONTRACTS["TestERC20A"],
        key=KEYS["validator"],
        exp_gas_used=619754,
    )
    tx = erc20.functions.transfer(ADDRS["community"], 10).build_transaction(
        {"from": ADDRS["validator"]}
    )
    txreceipt = send_transaction(w3, tx, KEYS["validator"])
    assert len(txreceipt.logs) == 1
    data = "0x000000000000000000000000000000000000000000000000000000000000000a"
    expect_log = {
        "address": erc20.address,
        "topics": [
            HexBytes(
                abi.event_signature_to_log_topic("Transfer(address,address,uint256)")
            ),
            HexBytes(b"\x00" * 12 + HexBytes(ADDRS["validator"])),
            HexBytes(b"\x00" * 12 + HexBytes(ADDRS["community"])),
        ],
        "data": HexBytes(data),
        "transactionIndex": 0,
        "logIndex": 0,
        "removed": False,
    }
    assert expect_log.items() <= txreceipt.logs[0].items()

    # check block bloom
    bloom = BloomFilter(
        big_endian_to_int(w3.eth.get_block(txreceipt.blockNumber).logsBloom)
    )
    assert HexBytes(erc20.address) in bloom
    for topic in expect_log["topics"]:
        assert topic in bloom


def test_minimal_gas_price(mantra):
    w3 = mantra.w3
    gas_price = w3.eth.gas_price
    tx = {
        "to": "0x0000000000000000000000000000000000000000",
        "value": 10000,
    }
    with pytest.raises(ValueError):
        send_transaction(
            w3,
            {**tx, "gasPrice": 1},
            KEYS["community"],
        )
    receipt = send_transaction(
        w3,
        {**tx, "gasPrice": gas_price},
        KEYS["validator"],
    )
    assert receipt.status == 1


def test_transaction(mantra):
    w3 = mantra.w3
    gas_price = w3.eth.gas_price

    # send transaction
    txhash_1 = send_transaction(
        w3,
        {"to": ADDRS["community"], "value": 10000, "gasPrice": gas_price},
        KEYS["validator"],
    )["transactionHash"]
    tx1 = w3.eth.get_transaction(txhash_1)
    assert tx1["transactionIndex"] == 0

    initial_block_number = w3.eth.get_block_number()

    # tx already in mempool
    with pytest.raises(ValueError) as exc:
        send_transaction(
            w3,
            {
                "to": ADDRS["community"],
                "value": 10000,
                "gasPrice": gas_price,
                "nonce": w3.eth.get_transaction_count(ADDRS["validator"]) - 1,
            },
            KEYS["validator"],
        )
    assert "tx already in mempool" in str(exc)

    # invalid sequence
    with pytest.raises(ValueError) as exc:
        send_transaction(
            w3,
            {
                "to": ADDRS["community"],
                "value": 10000,
                "gasPrice": w3.eth.gas_price,
                "nonce": w3.eth.get_transaction_count(ADDRS["validator"]) + 1,
            },
            KEYS["validator"],
        )
    assert "invalid sequence" in str(exc)

    # out of gas
    with pytest.raises(ValueError) as exc:
        send_transaction(
            w3,
            {
                "to": ADDRS["community"],
                "value": 10000,
                "gasPrice": w3.eth.gas_price,
                "gas": 1,
            },
            KEYS["validator"],
        )["transactionHash"]
    assert "out of gas" in str(exc)

    # insufficient fee
    with pytest.raises(ValueError) as exc:
        send_transaction(
            w3,
            {
                "to": ADDRS["community"],
                "value": 10000,
                "gasPrice": 1,
            },
            KEYS["validator"],
        )["transactionHash"]
    assert "insufficient fee" in str(exc)

    # check all failed transactions are not included in blockchain
    assert w3.eth.get_block_number() == initial_block_number

    # Deploy multiple contracts
    contracts = {
        "test_revert_1": RevertTestContract(
            CONTRACTS["TestRevert"],
            KEYS["validator"],
        ),
        "test_revert_2": RevertTestContract(
            CONTRACTS["TestRevert"],
            KEYS["community"],
        ),
        "greeter_1": Greeter(
            CONTRACTS["Greeter"],
            KEYS["signer1"],
        ),
        "greeter_2": Greeter(
            CONTRACTS["Greeter"],
            KEYS["signer2"],
        ),
    }

    with ThreadPoolExecutor(4) as executor:
        future_to_contract = {
            executor.submit(contract.deploy, w3): name
            for name, contract in contracts.items()
        }

        assert_receipt_transaction_and_block(w3, future_to_contract)

    # Do Multiple contract calls
    with ThreadPoolExecutor(4) as executor:
        futures = []
        futures.append(
            executor.submit(contracts["test_revert_1"].transfer, 5 * (10**18) - 1)
        )
        futures.append(
            executor.submit(contracts["test_revert_2"].transfer, 5 * (10**18))
        )
        futures.append(executor.submit(contracts["greeter_1"].transfer, "hello"))
        futures.append(executor.submit(contracts["greeter_2"].transfer, "world"))

        assert_receipt_transaction_and_block(w3, futures)

        # revert transaction
        assert futures[0].result()["status"] == 0
        # normal transaction
        assert futures[1].result()["status"] == 1
        # normal transaction
        assert futures[2].result()["status"] == 1
        # normal transaction
        assert futures[3].result()["status"] == 1


def assert_receipt_transaction_and_block(w3, futures):
    receipts = []
    for future in as_completed(futures):
        data = future.result()
        receipts.append(data)
    assert len(receipts) == 4

    block_number = w3.eth.get_block_number()
    tx_indexes = [0, 1, 2, 3]
    for receipt in receipts:
        assert receipt["blockNumber"] == block_number
        transaction_index = receipt["transactionIndex"]
        assert transaction_index in tx_indexes
        tx_indexes.remove(transaction_index)

    block = w3.eth.get_block(block_number)
    transactions = [
        w3.eth.get_transaction_by_block(block_number, receipt["transactionIndex"])
        for receipt in receipts
    ]
    assert len(transactions) == 4
    for i, transaction in enumerate(transactions):
        assert transaction["blockNumber"] == block_number
        assert transaction["transactionIndex"] == receipts[i]["transactionIndex"]
        assert transaction["hash"] == receipts[i]["transactionHash"]
        assert transaction["hash"] in block["transactions"]
        assert transaction["blockNumber"] == block["number"]


def test_exception(mantra):
    w3 = mantra.w3
    contract = deploy_contract(
        w3,
        CONTRACTS["TestRevert"],
    )
    with pytest.raises(web3.exceptions.ContractLogicError):
        send_transaction(
            w3, contract.functions.transfer(5 * (10**18) - 1).build_transaction()
        )
    assert 0 == contract.caller.query()

    receipt = send_transaction(
        w3, contract.functions.transfer(5 * (10**18)).build_transaction()
    )
    assert receipt.status == 1, "should be successfully"
    assert 5 * (10**18) == contract.caller.query()
