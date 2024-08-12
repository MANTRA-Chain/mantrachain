# ADR 007: KVStore key prefixes to be single/double byte(s)

## Context

Many of our KVStore have both store key as well as a index prefix. This effectively means that the the store keys are prefixed with both the store key as well as the index prefix. They are also all in descriptive strings.

## Decision

We should remove all store key for these stores that have index prefix. The index prefix should also be changed to single/double byte(s). This makes the set/get code much neater and maintainable. There will also be slight performance improvements with these changes.

## Status

Proposed
