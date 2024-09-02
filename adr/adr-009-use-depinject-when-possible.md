# ADR 009: Use depinject when possible

## Context

Depinject is a new cosmos SDK feature that is not yet fully supported up and down the cosmos stack.  

We observed that it causes an enormous reduction in the complexity of our code, and that it can be stubbed throughout our chain for modules that don't yet support it, like wasm.go and ibc.go in the rescaffolding branch.

## Decision

When possible, we'll use depinject to integrate mantra's external dependencies.


## Status

Accepted
