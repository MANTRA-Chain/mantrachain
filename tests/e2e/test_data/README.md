# Wasm Contracts for before-send-hook

This directory contains WASM contracts designed to be used as `before-send-hook` for TokenFactory tokens.

## Purpose

These contracts enable additional execution and validation whenever a specified TokenFactory token is sent. The `before-send-hook` mechanism is configured to be denom-specific, meaning a different contract or configuration can be applied to each unique TokenFactory denom.

## How it Works

When a TokenFactory token with an associated `before-send-hook` is transferred, the configured WASM contract is executed *before* the actual token transfer occurs. This allows the contract to:

*   **Perform Custom Logic:** Execute arbitrary logic defined within the WASM contract.
*   **Validate Transactions:** Check conditions, enforce rules, or verify sender/receiver addresses.
*   **Prevent Transfers:** Revert the transaction if certain conditions are not met.

## Use Cases

*   **Compliance Checks:** Ensure token transfers adhere to regulatory requirements.
*   **Fee Collection:** Implement custom fees for token transfers.
*   **Whitelisting/Blacklisting:** Restrict token transfers to/from specific addresses.
*   **Token Gating:** Require certain conditions (e.g., holding another token) for transfers.
*   **Inter-chain Security:** Add custom logic for tokens bridged from other chains.

## Contract Structure

Each WASM contract in this directory is designed to be deployed on the chain and then registered as a `before-send-hook` for a specific TokenFactory denom. The contract's entry points (e.g., `execute`, `query`) will be invoked by the `before-send-hook` mechanism.

## Deployment and Configuration

1.  **Compile:** Compile the Rust-based WASM contracts into `.wasm` bytecode.
2.  **Upload:** Upload the `.wasm` bytecode to the chain using the `wasm` module's `store-code` transaction.
3.  **Instantiate:** Instantiate the uploaded code to create a contract instance.
4.  **Set Hook:** Use the TokenFactory module's functionality (e.g., `set-before-send-hook`) to associate the instantiated contract's address with a specific TokenFactory denom.
## Included Contracts

The contracts included here are simplified for testing purposes and are utilized in `e2e_tokenfactory_test.go`.

*   **`percentage_cap.wasm`**: This contract ensures that no more than 50% of the total supply of a given denom can be transferred in a single bank transfer.
*   **`transfer_cap.wasm`**: This contract enforces a maximum transfer limit, ensuring that no more than 1,000,000 units of a denom can be transferred in a single bank transfer.
*   **`tax_stake_denom.wasm`**: This contract implements a tax mechanism. It requires users to deposit the staking denom (`amantra`) into the contract. Whenever the associated TokenFactory denom is transferred, regardless of the amount, `1000000000000000000amantra` will first be transferred to the contract's admin from the sender of the associated Tokenfactory denom.