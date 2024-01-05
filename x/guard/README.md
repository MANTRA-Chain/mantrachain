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

## Background

The Guard module is responsible for managing all the permissions and privileges of the users in the network. It allows other modules to define their own permissions and use the Guard module to enforce them.

The Guard module provides a set of functions for updating account privileges, updating required privileges, updating guard transfer coins and updating authz grant/revoke execute permissions.

The Guard module is responsible for:

1. Guarding against unauthorized actions.
2. Handling all the privileges and permissions.
3. Managing the permission system.

## Concepts

The Guard module defines several concepts:

- **Accounts privileges:** Actions that can be performed within the system. Each privilege has an associated bytes array and an account that have that privilege.
- **Required privileges:** The privileges required to perform an action. These can be set by the admin or by the system operators.
- **Guarded transfer of coins:** A setting that can be enabled or disabled which prevents coins from being transferred without proper authorization.

## Installation

To install the Guard module, you can use the following command:

```bash
go get aumega/x/guard
```

## State

The Guard module stores the following data:

- **Accounts privileges:** A list of all the privileges of associated accounts.
- **Required privileges:** A list of the required privileges for coins transfers.
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
`/aumega/guard/v1/params`

##### AccountPrivileges

Queries an AccountPrivileges item by account.

- Request Type: `QueryGetAccountPrivilegesRequest`
- Response Type: `QueryGetAccountPrivilegesResponse`

HTTP GET Path:
`/aumega/guard/v1/account_privileges/{account}`

##### AccountPrivilegesAll

Queries a list of AccountPrivileges items.

- Request Type: `QueryAllAccountPrivilegesRequest`
- Response Type: `QueryAllAccountPrivilegesResponse`

HTTP GET Path:
`/aumega/guard/v1/account_privileges`

##### GuardTransferCoins

Queries a GuardTransferCoins item.

- Request Type: `QueryGetGuardTransferCoinsRequest`
- Response Type: `QueryGetGuardTransferCoinsResponse`

HTTP GET Path
`/aumega/guard/v1/guard_transfer_coins`

##### RequiredPrivileges

Queries a RequiredPrivileges item by index.

- Request Type: `QueryGetRequiredPrivilegesRequest`
- Response Type: `QueryGetRequiredPrivilegesResponse`

HTTP GET Path
`/aumega/guard/v1/required_privileges/{index}`

##### RequiredPrivilegesAll

Queries a list of RequiredPrivileges items.

- Request Type: `QueryAllRequiredPrivilegesRequest`
- Response Type: `QueryAllRequiredPrivilegesResponse`

HTTP GET Path
`/aumega/guard/v1/required_privileges`

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

| Field                                       | Type   | Description                                                                           |
| ------------------------------------------- | ------ | ------------------------------------------------------------------------------------- |
| admin_account                               | string | The address of the admin account that can update the parameters.                      |
| account_privileges_token_collection_creator | string | The address of the user who created the token collection used for account privileges. |
| account_privileges_token_collection_id      | string | The ID of the token collection used for account privileges.                           |
| default_privileges                          | bytes  | The default privileges that will be set for new accounts.                             |
| base_denom                                  | string | The base denomination of the network.                                                 |

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
