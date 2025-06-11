import json
import os
import signal
import subprocess
from pathlib import Path

import web3
from pystarport import cluster, ports

from .cosmoscli import CosmosCLI
from .utils import supervisorctl, wait_for_block, wait_for_port


class Mantra:
    def __init__(self, base_dir, chain_binary="mantrachaind"):
        self._w3 = None
        self.base_dir = base_dir
        self.config = json.loads((base_dir / "config.json").read_text())
        self.chain_binary = chain_binary
        self._use_websockets = False

    def copy(self):
        return Mantra(self.base_dir)

    def w3_http_endpoint(self, i=0):
        port = ports.evmrpc_port(self.base_port(i))
        return f"http://localhost:{port}"

    def w3_ws_endpoint(self, i=0):
        port = ports.evmrpc_ws_port(self.base_port(i))
        return f"ws://localhost:{port}"

    @property
    def w3(self):
        if self._w3 is None:
            self._w3 = self.node_w3(0)
        return self._w3

    def node_w3(self, i=0):
        if self._use_websockets:
            return web3.Web3(web3.providers.WebsocketProvider(self.w3_ws_endpoint(i)))
        else:
            return web3.Web3(web3.providers.HTTPProvider(self.w3_http_endpoint(i)))

    def base_port(self, i):
        return self.config["validators"][i]["base_port"]

    def node_rpc(self, i):
        return "tcp://127.0.0.1:%d" % ports.rpc_port(self.base_port(i))

    def cosmos_cli(self, i=0) -> CosmosCLI:
        return CosmosCLI(self.node_home(i), self.node_rpc(i), self.chain_binary)

    def node_home(self, i=0):
        return self.base_dir / f"node{i}"

    def use_websocket(self, use=True):
        self._w3 = None
        self._use_websockets = use

    def supervisorctl(self, *args):
        return supervisorctl(self.base_dir / "../tasks.ini", *args)


def setup_mantra(path, base_port):
    cfg = Path(__file__).parent / ("configs/default.jsonnet")
    yield from setup_custom_mantra(path, base_port, cfg)


def setup_custom_mantra(
    path,
    base_port,
    config,
    post_init=None,
    chain_binary=None,
    wait_port=True,
    relayer=cluster.Relayer.HERMES.value,
):
    cmd = [
        "pystarport",
        "init",
        "--config",
        config,
        "--data",
        path,
        "--base_port",
        str(base_port),
        "--no_remove",
    ]
    if relayer == cluster.Relayer.RLY.value:
        cmd = cmd + ["--relayer", str(relayer)]
    if chain_binary is not None:
        cmd = cmd[:1] + ["--cmd", chain_binary] + cmd[1:]
    print(*cmd)
    subprocess.run(cmd, check=True)
    if post_init is not None:
        post_init(path, base_port, config)
    proc = subprocess.Popen(
        ["pystarport", "start", "--data", path, "--quiet"],
        preexec_fn=os.setsid,
    )
    try:
        if wait_port:
            wait_for_port(ports.rpc_port(base_port))
        c = Mantra(path / "mantra_5887-1", chain_binary=chain_binary or "mantrachaind")
        wait_for_block(c.cosmos_cli(), 1)
        yield c
    finally:
        os.killpg(os.getpgid(proc.pid), signal.SIGTERM)
        # proc.terminate()
        proc.wait()
