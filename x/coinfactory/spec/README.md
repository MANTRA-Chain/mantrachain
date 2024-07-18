<!-- order: 0 title: Coin Factory Overview parent: title: "coinfactory" -->

# x/coinfactory

The coinfactory module allows any account to create a new token with
the name `factory/{creator address}/{subdenom}`. Because tokens are
namespaced by creator address, this allows token minting to be
permissionless, due to not needing to resolve name collisions. A single
account can create multiple denoms, by providing a unique subdenom for each
created denom. Once a denom is created, the original creator is given
"admin" privileges over the asset. This allows them to:

- Mint their denom to any account
- Burn their denom from any account
- Create a transfer of their denom between any two accounts
- Change the admin. In the future, more admin capabilities may be added. Admins
  can choose to share admin privileges with other accounts using the authz
  module. The `ChangeAdmin` functionality, allows changing the master admin
  account, or even setting it to `""`, meaning no account has admin privileges
  of the asset.
- The `chain admin` can force transfer any denom(coinfactory denoms and any other native coin denoms), even if the admin is not the original creator. This is useful for resolving disputes.

## Contents

1. **[Transaction flows](01_txs_flows.md)**
2. **[Concepts](02_concepts.md)**
3. **[Examples](03_examples.md)**
