## [unreleased]

### FEATURES

### BUG-FIXES

- Update version extraction regex for SDK version detection ([\#487](https://github.com/MANTRA-Chain/mantrachain/issues/487))

### DEPENDENCIES

### DOCUMENTATION

- Add CHANGELOG.md ([\#516](https://github.com/MANTRA-Chain/mantrachain/issues/516))

### CI

### OTHER

- Add template for v7 upgrade ([\#486](https://github.com/MANTRA-Chain/mantrachain/issues/486))

## [6.1.1] - 2025-10-30

### BUG-FIXES

- Historical EvmCoinInfo don't exists ([\#503](https://github.com/MANTRA-Chain/mantrachain/issues/503)) ([\#505](https://github.com/MANTRA-Chain/mantrachain/issues/505))

### CI

- *(release)* Add gpg signing on the binaries checksum ([\#457](https://github.com/MANTRA-Chain/mantrachain/issues/457)) ([\#507](https://github.com/MANTRA-Chain/mantrachain/issues/507))

## [6.1.0] - 2025-10-27

### BUG-FIXES

- Evm fix not merged ([\#497](https://github.com/MANTRA-Chain/mantrachain/issues/497)) ([\#499](https://github.com/MANTRA-Chain/mantrachain/issues/499))
- V6.1.0 upgrade handler not added ([\#500](https://github.com/MANTRA-Chain/mantrachain/issues/500))

## [6.0.1] - 2025-10-30

### BUG-FIXES

- Historical EvmCoinInfo don't exists ([\#503](https://github.com/MANTRA-Chain/mantrachain/issues/503)) ([\#504](https://github.com/MANTRA-Chain/mantrachain/issues/504))
- Remove duplicate backport rule for release/v6.0.x branch

### CI

- *(release)* Add gpg signing on the binaries checksum ([\#457](https://github.com/MANTRA-Chain/mantrachain/issues/457)) ([\#506](https://github.com/MANTRA-Chain/mantrachain/issues/506))

## [6.0.0] - 2025-10-22

### FEATURES

- Deprecate x/group ([\#459](https://github.com/MANTRA-Chain/mantrachain/issues/459))
- Additional blacklist for wasm msgserver ([\#464](https://github.com/MANTRA-Chain/mantrachain/issues/464))
- Add comet query commands ([\#475](https://github.com/MANTRA-Chain/mantrachain/issues/475))

### BUG-FIXES

- Error ignored by scoped error variables ([\#452](https://github.com/MANTRA-Chain/mantrachain/issues/452))
- Change chain id back to default after sig verify ([\#453](https://github.com/MANTRA-Chain/mantrachain/issues/453))
- Owner is not aligned in tokenfactory token\_pairs for ibc transfer ([\#460](https://github.com/MANTRA-Chain/mantrachain/issues/460))
- Authz module is risky ([\#462](https://github.com/MANTRA-Chain/mantrachain/issues/462))
- Static precompiles are not enabled ([\#465](https://github.com/MANTRA-Chain/mantrachain/issues/465))
- Not set disable list of wasm breaker in upgrade ([\#466](https://github.com/MANTRA-Chain/mantrachain/issues/466))
- Cometbft is outdated ([\#477](https://github.com/MANTRA-Chain/mantrachain/issues/477))
- Evm is outdated ([\#478](https://github.com/MANTRA-Chain/mantrachain/issues/478))
- Upstream evm is not integrated ([\#468](https://github.com/MANTRA-Chain/mantrachain/issues/468))
- Miss coin info in connect test ([\#484](https://github.com/MANTRA-Chain/mantrachain/issues/484)) ([\#490](https://github.com/MANTRA-Chain/mantrachain/issues/490))
- Upstream evm v0.5.0 not used ([\#491](https://github.com/MANTRA-Chain/mantrachain/issues/491)) ([\#493](https://github.com/MANTRA-Chain/mantrachain/issues/493))
- Static precompiles are not enabled in upgrade handler ([\#494](https://github.com/MANTRA-Chain/mantrachain/issues/494))
- Static precompiles param is not set ([\#495](https://github.com/MANTRA-Chain/mantrachain/issues/495))
- Adding back v6rc0 ([\#496](https://github.com/MANTRA-Chain/mantrachain/issues/496))

### DEPENDENCIES

- *(deps)* Bump github.com/ulikunitz/xz from v0.5.11 to v0.5.15 ([\#450](https://github.com/MANTRA-Chain/mantrachain/issues/450))
- *(deps)* Bump [github.com/Zondax/ledger-go](https://github.com/Zondax/ledger-go) from 0.14.3 to 1.0.1 ([\#458](https://github.com/MANTRA-Chain/mantrachain/issues/458))
- *(deps)* Update cosmos-sdk fork from v0.53.3 to v0.53.4 ([\#488](https://github.com/MANTRA-Chain/mantrachain/issues/488)) ([\#489](https://github.com/MANTRA-Chain/mantrachain/issues/489))

### CI

- Read chain ID from app.toml before genesis ([\#454](https://github.com/MANTRA-Chain/mantrachain/issues/454))
- *(release)* Add github attestation on the binaries ([\#456](https://github.com/MANTRA-Chain/mantrachain/issues/456))

### OTHER

- Add template for v6 upgrade ([\#461](https://github.com/MANTRA-Chain/mantrachain/issues/461))

## [5.0.2] - 2025-10-15

### BUG-FIXES

- Cometbft is outdated ([\#477](https://github.com/MANTRA-Chain/mantrachain/issues/477)) ([\#480](https://github.com/MANTRA-Chain/mantrachain/issues/480))
- Evm is outdated ([\#478](https://github.com/MANTRA-Chain/mantrachain/issues/478)) ([\#481](https://github.com/MANTRA-Chain/mantrachain/issues/481))

### DEPENDENCIES

- Update backport rules for release/v5.0.x branch ([\#445](https://github.com/MANTRA-Chain/mantrachain/issues/445))

### OTHER

- Add template for v4 upgrade ([\#277](https://github.com/MANTRA-Chain/mantrachain/issues/277))
- Add template for v5 upgrade ([\#283](https://github.com/MANTRA-Chain/mantrachain/issues/283))
- *(golangci-lint)* Migrate to v2 config and apply lint to repo ([\#358](https://github.com/MANTRA-Chain/mantrachain/issues/358))
- Allow Unprotected Transactions ([\#372](https://github.com/MANTRA-Chain/mantrachain/issues/372))

## [5.0.0] - 2025-09-13

### FEATURES

- Ibc-v10 ([\#321](https://github.com/MANTRA-Chain/mantrachain/issues/321))
- Evm module with v5 upgrade ([\#345](https://github.com/MANTRA-Chain/mantrachain/issues/345))
- Block mint burn for disabled tokenfactory coins ([\#359](https://github.com/MANTRA-Chain/mantrachain/issues/359))
- V5.0.0 evm enhancements + bug fixes ([\#364](https://github.com/MANTRA-Chain/mantrachain/issues/364))

## [4.0.2] - 2025-09-08

### FEATURES

- Added x/circuit api ([\#305](https://github.com/MANTRA-Chain/mantrachain/issues/305)) ([\#306](https://github.com/MANTRA-Chain/mantrachain/issues/306))

### BUG-FIXES

- None ([\#437](https://github.com/MANTRA-Chain/mantrachain/issues/437))
- None ([\#433](https://github.com/MANTRA-Chain/mantrachain/issues/433)) ([\#438](https://github.com/MANTRA-Chain/mantrachain/issues/438))
- Behavioral change between proxygasmeter and standalone gas meter ([\#439](https://github.com/MANTRA-Chain/mantrachain/issues/439))

### CI

- Bump cosmos-sdk fork from v0.50.13-v3-mantra-1 to v0.50.14-v4-mantra-1 ([\#368](https://github.com/MANTRA-Chain/mantrachain/issues/368))

## [4.0.1] - 2025-04-03

### FEATURES

- Pump go in release.yml ([\#295](https://github.com/MANTRA-Chain/mantrachain/issues/295)) ([\#296](https://github.com/MANTRA-Chain/mantrachain/issues/296))

### BUG-FIXES

- Add openapi route ([\#299](https://github.com/MANTRA-Chain/mantrachain/issues/299)) ([\#300](https://github.com/MANTRA-Chain/mantrachain/issues/300))
- Remove broken links ([\#304](https://github.com/MANTRA-Chain/mantrachain/issues/304))
- Change okx websocket url ([\#310](https://github.com/MANTRA-Chain/mantrachain/issues/310))
- Update heighliner tag to v1.7.5 ([\#325](https://github.com/MANTRA-Chain/mantrachain/issues/325))
- Update workflow configurations for CodeQL and connect-test ([\#327](https://github.com/MANTRA-Chain/mantrachain/issues/327))

## [4.0.0] - 2025-03-25

### FEATURES

- Interchain accounts ([\#278](https://github.com/MANTRA-Chain/mantrachain/issues/278))
- Whitelist getprice query for wasm ([\#284](https://github.com/MANTRA-Chain/mantrachain/issues/284))
- Add grpc querier for wasm query plugins ([\#289](https://github.com/MANTRA-Chain/mantrachain/issues/289)) ([\#290](https://github.com/MANTRA-Chain/mantrachain/issues/290))

### DEPENDENCIES

- *(deps)* Bump [github.com/cosmos/ibc-go/v8](https://www.google.com/search?q=https://github.com/cosmos/ibc-go/v8) from v8.6.1 to v8.7.0 ([\#280](https://github.com/MANTRA-Chain/mantrachain/issues/280))

### CI

- Update connect to v2.3.0-v4-mantra-1 ([\#293](https://github.com/MANTRA-Chain/mantrachain/issues/293))
- Add license badge ([\#317](https://github.com/MANTRA-Chain/mantrachain/issues/317))
- Remove networks folder ([\#323](https://github.com/MANTRA-Chain/mantrachain/issues/323))
- Remove registerdenoms ([\#349](https://github.com/MANTRA-Chain/mantrachain/issues/349))
- Bump cosmos-sdk fork from v0.50.12-v2-mantra-1 to v0.50.12-v3-mantra-1 ([\#282](https://github.com/MANTRA-Chain/mantrachain/issues/282))
- Bump cosmos-sdk fork from v0.50.12-v3-mantra-1 to v0.50.13-v3-mantra-1 ([\#287](https://github.com/MANTRA-Chain/mantrachain/issues/287))
- Update connect tag ([\#292](https://github.com/MANTRA-Chain/mantrachain/issues/292))

## [3.0.3] - 2025-03-21

### CI

- Bump cosmos-sdk fork from v0.50.12-v3-mantra-1 to v0.50.13-v3-mantra-1 ([\#286](https://github.com/MANTRA-Chain/mantrachain/issues/286))

## [3.0.2] - 2025-03-18

### CI

- Bump cosmos-sdk fork from v0.50.12-v2-mantra-1 to v0.50.12-v3-mantra-1 ([\#281](https://github.com/MANTRA-Chain/mantrachain/issues/281))

## [3.0.1] - 2025-03-14

### DEPENDENCIES

- *(deps)* Bump [github.com/cosmos/ibc-go/v8](https://www.google.com/search?q=https://github.com/cosmos/ibc-go/v8) from v8.6.1 to v8.7.0 ([\#279](https://github.com/MANTRA-Chain/mantrachain/issues/279))

## [3.0.0] - 2025-03-04

### FEATURES

- Increase max-wasm-size to 3mb ([\#266](https://github.com/MANTRA-Chain/mantrachain/issues/266))
- Sanction module ([\#267](https://github.com/MANTRA-Chain/mantrachain/issues/267))

### DEPENDENCIES

- *(deps)* Bump cosmos-sdk MANTRA fork to v0.50.12 ([\#269](https://github.com/MANTRA-Chain/mantrachain/issues/269))
- *(deps)* Bump lycheeverse/lychee-action from 2.1.0 to 2.3.0 ([\#258](https://github.com/MANTRA-Chain/mantrachain/issues/258))
- *(deps)* Bump [github.com/cosmos/ibc-go/v8](https://www.google.com/search?q=https://github.com/cosmos/ibc-go/v8) from v8.5.2 to v8.6.1 ([\#272](https://github.com/MANTRA-Chain/mantrachain/issues/272))

### CI

- Upgrade cometbft to v0.38.17 + wasmvm to v2.2.2 ([\#265](https://github.com/MANTRA-Chain/mantrachain/issues/265))
- Use new connect fork that is based off v2.3.0 ([\#270](https://github.com/MANTRA-Chain/mantrachain/issues/270))

## [2.0.2] - 2025-02-21

### DEPENDENCIES

- *(deps)* Bump cosmos-sdk MANTRA fork to v0.50.12 ([\#268](https://github.com/MANTRA-Chain/mantrachain/issues/268))

## [2.0.0] - 2025-02-13

### FEATURES

- Add e2e tests ([\#252](https://github.com/MANTRA-Chain/mantrachain/issues/252))
- Upgrade to v2

### DEPENDENCIES

- *(deps)* Bump cosmossdk.io/client/v2 ([\#109](https://github.com/MANTRA-Chain/mantrachain/issues/109))
- *(deps)* Bump [github.com/bufbuild/buf](https://github.com/bufbuild/buf) from 1.41.0 to 1.42.0 ([\#119](https://github.com/MANTRA-Chain/mantrachain/issues/119))
- *(deps)* Bump github.com/skip-mev/connect/v2 ([\#120](https://github.com/MANTRA-Chain/mantrachain/issues/120))
- *(deps)* Bump [github.com/skip-mev/connect/v2](https://www.google.com/search?q=https://github.com/skip-mev/connect/v2) in /tests/connect ([\#121](https://github.com/MANTRA-Chain/mantrachain/issues/121))
- *(deps)* Bump [github.com/strangelove-ventures/interchaintest/v8](https://www.google.com/search?q=https://github.com/strangelove-ventures/interchaintest/v8) ([\#128](https://github.com/MANTRA-Chain/mantrachain/issues/128))
- *(deps)* Bump google.golang.org/grpc from 1.66.2 to 1.67.0 ([\#127](https://github.com/MANTRA-Chain/mantrachain/issues/127))
- *(deps)* Bump github.com/btcsuite/btcd from 0.22.0-beta to 0.24.0 in /tests/connect ([\#8](https://github.com/MANTRA-Chain/mantrachain/issues/8))
- *(deps)* Bump cosmossdk.io/api from 0.7.5 to 0.7.6 ([\#134](https://github.com/MANTRA-Chain/mantrachain/issues/134))
- *(deps)* Bump [github.com/CosmWasm/wasmvm/v2](https://www.google.com/search?q=https://github.com/CosmWasm/wasmvm/v2) from 2.1.2 to 2.1.3 ([\#133](https://github.com/MANTRA-Chain/mantrachain/issues/133))
- *(deps)* Bump github.com/skip-mev/connect/v2 in /tests/connect ([\#163](https://github.com/MANTRA-Chain/mantrachain/issues/163))
- *(deps)* Bump [github.com/bufbuild/buf](https://github.com/bufbuild/buf) from 1.42.0 to 1.43.0 ([\#165](https://github.com/MANTRA-Chain/mantrachain/issues/165))
- *(deps)* Bump google.golang.org/grpc from 1.67.0 to 1.67.1 ([\#166](https://github.com/MANTRA-Chain/mantrachain/issues/166))
- *(deps)* Bump github.com/skip-mev/connect/v2 from 2.0.1 to 2.1.0 ([\#164](https://github.com/MANTRA-Chain/mantrachain/issues/164))
- *(deps)* Bump google/osv-scanner-action from 1.8.5 to 1.9.0 ([\#170](https://github.com/MANTRA-Chain/mantrachain/issues/170))
- *(deps)* Bump golang.org/x/tools from 0.25.0 to 0.26.0 ([\#174](https://github.com/MANTRA-Chain/mantrachain/issues/174))
- *(deps)* Bump github.com/bufbuild/buf from 1.43.0 to 1.44.0 ([\#172](https://github.com/MANTRA-Chain/mantrachain/issues/172))
- *(deps)* Bump lycheeverse/lychee-action from 1.10.0 to 2.0.0 ([\#179](https://github.com/MANTRA-Chain/mantrachain/issues/179))
- *(deps)* Bump google.golang.org/protobuf ([\#177](https://github.com/MANTRA-Chain/mantrachain/issues/177))
- *(deps)* Bump github.com/bufbuild/buf from 1.44.0 to 1.45.0 ([\#183](https://github.com/MANTRA-Chain/mantrachain/issues/183))
- *(deps)* Bump lycheeverse/lychee-action from 2.0.0 to 2.0.2 ([\#188](https://github.com/MANTRA-Chain/mantrachain/issues/188))
- *(deps)* Bump elliptic ([\#198](https://github.com/MANTRA-Chain/mantrachain/issues/198))
- *(deps)* Bump github.com/bufbuild/buf from 1.45.0 to 1.46.0 ([\#197](https://github.com/MANTRA-Chain/mantrachain/issues/197))
- *(deps)* Bump [github.com/grpc-ecosystem/grpc-gateway/v2](https://www.google.com/search?q=https://github.com/grpc-ecosystem/grpc-gateway/v2) ([\#206](https://github.com/MANTRA-Chain/mantrachain/issues/206))
- *(deps)* Bump github.com/cometbft/cometbft from 0.38.12 to 0.38.13 ([\#196](https://github.com/MANTRA-Chain/mantrachain/issues/196))
- *(deps)* Bump [github.com/prometheus/client\_golang](https://github.com/prometheus/client_golang) ([\#193](https://github.com/MANTRA-Chain/mantrachain/issues/193))
- *(deps)* Bump [github.com/skip-mev/connect/v2](https://www.google.com/search?q=https://github.com/skip-mev/connect/v2) in /tests/connect ([\#209](https://github.com/MANTRA-Chain/mantrachain/issues/209))
- *(deps)* Bump [github.com/skip-mev/connect/v2](https://www.google.com/search?q=https://github.com/skip-mev/connect/v2) from 2.1.0 to 2.1.2 ([\#208](https://github.com/MANTRA-Chain/mantrachain/issues/208))
- *(deps)* Bump github.com/cometbft/cometbft in the go\_modules group ([\#215](https://github.com/MANTRA-Chain/mantrachain/issues/215))
- *(deps)* Bump cosmossdk.io/log from 1.4.1 to 1.5.0 ([\#220](https://github.com/MANTRA-Chain/mantrachain/issues/220))
- *(deps)* Bump [github.com/cosmos/ibc-go/v8](https://www.google.com/search?q=https://github.com/cosmos/ibc-go/v8) from 8.5.1 to 8.5.2 ([\#212](https://github.com/MANTRA-Chain/mantrachain/issues/212))
- *(deps)* Bump github.com/cometbft/cometbft ([\#214](https://github.com/MANTRA-Chain/mantrachain/issues/214))
- *(deps)* Bump cosmossdk.io/math ([\#227](https://github.com/MANTRA-Chain/mantrachain/issues/227))
- *(deps)* Bump cosmossdk.io/math in the go\_modules group ([\#228](https://github.com/MANTRA-Chain/mantrachain/issues/228))
- *(deps)* Bump github.com/stretchr/testify in /tests/connect ([\#229](https://github.com/MANTRA-Chain/mantrachain/issues/229))
- *(deps)* Bump google.golang.org/grpc from 1.67.1 to 1.68.1 ([\#235](https://github.com/MANTRA-Chain/mantrachain/issues/235))
- *(deps)* Bump golang.org/x/tools from 0.26.0 to 0.28.0 ([\#234](https://github.com/MANTRA-Chain/mantrachain/issues/234))
- *(deps)* Bump codecov/codecov-action from 4 to 5 ([\#223](https://github.com/MANTRA-Chain/mantrachain/issues/223))
- *(deps)* Bump lycheeverse/lychee-action from 2.0.2 to 2.1.0 ([\#218](https://github.com/MANTRA-Chain/mantrachain/issues/218))
- *(deps)* Bump cosmossdk.io/depinject from 1.0.0 to 1.1.0 ([\#216](https://github.com/MANTRA-Chain/mantrachain/issues/216))
- *(deps)* Bump github.com/cosmos/cosmos-db from 1.0.2 to 1.1.0 ([\#237](https://github.com/MANTRA-Chain/mantrachain/issues/237))

### CI

- Adjust cron schedule for CodeQL ([\#230](https://github.com/MANTRA-Chain/mantrachain/issues/230))
- Add sonar project properties ([\#232](https://github.com/MANTRA-Chain/mantrachain/issues/232))
- Use new cosmos-sdk fork that is based off v0.50.10 ([\#238](https://github.com/MANTRA-Chain/mantrachain/issues/238))
- Bump cosmos-sdk MANTRA fork to v0.50.11 ([\#250](https://github.com/MANTRA-Chain/mantrachain/issues/250))

## [1.0.3] - 2024-12-11

### CI

- Upgrade cosmossdk.io/math to v1.4.0 ([\#239](https://github.com/MANTRA-Chain/mantrachain/issues/239))

## [1.0.2] - 2024-11-08

### FEATURES

- Add the v1 - v2 migration logic
- Add cosmwasm
- Update guard to allow transfers when required coin privileges is not being set ([\#53](https://github.com/MANTRA-Chain/mantrachain/issues/53))
- Allow to delete unused image tag
- Integrate slinky ([\#66](https://github.com/MANTRA-Chain/mantrachain/issues/66))
- Feemarket with dependency injection ([\#177](https://github.com/MANTRA-Chain/mantrachain/issues/177))

### BUG-FIXES

- The name of static lib changes
- Update swagger for ibc ([\#201](https://github.com/MANTRA-Chain/mantrachain/issues/201))

### DEPENDENCIES

- *(deps)* Bump [github.com/cometbft/cometbft](https://github.com/cometbft/cometbft) from v0.38.13 to v0.38.14 ([\#213](https://github.com/MANTRA-Chain/mantrachain/issues/213))
- *(deps)* Bump [github.com/cometbft/cometbft](https://github.com/cometbft/cometbft) from v0.38.14 to v0.38.15 ([\#217](https://github.com/MANTRA-Chain/mantrachain/issues/217))

### CI

- Remove use of osv scanner ([\#204](https://github.com/MANTRA-Chain/mantrachain/issues/204))
- Remove unnecessary ci jobs ([\#205](https://github.com/MANTRA-Chain/mantrachain/issues/205))
- Remove legacy ADRs ([\#207](https://github.com/MANTRA-Chain/mantrachain/issues/207))

## [1.0.0] - 2024-10-23

### FEATURES

- Initial chain scaffolding and module setup, including Bridge, EVM, Guard, Coin Factory, DEX (Crescent), DID, Rewards, Airdrop, and CosmWasm.
- Core logic for Guard module, including account privileges, coin locking, and authz decorators.
- Major SDK upgrades (targeting v0.47 & v0.50) and Interchain CCV implementation.
- Added oracle genesis ([\#47](https://github.com/MANTRA-Chain/mantrachain/issues/47))
- Ibc rate limit ([\#67](https://github.com/MANTRA-Chain/mantrachain/issues/67))
- Visualize genesis ([\#75](https://github.com/MANTRA-Chain/mantrachain/issues/75))
- Add ability to query the tokenfactory module parameters and denom info ([\#168](https://github.com/MANTRA-Chain/mantrachain/issues/168))
- Update Genesis for RC1 and Tokenomics Gov Proposal alignment ([\#173](https://github.com/MANTRA-Chain/mantrachain/issues/173))
- Add default denom resolver ([\#175](https://github.com/MANTRA-Chain/mantrachain/issues/175))
- Wasm api ([\#191](https://github.com/MANTRA-Chain/mantrachain/issues/191))
- Release v1.0.0

### BUG-FIXES

- Extensive bug fixes across all new modules, focusing on Guard, CoinFactory, DEX, and swap fee logic.
- Addressed v1-v2 migration issues, including proto type mismatches and module parameter migration.
- Fixed numerous issues related to the rewards module (snapshots, purging) and Guard transfer/swap privilege checks.
- *(lint)* Commits from previous main ([\#63](https://github.com/MANTRA-Chain/mantrachain/issues/63))
- E2e tests jest global setup ([\#65](https://github.com/MANTRA-Chain/mantrachain/issues/65))
- Set local setup default swap fee rate param to zero ([\#72](https://github.com/MANTRA-Chain/mantrachain/issues/72))
- *(lint)* Revert lint-install to previous script
- Cli get params ([\#133](https://github.com/MANTRA-Chain/mantrachain/issues/133))
- Generated swagger + openapi
- Lint
- Simplify TokenFactory depinject ([\#36](https://github.com/MANTRA-Chain/mantrachain/issues/36))
- Use goreleaser for release ([\#41](https://github.com/MANTRA-Chain/mantrachain/issues/41))
- Enable vote extension ([\#48](https://github.com/MANTRA-Chain/mantrachain/issues/48))
- Update denoms and timings in genesis file ([\#52](https://github.com/MANTRA-Chain/mantrachain/issues/52))
- Missing rest api ([\#64](https://github.com/MANTRA-Chain/mantrachain/issues/64))
- Harmonize block gas ([\#91](https://github.com/MANTRA-Chain/mantrachain/issues/91))
- Allow docker and static linked binary built for PR ci and fix no… ([\#93](https://www.google.com/search?q=https%22github.com/MANTRA-Chain/mantrachain/issues/93))
- Use github.ref\_name instead ([\#95](https://github.com/MANTRA-Chain/mantrachain/issues/95))
- Align min\_base\_gas\_price for fee market ([\#123](https://github.com/MANTRA-Chain/mantrachain/issues/123))
- Remove the use of app\_config authority for tax ([\#135](https://github.com/MANTRA-Chain/mantrachain/issues/135))
- Error log in connect test ([\#143](https://github.com/MANTRA-Chain/mantrachain/issues/143))
- Discontinue the use of depinject ([\#160](https://github.com/MANTRA-Chain/mantrachain/issues/160))
- Simulate gas calculations ([\#187](https://github.com/MANTRA-Chain/mantrachain/issues/187))
- Downgrade cometbft-db due to forward compatibility issue ([\#189](https://github.com/MANTRA-Chain/mantrachain/issues/189))
- No license ([\#192](https://github.com/MANTRA-Chain/mantrachain/issues/192))
- Register uom as DefaultBondDenom ([\#195](https://github.com/MANTRA-Chain/mantrachain/issues/195))

### DEPENDENCIES

- *(deps)* Bump github.com/cosmos/interchain-security/v5 ([\#30](https://github.com/MANTRA-Chain/mantrachain/issues/30))
- *(deps)* Bump [github.com/bufbuild/buf](https://github.com/bufbuild/buf) from 1.40.0 to 1.40.1 ([\#31](https://github.com/MANTRA-Chain/mantrachain/issues/31))
- *(deps)* Bump golang.org/x/tools from 0.24.0 to 0.25.0 ([\#39](https://github.com/MANTRA-Chain/mantrachain/issues/39))
- *(deps)* Bump github.com/bufbuild/buf from 1.30.0 to 1.36.0 ([\#56](https://github.com/MANTRA-Chain/mantrachain/issues/56))
- *(deps)* Bump actions/setup-go from 3 to 5 ([\#35](https://github.com/MANTRA-Chain/mantrachain/issues/35))
- *(deps)* Bump github.com/spf13/viper from 1.18.2 to 1.19.0 ([\#48](https://github.com/MANTRA-Chain/mantrachain/issues/48))
- *(deps)* Bump cosmossdk.io/x/evidence from 0.1.0 to 0.1.1 ([\#51](https://www.google.com/search?q=https-%3Egithub.com/MANTRA-Chain/mantrachain/issues/51))
- *(deps)* Bump [github.com/cosmos/cosmos-sdk](https://github.com/cosmos/cosmos-sdk) from 0.50.7 to 0.50.9 ([\#62](https://github.com/MANTRA-Chain/mantrachain/issues/62))
- *(deps)* Bump github.com/cosmos/gogoproto from 1.5.0 to 1.6.0 ([\#78](https://github.com/MANTRA-Chain/mantrachain/issues/78))
- *(deps)* Bump cosmossdk.io/x/nft from 0.1.0 to 0.1.1 ([\#77](https://github.com/MANTRA-Chain/mantrachain/issues/77))
- *(deps)* Bump cosmossdk.io/x/upgrade from 0.1.3 to 0.1.4 ([\#76](https://github.com/MANTRA-Chain/mantrachain/issues/76))
- *(deps)* Bump cosmossdk.io/log from 1.3.1 to 1.4.0 ([\#75](https://github.com/MANTRA-Chain/mantrachain/issues/75))
- *(deps)* Bump [github.com/CosmWasm/wasmvm/v2](https://www.google.com/search?q=https://github.com/CosmWasm/wasmvm/v2) from 2.0.0 to 2.0.3
- *(deps)* Bump [github.com/cometbft/cometbft](https://github.com/cometbft/cometbft) from 0.38.10 to 0.38.11 ([\#86](https://github.com/MANTRA-Chain/mantrachain/issues/86))
- *(deps)* Bump [github.com/spf13/cast](https://github.com/spf13/cast) from 1.6.0 to 1.7.0 ([\#85](https://github.com/MANTRA-Chain/mantrachain/issues/85))
- *(deps)* Bump golang.org/x/tools from 0.22.0 to 0.24.0 ([\#84](https://github.com/MANTRA-Chain/mantrachain/issues/84))
- *(deps)* Bump github.com/cosmos/ibc-go/modules/capability ([\#83](https://github.com/MANTRA-Chain/mantrachain/issues/83))
- *(deps)* Bump [github.com/cosmos/gogoproto](https://github.com/cosmos/gogoproto) from 1.6.0 to 1.7.0 ([\#88](https://github.com/MANTRA-Chain/mantrachain/issues/88))
- *(deps)* Bump google.golang.org/grpc/cmd/protoc-gen-go-grpc ([\#91](https://github.com/MANTRA-Chain/mantrachain/issues/91))
- *(deps)* Bump cosmossdk.io/x/feegrant from 0.1.0 to 0.1.1 ([\#89](https://github.com/MANTRA-Chain/mantrachain/issues/89))
- *(deps)* Bump [github.com/grpc-ecosystem/grpc-gateway/v2](https://www.google.com/search?q=https://github.com/grpc-ecosystem/grpc-gateway/v2) ([\#98](https://github.com/MANTRA-Chain/mantrachain/issues/98))
- *(deps)* Bump cosmossdk.io/x/circuit from 0.1.0 to 0.1.1 ([\#99](https://github.com/MANTRA-Chain/mantrachain/issues/99))
- *(deps)* Bump cosmossdk.io/tools/confix from 0.1.1 to 0.1.2 ([\#101](https://github.com/MANTRA-Chain/mantrachain/issues/101))
- *(deps)* Bump [github.com/prometheus/client\_golang](https://github.com/prometheus/client_golang) ([\#97](https://github.com/MANTRA-Chain/mantrachain/issues/97))
- *(deps)* Bump cosmossdk.io/log from 1.4.0 to 1.4.1 ([\#104](https://github.com/MANTRA-Chain/mantrachain/issues/104))
- *(deps)* Bump github.com/grpc-ecosystem/grpc-gateway/v2 ([\#106](https://www.google.com/search?q=https-%3Egithub.com/MANTRA-Chain/mantrachain/issues/106))
- *(deps)* Bump [github.com/bufbuild/buf](https://github.com/bufbuild/buf) from 1.36.0 to 1.37.0 ([\#105](https://github.com/MANTRA-Chain/mantrachain/issues/105))
- *(deps)* Bump cosmossdk.io/client/v2
- *(deps)* Bump github.com/cometbft/cometbft-db from 0.11.0 to 0.14.0
- *(deps)* Bump github.com/prometheus/client\_golang
- *(deps)* Bump axios
- *(deps)* Bump github.com/skip-mev/slinky from 1.0.8 to 1.0.10 ([\#132](https://github.com/MANTRA-Chain/mantrachain/issues/132))
- *(deps)* Bump elliptic ([\#126](https://github.com/MANTRA-Chain/mantrachain/issues/126))
- *(deps)* Bump the npm\_and\_yarn group across 1 directory with 3 updates ([\#136](https://github.com/MANTRA-Chain/mantrachain/issues/136))
- *(deps)* Bump ws ([\#125](https://github.com/MANTRA-Chain/mantrachain/issues/125))
- *(deps)* Bump github.com/btcsuite/btcd ([\#127](https://github.com/MANTRA-Chain/mantrachain/issues/127))
- *(deps)* Bump [github.com/CosmWasm/wasmd](https://github.com/CosmWasm/wasmd) ([\#135](https://github.com/MANTRA-Chain/mantrachain/issues/135))
- *(deps)* Bump github.com/prometheus/client\_golang ([\#148](https://github.com/MANTRA-Chain/mantrachain/issues/148))
- *(deps)* Bump google.golang.org/grpc from 1.65.0 to 1.66.0 ([\#155](https://github.com/MANTRA-Chain/mantrachain/issues/155))
- *(deps)* Bump [github.com/docker/docker](https://github.com/docker/docker) ([\#169](https://github.com/MANTRA-Chain/mantrachain/issues/169))
- *(deps)* Bump axios ([\#170](https://github.com/MANTRA-Chain/mantrachain/issues/170))
- *(deps)* Bump github.com/cosmos/ibc-go/v8 from 8.4.0 to 8.5.0 ([\#173](https://github.com/MANTRA-Chain/mantrachain/issues/173))
- *(deps)* Bump google/osv-scanner-action from 1.7.1 to 1.8.4 ([\#180](https://github.com/MANTRA-Chain/mantrachain/issues/180))
- *(deps)* Bump dev-drprasad/delete-tag-and-release ([\#181](https://github.com/MANTRA-Chain/mantrachain/issues/181))
- *(deps)* Bump softprops/action-gh-release from 1 to 2 ([\#182](https://github.com/MANTRA-Chain/mantrachain/issues/182))
- *(deps)* Bump [github.com/cometbft/cometbft](https://github.com/cometbft/cometbft) from 0.38.11 to 0.38.12 ([\#188](https://github.com/MANTRA-Chain/mantrachain/issues/188))
- *(deps)* Bump github.com/bufbuild/buf from 1.37.0 to 1.39.0 ([\#187](https://www.google.com/search?q=https-%3Egithub.com/MANTRA-Chain/mantrachain/issues/187))
- *(deps)* Bump dev-drprasad/delete-tag-and-release from 1.0.1 to 1.1 ([\#189](https://github.com/MANTRA-Chain/mantrachain/issues/189))
- *(deps)* Bump the npm\_and\_yarn group across 1 directory with 7 updates ([\#195](https://github.com/MANTRA-Chain/mantrachain/issues/195))
- *(deps)* Bump github.com/bufbuild/buf from 1.39.0 to 1.40.0 ([\#199](https://github.com/MANTRA-Chain/mantrachain/issues/199))
- *(deps)* Bump actions/checkout from 3 to 4
- *(deps)* Bump technote-space/get-diff-action from 3 to 6
- *(deps)* Bump cosmossdk.io/store from 1.1.0 to 1.1.1
- *(deps)* Bump google.golang.org/grpc from 1.66.0 to 1.66.1 ([\#44](https://github.com/MANTRA-Chain/mantrachain/issues/44))
- *(deps)* Bump [github.com/skip-mev/slinky](https://github.com/skip-mev/slinky) in /tests/connect ([\#62](https://github.com/MANTRA-Chain/mantrachain/issues/62))
- *(deps)* Bump google/osv-scanner-action from 1.8.4 to 1.8.5 ([\#69](https://github.com/MANTRA-Chain/mantrachain/issues/69))
- *(deps)* Bump [github.com/bufbuild/buf](https://github.com/bufbuild/buf) from 1.40.1 to 1.41.0 ([\#71](https://github.com/MANTRA-Chain/mantrachain/issues/71))
- *(deps)* Bump [github.com/cosmos/ibc-go/v8](https://www.google.com/search?q=https://github.com/cosmos/ibc-go/v8) from 8.5.0 to 8.5.1 ([\#77](https://github.com/MANTRA-Chain/mantrachain/issues/77))
- *(deps)* Bump [github.com/skip-mev/slinky](https://github.com/skip-mev/slinky) in /tests/connect ([\#76](https://www.google.com/search?q=https-%3Egithub.com/MANTRA-Chain/mantrachain/issues/76))
- *(deps)* Bump google.golang.org/grpc from 1.66.1 to 1.66.2 ([\#78](https://github.com/MANTRA-Chain/mantrachain/issues/78))
- *(deps)* Bump actions/setup-python from 4 to 5 ([\#85](https://github.com/MANTRA-Chain/mantrachain/issues/85))
- *(deps)* Bump actions/checkout from 3 to 4 ([\#86](https://github.com/MANTRA-Chain/mantrachain/issues/86))
- *(deps)* Bump actions/upload-artifact from 3 to 4 ([\#87](https://www.google.com/search?q=https-%3Egithub.com/MANTRA-Chain/mantrachain/issues/87))
- *(deps)* Bump dompurify ([\#88](https://github.com/MANTRA-Chain/mantrachain/issues/88))
- *(deps)* Bump lycheeverse/lychee-action from 1.10.0 to 2.0.0 ([\#179](https://github.com/MANTRA-Chain/mantrachain/issues/179))
- *(deps)* Bump google.golang.org/protobuf ([\#177](https://github.com/MANTRA-Chain/mantrachain/issues/177))
- *(deps)* Bump lycheeverse/lychee-action from 2.0.0 to 2.0.2 ([\#188](https://github.com/MANTRA-Chain/mantrachain/issues/188))

### DOCUMENTATION

- Summarized: Added/updated READMEs, ADRs, and tx documentation for new modules (Guard, CoinFactory, etc.).
- Added adr-006 adr-007 ([\#74](https://github.com/MANTRA-Chain/mantrachain/issues/74))
- *(adr)* Update adr-006-standardise-coinfactory.md ([\#154](https://github.com/MANTRA-Chain/mantrachain/issues/154))
- *(adr)* Create adr-008-use-neutron-tokenfactory.md ([\#153](https://github.com/MANTRA-Chain/mantrachain/issues/153))

### TESTING

- Established initial e2e and unit test suites for new modules (Guard, Token, DEX, NFT, etc.).

### CI

- Summarized: Established initial CI pipelines; set up devcontainers, build profiles, and Makefiles; configured repo naming and genesis validation.
- Lint all tests ([\#113](https://github.com/MANTRA-Chain/mantrachain/issues/113))
- Add script for go-proto generation ([\#174](https://github.com/MANTRA-Chain/mantrachain/issues/174))
- Update makefile + generate docs ([\#4](https://github.com/MANTRA-Chain/mantrachain/issues/4))
- Update genesis.json ([\#50](https://github.com/MANTRA-Chain/mantrachain/issues/50))
- Assign uom denom ([\#51](https://github.com/MANTRA-Chain/mantrachain/issues/51))
- Validate genesis ([\#57](https://github.com/MANTRA-Chain/mantrachain/issues/57))
- Update genesis configuration. increment app version, adjust transaction parameters, and modify voting settings ([\#122](https://www.google.com/search?q=https-%3Egithub.com/MANTRA-Chain/mantrachain/issues/122))
- Adjust voting and inflation parameters in genesis.json ([\#145](https://github.com/MANTRA-Chain/mantrachain/issues/145))
- Update genesis ([\#178](https://github.com/MANTRA-Chain/mantrachain/issues/178))
- Add genesis for testnet and mainnet ([\#190](https://github.com/MANTRA-Chain/mantrachain/issues/190))

### OTHER

- Summarized: Performed major SDK upgrades and refactoring. Fixed swap issues and updated `Params.md`.
- Remove unnecessary whitelisting code ([\#73](https://github.com/MANTRA-Chain/mantrachain/issues/73))
- Did store code to allow lint forcetypeassert ([\#93](https://github.com/MANTRA-Chain/mantrachain/issues/93))
- Patching iavl app hash mismatch glitch ([\#95](https://github.com/MANTRA-Chain/mantrachain/issues/95))
- Remove tools ([\#109](https://github.com/MANTRA-Chain/mantrachain/issues/109))
- Bump deps ([\#28](https://github.com/MANTRA-Chain/mantrachain/issues/28))
- Use osmosis sdk ([\#37](https://github.com/MANTRA-Chain/mantrachain/issues/37))
