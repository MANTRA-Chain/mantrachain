# ADR 009: Use depinject when possible

## Context

Depinject is a new cosmos SDK feature that is not yet fully supported up and down the cosmos stack.  

We observed that it causes an enormous reduction in the complexity of our code, and that it can be stubbed throughout our chain for modules that don't yet support it, like wasm.go and ibc.go in the rescaffolding branch.

## Decision

We decided not to use depinject at this time because it is not yet fully supported up and down the cosmos stack.  We will revisit this decision when it is more widely supported.
For now, we will continue to manually register all of our dependencies in the app.go file.

## Status

Rejected
