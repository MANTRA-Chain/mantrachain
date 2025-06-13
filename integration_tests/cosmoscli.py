import json
import subprocess

from pystarport.utils import build_cli_args_safe, interact

from .utils import DEFAULT_GAS_PRICE


class ChainCommand:
    def __init__(self, cmd):
        self.cmd = cmd

    def __call__(self, cmd, *args, stdin=None, stderr=subprocess.STDOUT, **kwargs):
        "execute mantrachaind"
        args = " ".join(build_cli_args_safe(cmd, *args, **kwargs))
        return interact(f"{self.cmd} {args}", input=stdin, stderr=stderr)


class CosmosCLI:
    "the apis to interact with wallet and blockchain"

    def __init__(
        self,
        data_dir,
        node_rpc,
        cmd,
    ):
        self.data_dir = data_dir
        self._genesis = json.loads(
            (self.data_dir / "config" / "genesis.json").read_text()
        )
        self.chain_id = self._genesis["chain_id"]
        self.node_rpc = node_rpc
        self.raw = ChainCommand(cmd)
        self.output = None
        self.error = None

    @classmethod
    def init(cls, moniker, data_dir, node_rpc, cmd, chain_id):
        "the node's config is already added"
        ChainCommand(cmd)(
            "init",
            moniker,
            chain_id=chain_id,
            home=data_dir,
        )
        return cls(data_dir, node_rpc, cmd)

    def validators(self):
        return json.loads(
            self.raw(
                "query", "staking", "validators", output="json", node=self.node_rpc
            )
        )["validators"]

    def status(self):
        return json.loads(self.raw("status", node=self.node_rpc))

    def balances(self, addr, height=0):
        return json.loads(
            self.raw(
                "query",
                "bank",
                "balances",
                addr,
                height=height,
                output="json",
                home=self.data_dir,
                node=self.node_rpc,
            )
        )["balances"]

    def balance(self, addr, denom="uom", height=0):
        denoms = {
            coin["denom"]: int(coin["amount"])
            for coin in self.balances(addr, height=height)
        }
        return denoms.get(denom, 0)

    def address(self, name, bech="acc", field="address"):
        output = self.raw(
            "keys",
            "show",
            name,
            f"--{field}",
            home=self.data_dir,
            keyring_backend="test",
            bech=bech,
        )
        return output.strip().decode()

    def account(self, addr):
        return json.loads(
            self.raw(
                "query", "auth", "account", addr, output="json", node=self.node_rpc
            )
        )

    def transfer(
        self,
        from_,
        to,
        coins,
        generate_only=False,
        event_query_tx=True,
        fees=None,
        **kwargs,
    ):
        kwargs.setdefault("gas_prices", DEFAULT_GAS_PRICE)
        rsp = json.loads(
            self.raw(
                "tx",
                "bank",
                "send",
                from_,
                to,
                coins,
                "-y",
                "--generate-only" if generate_only else None,
                home=self.data_dir,
                fees=fees,
                **kwargs,
            )
        )
        if rsp["code"] == 0 and event_query_tx:
            rsp = self.event_query_tx_for(rsp["txhash"])
        return rsp

    def event_query_tx_for(self, hash):
        return json.loads(
            self.raw(
                "query",
                "event-query-tx-for",
                hash,
                home=self.data_dir,
            )
        )

    def query_all_txs(self, addr):
        txs = self.raw(
            "query",
            "txs-all",
            addr,
            home=self.data_dir,
            keyring_backend="test",
            node=self.node_rpc,
        )
        return json.loads(txs)
