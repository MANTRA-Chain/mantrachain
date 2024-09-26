# ADR 010: Maintain a cosmos-sdk fork for the bank module

## Context

Neutron and Osmosis both support so-called before send hooks via changes to x/bank.  Our initial plan was to have no fork at all, but the BeforeSend features proved very attractive for RWAs, and a fork was created at:

* https://github.com/MANTRA-Chain/cosmos-sdk

## Decision

We'll maintain a cosmos-sdk fork, currently from v0.50.x.  The version of the SDK we are currently working with will be set to the default branch, and dependabot and other tools will be enabled.  

* We will merge dependabot updates when possible, because these will provide us with indicaitons of compatiblity. 
* We will not fix lints, since that would create too large of a diff (upstream doesn't fix lints once a branch has had a release).
* We will routinely update the cosmos-sdk fork to ensure that our work benefits from the latest fixes from upstream.
* Changes to our SDK fork outside of what exists in x/bank require an additional ADR.
* If BeforeSend hooks were supported in the mainline SDK, we would stop using a fork.

## Status

Accepted
