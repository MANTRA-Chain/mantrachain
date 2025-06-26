[![License](https://img.shields.io/github/license/MANTRA-Chain/mantrachain)](https://github.com/MANTRA-Chain/mantrachain/blob/main/LICENSE)

# Mantrachain

Mantrachain is a global real-world assets platform built on blockchain technology. It leverages advanced blockchain features to facilitate the tokenization and trading of real-world assets.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Joining the Mainnet](#joining-the-mainnet)
- [Getting Started](#getting-started)
- [Development](#development)
- [Architecture](#architecture)
- [Modules](#modules)
- [Security](#security)

## Overview

Mantrachain is designed to bridge the gap between traditional assets and the blockchain world. By enabling the tokenization of real-world assets, it opens up new possibilities for asset management, trading, and financial innovation.

## Features

- Real-world asset tokenization
- Advanced blockchain technology integration
- Multi-token support for transaction fees
- Custom fee market implementation
- Cosmos SDK-based architecture

## Joining the Mainnet

[Please visit the official instructions on how to join the Mainnet here.](https://docs.mantrachain.io/node-and-validator-operations/node-setup-and-deployment/running-a-node)

## Getting Started

To get started with Mantrachain, you'll need to set up your development environment.

### Prerequisites

- Go 1.23 or later

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MANTRA-Chain/mantrachain.git
   cd mantrachain
   ```

2. Build the project:
   ```bash
   make install
   ```

## Development



### Testing

To run unit tests:

```bash
make test-unit
```



## Architecture

Mantrachain follows the Cosmos SDK architecture and implements several custom modules to achieve its functionality. The project uses Architecture Decision Records (ADRs) to document important architectural decisions.

For more information on the architecture and design decisions, please refer to the [ADR directory](adr/).

## Modules

Mantrachain includes several custom modules:

- `x/sanction`: Blacklisting of addresses to prevent transactions from sanctioned entities.
- `x/tokenfactory`: Allows for the creation and management of new tokens (based on Neutron's implementation).
- `x/tax`: Handles tax-related operations within the chain.

For detailed information on each module, please refer to their respective README files in the `x/` directory.

## Security

We take security seriously. If you discover a security issue, please bring it to our attention right away!

Please refer to our [Security Policy](SECURITY.md) for more details on reporting vulnerabilities.




---

For more detailed information, please check the documentation in the respective directories and files within the repository.