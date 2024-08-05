# ADR 002: Moving of guard permission checks to Message level

## Context

Current guard checks are all done at the cosmos-sdk/x/bank SendCoins function level with the restriction registered [here](/x/guard/keeper/keeper.go#L72). This is extremely problematic as the SendCoins function is called by multiple internal functions and is critical to the functionality of mantrachain. The current architecture of allowing the native fee token to always bypass the check was also due to the fact that blocking the SendCoins function for the native fee token can have unintended consequences.

Currently, whenever a module/generated address were to send/receive tokens, they need to be whitelisted before the SendCoins function is called and immediately removed from whitelist afterwards. Below is a code snippert from [airdrop CreateCampaign](/x/airdrop/keeper/msg_server_campaign.go#L34) to show how it is currently.

```go
whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{campaign.GetReserveAddress().String()})
err = k.bankKeeper.SendCoins(ctx, creator, campaign.GetReserveAddress(), campaign.Amounts)
k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
```

## Alternatives

### Keep the checks at the keeper.SendCoins as it is

This ensure that there is no 'accidental' way users without authorization can transfer tokens. However, it makes our code base difficult to maintain and forces us to fork and modify the IBC-go transfer module.

### Move checks to the anteHandle

The checks can be done at the antehandler where it look for specific messages to check. This will cause the transactions to not cause any gas should the user not have the required authorization. Unmitigated, this could leave the chain vulnerable to DDoS attacks, potentially halting block production.

## Decision

Move all checks from the keeper SendCoins level to the message server level for functions Send and Multisend.

This way all the AddTransferAccAddressesWhitelist and RemoveTransferAccAddressesWhitelist code can be removed that are associated with module/generated addresses

## Status

Draft

## Consequences

### Positive

* Improved readability and maintainability of code.
* We no longer require a forked version of the IBC-go and can switch to using the official IBC-go.

### Negative

* We might need to either fork the cosmos-sdk/x/bank folder or copy it into mantrachain to add the checks at the message server.
* There might be other messages or a combination of them that users can invoke to bypass these checks and transfer tokens from/to unauthorized accounts.

## Further Discussions

We should check each message callable by EOAs to check which messages can lead to a transfer of tokens. We also need to know if the checks can be bypassed with CosmWasm contracts. This can be done with the help of auditors.
