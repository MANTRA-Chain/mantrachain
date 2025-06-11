import json
import os
import socket
import subprocess
import sys
import time
from pathlib import Path

from dotenv import load_dotenv
from eth_account import Account
from web3._utils.transactions import fill_nonce, fill_transaction_defaults

load_dotenv(Path(__file__).parent.parent / "scripts/.env")
load_dotenv(Path(__file__).parent.parent / "scripts/.env")
Account.enable_unaudited_hdwallet_features()
ACCOUNTS = {
    "validator": Account.from_mnemonic(os.getenv("VALIDATOR1_MNEMONIC")),
    "validator2": Account.from_mnemonic(os.getenv("VALIDATOR2_MNEMONIC")),
    "community": Account.from_mnemonic(os.getenv("COMMUNITY_MNEMONIC")),
    "signer1": Account.from_mnemonic(os.getenv("SIGNER1_MNEMONIC")),
    "signer2": Account.from_mnemonic(os.getenv("SIGNER2_MNEMONIC")),
}
KEYS = {name: account.key for name, account in ACCOUNTS.items()}
ADDRS = {name: account.address for name, account in ACCOUNTS.items()}

DEFAULT_DENOM = "uom"
# the default initial base fee used by integration tests
DEFAULT_GAS_PRICE = f"100000000000{DEFAULT_DENOM}"

TEST_CONTRACTS = {
    "TestERC20A": "TestERC20A.sol",
    "TestRevert": "TestRevert.sol",
    "Greeter": "Greeter.sol",
}


def contract_path(name, filename):
    return (
        Path(__file__).parent
        / "contracts/artifacts/contracts"
        / filename
        / (name + ".json")
    )


CONTRACTS = {
    **{
        name: contract_path(name, filename) for name, filename in TEST_CONTRACTS.items()
    },
}


def wait_for_fn(name, fn, *, timeout=240, interval=1):
    for i in range(int(timeout / interval)):
        result = fn()
        if result:
            return result
        time.sleep(interval)
    else:
        raise TimeoutError(f"wait for {name} timeout")


def get_sync_info(s):
    return s.get("SyncInfo") or s.get("sync_info")


def wait_for_block(cli, height, timeout=240):
    for i in range(timeout * 2):
        try:
            status = cli.status()
        except AssertionError as e:
            print(f"get sync status failed: {e}", file=sys.stderr)
        else:
            current_height = int(get_sync_info(status)["latest_block_height"])
            print("current block height", current_height)
            if current_height >= height:
                break
        time.sleep(0.5)
    else:
        raise TimeoutError(f"wait for block {height} timeout")


def wait_for_port(port, host="127.0.0.1", timeout=40.0):
    print("wait for port", port, "to be available")
    start_time = time.perf_counter()
    while True:
        try:
            with socket.create_connection((host, port), timeout=timeout):
                break
        except OSError as ex:
            time.sleep(0.1)
            if time.perf_counter() - start_time >= timeout:
                raise TimeoutError(
                    "Waited too long for the port {} on host {} to start accepting "
                    "connections.".format(port, host)
                ) from ex


def supervisorctl(inipath, *args):
    return subprocess.check_output(
        (sys.executable, "-msupervisor.supervisorctl", "-c", inipath, *args),
    ).decode()


def find_log_event_attrs(events, ev_type, cond=None):
    for ev in events:
        if ev["type"] == ev_type:
            attrs = {attr["key"]: attr["value"] for attr in ev["attributes"]}
            if cond is None or cond(attrs):
                return attrs
    return None


def sign_transaction(w3, tx, key=KEYS["validator"]):
    "fill default fields and sign"
    acct = Account.from_key(key)
    tx["from"] = acct.address
    tx = fill_transaction_defaults(w3, tx)
    tx = fill_nonce(w3, tx)
    return acct.sign_transaction(tx)


def send_transaction(w3, tx, key=KEYS["validator"]):
    signed = sign_transaction(w3, tx, key)
    txhash = w3.eth.send_raw_transaction(signed.rawTransaction)
    return w3.eth.wait_for_transaction_receipt(txhash)


def deploy_contract(w3, jsonfile, args=(), key=KEYS["validator"], exp_gas_used=None):
    """
    deploy contract and return the deployed contract instance
    """
    acct = Account.from_key(key)
    info = json.loads(jsonfile.read_text())
    bytecode = ""
    if "bytecode" in info:
        bytecode = info["bytecode"]
    if "byte" in info:
        bytecode = info["byte"]
    contract = w3.eth.contract(abi=info["abi"], bytecode=bytecode)
    tx = contract.constructor(*args).build_transaction({"from": acct.address})
    txreceipt = send_transaction(w3, tx, key)
    assert txreceipt.status == 1
    if exp_gas_used is not None:
        assert (
            exp_gas_used == txreceipt.gasUsed
        ), f"exp {exp_gas_used}, got {txreceipt.gasUsed}"
    address = txreceipt.contractAddress
    return w3.eth.contract(address=address, abi=info["abi"])


class Contract:
    def __init__(self, contract_path, private_key=KEYS["validator"], chain_id=5887):
        self.chain_id = chain_id
        self.account = Account.from_key(private_key)
        self.address = self.account.address
        self.private_key = private_key
        with open(contract_path) as f:
            json_data = f.read()
            contract_json = json.loads(json_data)
        self.bytecode = contract_json["bytecode"]
        self.abi = contract_json["abi"]
        self.contract = None
        self.w3 = None

    def deploy(self, w3):
        "Deploy contract on `w3` and return the receipt."
        if self.contract is None:
            self.w3 = w3
            contract = self.w3.eth.contract(abi=self.abi, bytecode=self.bytecode)
            transaction = contract.constructor().build_transaction(
                {"chainId": self.chain_id, "from": self.address}
            )
            receipt = send_transaction(self.w3, transaction, self.private_key)
            self.contract = self.w3.eth.contract(
                address=receipt.contractAddress, abi=self.abi
            )
            return receipt
        else:
            return receipt

class Greeter(Contract):
    "Greeter contract."

    def transfer(self, string):
        "Call contract on `w3` and return the receipt."
        transaction = self.contract.functions.setGreeting(string).build_transaction(
            {
                "chainId": self.chain_id,
                "from": self.address,
            }
        )
        receipt = send_transaction(self.w3, transaction, self.private_key)
        assert string == self.contract.functions.greet().call()
        return receipt

class RevertTestContract(Contract):
    "RevertTestContract contract."

    def transfer(self, value):
        "Call contract on `w3` and return the receipt."
        transaction = self.contract.functions.transfer(value).build_transaction(
            {
                "chainId": self.chain_id,
                "from": self.address,
                "gas": 100000,  # skip estimateGas error
            }
        )
        receipt = send_transaction(self.w3, transaction, self.private_key)
        return receipt