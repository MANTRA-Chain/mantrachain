<!-- order: 1 -->

# Concepts

It use of bitwise operations, with each bit representing a **permission/category** for the user. So in this way we are able to create a multi permission requirement performant and cost effectively.

The guard module has 3 purposes:

1. To check if the tx sender is the admin. This is necessarly because some transactions are restricted to the admin. In some transactions handlers accross the modules we call the `CheckIsAdmin`, which returns an error if the caller is not the admin and the transactions fails with `unathorized`.
2. The second purpose is to check if a wallet can transfer tokens to other wallets. More in [operations](03_operation.md).
3. The third is to give a way for a specific wallet to be able to execute a tx restricted only to the admin. Currently such a tx is `coin factory module` -> `create denom`. We use the 100th bit ot the account privileges and required priveleges for this. In the `./scripts/init-guard.sh` we call
```
./build/mantrachaind tx guard update-required-privileges module:coinfactory:CreateDenom AAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA//////////8= authz ...
```
which sets the 100th bit to `1` along with the default privileges, so after that the admin can execute a tx with with it can set the 100th bit of the account privileges for any wallet to `1` and this gives permission to that wallet to be able to execute `CreateDenom` transaction.
