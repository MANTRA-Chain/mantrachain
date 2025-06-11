from itertools import takewhile

from .utils import DEFAULT_DENOM, find_log_event_attrs


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
