# Guard Module

## Table of Contents

- **[Background](#background)**
- **[Concepts](#concepts)**
- **[Installation](#installation)**
- **[State](#state)**
- **[Messages](#messages)**
- **[Queries](#queries)**
- **[Events](#events)**
- **[Params](#params)**
- **[Dependencies](#dependencies)**
- **[Example Usage](#example-usage)**
  - **[Client](#client)**
    - **[SDK](#sdk)**

## Background

The Guard module is responsible for managing all the permissions and privileges of the users in the network. It allows other modules to define their own permissions and use the Guard module to enforce them.

The Guard module provides a set of functions for updating account privileges, updating required privileges, updating locked status, updating guard transfer coins and updating authz grant/revoke execute permissions.

The Guard module is responsible for:

1. Guarding against unauthorized actions.
2. Handling all the privileges and permissions.
3. Managing the permission system.

## Concepts

The Guard module defines several concepts:

- **Accounts privileges:** Actions that can be performed within the system. Each privilege has an associated bytes array and an account that have that privilege.
- **Required privileges:** The privileges required to perform an action. These can be set by the admin or by the system operators.
- **Locked coins:** Coins that have been locked and cannot perform any actions on them.
- **Guarded transfer of coins:** A setting that can be enabled or disabled which prevents coins from being transferred without proper authorization.

## Installation

To install the Guard module, you can use the following command:

```bash
go get github.com/MANTRA-Finance/mantrachain/x/guard
```

## State

The Guard module stores the following data:

- **Accounts privileges:** A list of all the privileges of associated accounts.
- **Required privileges:** A list of the required privileges for coins transfers.
- **Locked accounts:** A list of all the locked coins.
- **Guarded transfer of coins:** A flag that can turn enable/disable the coins tranfers.

## Messages

### MsgUpdateAccountPrivileges

This message is used to update the privileges associated with a specific account. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `account`: The address of the account whose privileges will be updated.
- `privileges`: The new privileges that will be associated with the account.

### MsgUpdateAccountPrivilegesBatch

This message is used to update the privileges associated with a batch of accounts. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `accounts_privileges`: The list of accounts and their new privileges.

### MsgUpdateAccountPrivilegesGroupedBatch

This message is used to update the privileges associated with a grouped batch of accounts. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `accounts_privileges_grouped`: The list of grouped accounts and their new privileges.

### MsgUpdateGuardTransferCoins

This message is used to enable or disable the guard transfer coins. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `enabled`: Whether or not to enable the guard transfer coins.

### MsgUpdateRequiredPrivileges

This message is used to update the required privileges associated with a specific index. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `index`: The index of the required privileges to update.
- `privileges`: The new privileges that will be associated with the index.
- `kind`: The type of the required privileges.

### MsgUpdateRequiredPrivilegesBatch

This message is used to update the required privileges associated with a batch of indexes. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `required_privileges_list`: The list of indexes and their new privileges.
- `kind`: The type of the required privileges.

### MsgUpdateRequiredPrivilegesGroupedBatch

This message is used to update the required privileges associated with a grouped batch of indexes. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `required_privileges_list_grouped`: The list of grouped indexes and their new privileges.
- `kind`: The type of the required privileges.

### MsgUpdateLocked

This message is used to update the locked status associated with a specific index. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `index`: The index of the locked status to update.
- `locked`: Whether or not to lock the index.
- `kind`: The type of the locked status.

### MsgUpdateAuthzGenericGrantRevokeBatch

This message is used to grant or revoke permissions associated with a specific grantee. It takes the following parameters:

- `creator`: The address of the user who is creating the message.
- `grantee`: The address of the grantee whose permissions will be updated.
- `authz_grant_revoke_msgs_types`: The list of permissions to grant or revoke.

## Queries

### MantraChain Guard v1 Queries

This documentation provides information about the gRPC querier service and message types used in the MantraChain Guard v1 package.

#### Description

Defines the gRPC querier service for MantraChain Guard v1.

#### Methods

##### Params

Queries the parameters of the Guard module.

- Request Type: `QueryParamsRequest`
- Response Type: `QueryParamsResponse`

HTTP GET Path:
`/mantrachain/guard/v1/params`

##### AccountPrivileges

Queries an AccountPrivileges item by account.

- Request Type: `QueryGetAccountPrivilegesRequest`
- Response Type: `QueryGetAccountPrivilegesResponse`

HTTP GET Path:
`/mantrachain/guard/v1/account_privileges/{account}`

##### AccountPrivilegesAll

Queries a list of AccountPrivileges items.

- Request Type: `QueryAllAccountPrivilegesRequest`
- Response Type: `QueryAllAccountPrivilegesResponse`

HTTP GET Path:
`/mantrachain/guard/v1/account_privileges`

##### GuardTransferCoins

Queries a GuardTransferCoins item.

- Request Type: `QueryGetGuardTransferCoinsRequest`
- Response Type: `QueryGetGuardTransferCoinsResponse`

HTTP GET Path
`/mantrachain/guard/v1/guard_transfer_coins`

##### RequiredPrivileges

Queries a RequiredPrivileges item by index.

- Request Type: `QueryGetRequiredPrivilegesRequest`
- Response Type: `QueryGetRequiredPrivilegesResponse`

HTTP GET Path
`/mantrachain/guard/v1/required_privileges/{index}`

##### RequiredPrivilegesAll

Queries a list of RequiredPrivileges items.

- Request Type: `QueryAllRequiredPrivilegesRequest`
- Response Type: `QueryAllRequiredPrivilegesResponse`

HTTP GET Path
`/mantrachain/guard/v1/required_privileges`

##### Locked

Queries a Locked item by index.

- Request Type: `QueryGetLockedRequest`
- Response Type: `QueryGetLockedResponse`

HTTP GET Path
`/mantrachain/guard/v1/locked/{index}`

##### LockedAll

Queries a list of Locked items.

- Request Type: `QueryAllLockedRequest`
- Response Type: `QueryAllLockedResponse`

HTTP GET Path
`/mantrachain/guard/v1/locked`

#### Message Types

##### QueryParamsRequest

Request type for the Query/Params RPC method.

##### QueryParamsResponse

Response type for the Query/Params RPC method.

Fields:

- `params`: holds all the parameters of this module.

##### QueryGetAccountPrivilegesRequest

Request type for the AccountPrivileges RPC method.

Fields:

- `account`: account for which the AccountPrivileges is to be queried.

##### QueryGetAccountPrivilegesResponse

Response type for the AccountPrivileges RPC method.

Fields:

- `account`: account for which the AccountPrivileges was queried.
- `privileges:` privileges associated with the account.

##### QueryAllAccountPrivilegesRequest

Request type for the AccountPrivilegesAll RPC method.

Fields:

- `pagination`: page request parameters.

##### QueryAllAccountPrivilegesResponse

Response type for the AccountPrivilegesAll RPC method.

Fields:

- `accounts`: list of accounts for which AccountPrivileges are to be queried.
- `privileges`: list of privileges associated with each account.
- `pagination`: page response parameters.

##### QueryGetGuardTransferCoinsRequest

Request type for the GuardTransferCoins RPC method.

##### QueryGetGuardTransferCoinsResponse

Response type for the GuardTransferCoins RPC method.

Fields:

- `guard_transfer_coins`: whether GuardTransferCoins is enabled or not

##### QueryGetRequiredPrivilegesRequest

Request type for the RequiredPrivileges RPC method.

Fields:

- `index`: index of the RequiredPrivileges item to be queried.
- `kind`: kind of RequiredPrivileges item to be queried.

##### QueryGetRequiredPrivilegesResponse

Response type for the RequiredPrivileges RPC method.

Fields:

- `index`: index of the RequiredPrivileges item that was queried.
- `privileges`: privileges associated with the RequiredPrivileges item.
- `kind`: kind of RequiredPrivileges item that was queried.

##### QueryAllRequiredPrivilegesRequest

Request type for the RequiredPrivilegesAll RPC method.

Fields:

- `pagination`: page request parameters.
- `kind`: kind of RequiredPrivileges items to be queried.

##### QueryAllRequiredPrivilegesResponse

Response type for the RequiredPrivilegesAll RPC method.

Fields:

- `indexes`: list of indexes of the RequiredPrivileges items.
- `privileges`: list of privileges associated with each RequiredPrivileges item.
- `kind`: kind of RequiredPrivileges items that were queried.
- `pagination`: page response parameters.

##### QueryGetLockedRequest

Request type for the Locked RPC method.

Fields:

- `index`: index of the Locked item to be queried.
- `kind`: kind of Locked item to be queried.

##### QueryGetLockedResponse

Response type for the Locked RPC method.

Fields:

- `index`: index of the Locked item that was queried.
- `locked`: whether the item is locked or not.
- `kind`: kind of Locked item that was queried.

##### QueryAllLockedRequest

Request type for the LockedAll RPC method.

Fields:

`pagination`: page request parameters.
`kind`: kind of Locked items to be queried.

##### QueryAllLockedResponse

Response type for the LockedAll RPC method.

Fields:

- `indexes`: list of indexes of the Locked items.
- `locked`: list of boolean values indicating whether each item is locked or not.
- `kind`: kind of Locked items that were queried.
- `pagination`: page response parameters.

## Events

The Guard module emits the following events:

### MsgUpdateAccountPrivileges: Emits an event when the account privileges are updated

- action: `update_account_privileges`
- creator: The address of the user who created the message.
- account: The address of the account whose privileges were updated.

### MsgUpdateAccountPrivilegesBatch: Emits an event when the account privileges are updated in a batch

- action: `update_account_privileges_batch`
- creator: The address of the user who created the message.
- accounts: The list of addresses of the accounts whose privileges were updated.

### MsgUpdateAccountPrivilegesGroupedBatch: Emits an event when the account privileges are updated in a grouped batch

- action: `update_account_privileges_grouped_batch`
- creator: The address of the user who created the message.
- accounts: The list of addresses of the accounts whose privileges were updated.

### MsgUpdateGuardTransferCoins: Emits an event when the guard transfer coins functionality is updated

- action: `update_guard_transfer_coins`
- creator: The address of the user who created the message.
- guard_transfer_coins: The updated guard transfer coins status.

### MsgUpdateLocked: Emits an event when the locked status is updated

- action: `update_locked`
- creator: The address of the user who created the message.
- locked: The updated locked status.
- index: The index of the locked status that was updated.
- kind: The type of the locked status.

### MsgUpdateRequiredPrivileges: Emits an event when the required privileges are updated for a specific index

- action: `update_required_privileges`
- creator: The address of the user who created the message.
- index: The index of the required privileges that was updated.
- kind: The type of the required privileges.

### MsgUpdateRequiredPrivilegesBatch: Emits an event when the required privileges are updated for a batch of indexes

- action: `update_required_privileges_batch`
- creator: The address of the user who created the message.
- indexes: The list of indexes of the required privileges that were updated.
- kind: The type of the required privileges.

### MsgUpdateRequiredPrivilegesGroupedBatch: Emits an event when the required privileges are updated for a grouped batch of indexes

- action: `update_required_privileges_grouped_batch`
- creator: The address of the user who created the message.
- indexes: The list of indexes of the required privileges that were updated.
- kind: The type of the required privileges.

## MsgParams

The `Params` message defines the parameters for the `guard` module.

| Field | Type | Description |
|-------|------|-------------|
| admin_account | string | The address of the admin account that can update the parameters. |
| account_privileges_token_collection_creator | string | The address of the user who created the token collection used for account privileges. |
| account_privileges_token_collection_id | string | The ID of the token collection used for account privileges. |
| default_privileges | bytes | The default privileges that will be set for new accounts. |

## Dependencies

The Guard module has the following dependencies:

- The codec of the application.
- The store key for the Guard module.
- The memory store key for the Guard module.
- The subspace for the Guard module in the main application's parameter store.
- The addresses of the module accounts.
- The message service router of the application.
- The account keeper of the application.
- The authorization keeper of the application.
- The token keeper of the application.
- The NFT keeper of the application.
- The coin factory keeper of the application.

## Example Usage

### Client

#### SDK

##### Prerequisites
 
1. Install the @mantra@sdk. 

```
npm i --save @mantra@sdk
```
2. Initialize the client.
```javascript
import { Client } from '@mantrachain/sdk';
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

(async () => {
  const mnemonic ="..."; // MNEMONIC
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
    prefix: "mantrachain",
  });

  const creator = (await wallet.getAccounts())[0].address;

  const client = new Client(
    {
      apiURL: "http://127.0.0.1:1317",
      rpcURL: "http://127.0.0.1:26657",
      prefix: "mantrachain",
    },
    wallet
  );
})();
```

##### Transactions

<span style="text-decoration: underline">ACCOUNT PRIVILEGES AND REQUIRED PRIVILEGES</span>

This module is primarily responsible for handling all the public facing user privileges permissions. 
Separate from the network permissions, these permissions are on the basis of a user account.

It works via the use of bitwise operations, with each bit representing a permission for the user. So
in this way we are able to create a multi permission requirement performant and cost effectively.

Currently, thee only supported permissions are related to tokens transfer for coins minted by the `CoinFactory` module.
NOTE: The default privileges are set to 0, which means that the user has no permissions.

Account privileges and required privileges are stored in the form of a bytes array.

When a user wants to transfer a coin minted by the `CoinFactory` module, the following checks are performed:

- The user must have `soul-bond nft` in a restricted collection. That nft is minted by the quilified operators of the chain.  
- The user must have the `transfer` permission set to `true` in the `account_privileges` module. Which byte of the array is the `transfer` permission is defined by the operator.
- The coin minted by the `CoinFactory` module must have the `transfer` permission set to `true` in the `required_privileges` module. Which byte of the array is the `transfer` permission is defined by the operator.
- The coin should not be locked.

The privileges array are represented by 32 bytes arrays. Each bit represents a permission.
The first 8 bytes are reserved. The rest of the bytes are used for the coins transfers permissions.

- 0 reserved for `not-blacklisted`
- 1 reserved for `not-locked`
- 2-63 reserved for system values in the future
- 0-63 will be 1 by default (affirmations) (2^64-1)

When the operator wants to set the `transfer` permission for the `account_privileges` module, they need to set the first reserved 8 bytes to 1 along with the byte that they want to set the `transfer` permission.

The same applies for the `required_privileges` for the `CoinFactory` issued coins.

1. Set the first reserved 8 bytes to 1 along with the byte that the operators want to set the `transfer` permission.

```javascript
  await client.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivileges({
    value: {
      creator,
      index: Buffer.from("testcoin"), // denom of the coin from the `CoinFactory` module represented in bytes
      privileges: Buffer.from([1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]),
      kind: 'coin'
    },
  })
```

2. Set the first reserved 8 bytes to 1 along with the byte that the operators want to set the `transfer` permission.

```javascript
  await client.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivileges({
    value: {
      creator,
      account: 'mantrachain...',
      privileges: Buffer.from([1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]),
    }
  })
```

Once it's done the user will be able to transfer the coin minted by the `CoinFactory` module if the have issued the restricted soul-bond nft and the coin is not locked.

NOTE: If the user doesn't have the `transfer` permission set to `true` in the `account_privileges` module, the user will not be able to transfer the coin even if the coin has the `transfer` permission set to `true` in the `required_privileges` module. 

NOTE: If the user doesn't have account privileges set, the default privileges will be used. The default privileges have 8 bytes set to 0xff, but that doesn't means that the user has permissions to transfer coins.

<b>Additional transactions includes update the account privileges for a batch of accounts, update the account privileges for a batch of grouped accounts, update the required privileges for a batch of coins, and update the required privileges for a grouped batch of coins.</b>

<span style="text-decoration: underline">UPDATE GUARD TRANSFER COINS</span>

```javascript
  await client.MantrachainGuardV1.tx.sendMsgUpdateGuardTransferCoins({
    value: {
      creator,
      enabled: true
    }
  })
```

Once set to true the user won't be able to transfer coins if they don't have the `transfer` permission set to `true` in the `account_privileges` permissions. The admin of the coin and the general admin account make an exception.
If the coin is locked or the coin doesn't have required privileges set by a operator, the user won't be able to transfer the coin.

<span style="text-decoration: underline">UPDATE LOCKED</span>

```javascript
  await client.MantrachainGuardV1.tx.sendMsgUpdateLocked({
    value: {
      creator,
      index: Buffer.from("testcoin"), // denom of the coin from the `CoinFactory` module represented in bytes
      locked: true,
      kind: 'coin'
    }
  })
```

If a coin is locked the users won't be able to transfer that coin.

<span style="text-decoration: underline">UPDATE AUTZ GENERIC GRANT REVOKE BATCH</span>

MsgUpdateAuthzGenericGrantRevokeBatch is used by the main admin of the chain to grant or revoke roles permissions to the operators. It uses the cosmos-sdk `authz` module. It performs batch operation with the messages a operator can execute on a behalf of the main admin.

```javascript
  await client.MantrachainGuardV1.tx.sendMsgUpdateAuthzGenericGrantRevokeBatch({
    value: {
      creator,
      grantee: 'mantrachain...',
      authzGrantRevokeMsgsTypes: {
        msgs: [
          { typeUrl: '/mantrachain.guard.v1.MsgUpdateAccountPrivileges', grant: true},
          { typeUrl: '/mantrachain.guard.v1.MsgUpdateLocked', grant: false}
        ],
      }
    }
  })
```

Once some operator has been granted the `MsgUpdateAccountPrivileges` permission, they will be able to update the account privileges for accounts. This happens through the cosmos-sdk `autz` module.

##### Queries

<span style="text-decoration: underline">PARAMS</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryParams()
  console.log(res.data.params);
```

Output:

```javascript
{
  admin_account: 'mantrachain...',
  account_privileges_token_collection_creator: 'mantrachain...',
  account_privileges_token_collection_id: 'account_privileges_guard_nft_collection',
  default_privileges: '//////////8='
}
```

The `default_privileges` is the default privileges that will be used for the accounts that don't have account privileges set. The default privileges have 8 bytes set to 0xff, but that doesn't means that the user has permissions to transfer coins. They are base64 stringified representation of 8 bytes set to 0xff.

To get the default privileges in bytes:

```javascript
  const defaultPrivileges = Buffer.from(res.data.params.default_privileges, "base64")
  console.log(defaultPrivileges)
```

Output:

```shell
<Buffer ff ff ff ff ff ff ff ff>
```

<span style="text-decoration: underline">ACCOUNT PRIVILEGES</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryAccountPrivileges("mantrachain...")
  console.log(res.data);
```

Output:

```shell
{
  account: 'mantrachain...',
  privileges: 'Af//////////'
}
```

To get the account privileges in bytes:

```javascript
  const privileges = Buffer.from(res.data.privileges, "base64")
  console.log(privileges)
```

Output:

```shell
<Buffer 01 ff ff ff ff ff ff ff ff>
```

<span style="text-decoration: underline">ACCOUNT PRIVILEGES ALL</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryAccountPrivilegesAll()
  console.log(res.data)
```

Output:

```shell
{
  accounts: [ 'mantrachain...' ],
  privileges: [ 'Af//////////' ],
  pagination: { next_key: null, total: '1' }
}
```

To get the account privileges in bytes:
  
```javascript
  const privileges = Buffer.from(res.data.privileges[0], "base64")
  console.log(privileges)
```

Output:

```shell
<Buffer 01 ff ff ff ff ff ff ff ff>
```

<span style="text-decoration: underline">GUARD TRANSFER COINS</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryGuardTransferCoins()
  console.log(res.data)
```

Output:

```shell
{ guard_transfer_coins: true }
```

<span style="text-decoration: underline">REQUIRED PRIVILEGES</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryRequiredPrivileges(Buffer.from("testcoin").toString("base64"), {
    kind: 'coin'
  })
  console.log(res.data)
```

Output:

```shell
{ index: 'bXljb2lu', privileges: 'Af//////////', kind: 'coin' }
```

To get the required privileges in bytes:
  
```javascript
  const privileges = Buffer.from(res.data.privileges, "base64")
  console.log(privileges)
```

Output:

```shell
<Buffer 01 ff ff ff ff ff ff ff ff>
```

To get the index in string:
  
```javascript
  const index = Buffer.from(res.data.index, "base64").toString('utf8')
  console.log(index)
```

Output:

```shell
testcoin
```

<span style="text-decoration: underline">REQUIRED PRIVILEGES ALL</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryRequiredPrivilegesAll({
    kind: 'coin'
  })
  console.log(res)
```

Output:

```shell
{
  indexes: [ 'bXljb2lu' ],
  privileges: [ 'Af//////////' ],
  kind: 'coin',
  pagination: { next_key: null, total: '1' }
}
```

To get the required privileges in bytes:
  
```javascript
  const privileges = Buffer.from(res.data.privileges[0], "base64")
  console.log(privileges)
```

Output:

```shell
<Buffer 01 ff ff ff ff ff ff ff ff>
```

To get the index in string:
  
```javascript
  const index = Buffer.from(res.data.indexes[0], "base64").toString('utf8')
  console.log(index)
```

Output:

```shell
testcoin
```

<span style="text-decoration: underline">LOCKED</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryLocked(Buffer.from("testcoin").toString("base64"), {
    kind: 'coin'
  })
  console.log(res.data)
```

Output:

```shell
{ index: 'dGVzdGNvaW4=', locked: true, kind: 'coin' }
```

To get the index in string:
  
```javascript
  const index = Buffer.from(res.data.index, "base64").toString('utf8')
  console.log(index)
```

Output:

```shell
testcoin
```

Note: If the coin doesn't exist on the chain, the `locked` field will be `false`.

<span style="text-decoration: underline">LOCKED ALL</span>

```javascript
  const res = await client.MantrachainGuardV1.query.queryLockedAll({
    kind: 'coin'
  })
  console.log(res.data)
```

Output:

```shell
{
  indexes: [ 'dGVzdGNvaW4=' ],
  locked: [ true ],
  kind: 'coin',
  pagination: { next_key: null, total: '1' }
}
```

To get the index in string:
  
```javascript
  const index = Buffer.from(res.data.indexes[0], "base64").toString('utf8')
  console.log(index)
```

Output:

```shell
testcoin
```
