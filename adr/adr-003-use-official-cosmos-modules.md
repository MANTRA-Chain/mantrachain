# ADR 003: Using of official Cosmos-SDK and IBC-Go

## Context

We are currently importing a forked version of IBC-Go due to the way our guard checks are setup. Previously we were using another forked version of Cosmos-SDK but switched to the official version since 0.50 since they added `SetSendRestrictions`.

## Alternatives

### Continue using forked versions if we require module modifications but minimize it

Reassess whether official module modifications are necessary. If they are without any other workarounds, we continue to import the forked repositories and make the changes there.

## Decision

Move to only use official imports of Comsos-SDK, IBC-Go and other future official modules maintained by the Cosmos team. All code modifications should be in mantrachain and if a official module need to be modified, we should copy the module and modify it in `/x`. The modified module should be written in a way where it imports from the files of the copied module as much as possible with minimal copied code.

## Status

Accepted

## Consequences

### Positive

* Improved readability and maintainability of code.
* We no longer require forked versions of Cosmos modules.

### Neutral

* If we have modified modules in `/x` in the future, we will need to see if they are compatible whenever we want to upgrade the Cosmos-SDK or whichever related official Cosmos repository.
