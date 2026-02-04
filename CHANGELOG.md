## [unreleased]

### DEPENDENCIES

- *(deps)* Bump github.com/cometbft/cometbft from 0.38.20 to 0.38.21 ([#568](https://github.com/MANTRA-Chain/mantrachain/issues/568))

### BUG-FIXES

- Return original error if no evm chain-id found ([#582](https://github.com/MANTRA-Chain/mantrachain/pull/582))

## v7.0.0

*January 26, 2026*

### FEATURES

- feat: token migration OM to MANTRA upgrade ([\#557](https://github.com/MANTRA-Chain/mantrachain/issues/557))

### BUG-FIXES

- Update version extraction regex for SDK version detection ([\#487](https://github.com/MANTRA-Chain/mantrachain/issues/487))

### DEPENDENCIES

- *(deps)* Bump github.com/consensys/gnark-crypto from 0.18.0 to 0.18.1 ([#519](https://github.com/MANTRA-Chain/mantrachain/issues/519))
- *(deps)* Bump github.com/cometbft/cometbft from 0.38.19 to 0.38.20 ([#550](https://github.com/MANTRA-Chain/mantrachain/issues/550))
- *(deps)* Bump github.com/cometbft/cometbft from 0.38.20 to 0.38.21 (backport [#568](https://github.com/MANTRA-Chain/mantrachain/issues/568)) ([#570](https://github.com/MANTRA-Chain/mantrachain/issues/570))

### DOCUMENTATION

- Add CHANGELOG.md ([\#516](https://github.com/MANTRA-Chain/mantrachain/issues/516))

### CI

- Updated Nix to version 25.05 and optimized e2e test execution workflow. ([\#525](https://github.com/MANTRA-Chain/mantrachain/issues/525))

### OTHER

- Add template for v7 upgrade ([\#486](https://github.com/MANTRA-Chain/mantrachain/issues/486))
- Remove depinject boilerplate and proto gen ([\#527](https://github.com/MANTRA-Chain/mantrachain/issues/527))

## v6.1.2

*November 14, 2025*

### BUG-FIXES

- Avoid crash on nil evmCoinInfo with grpc only mode ([#520](https://github.com/MANTRA-Chain/mantrachain/issues/520))

## v6.1.1

*October 30, 2025*

### BUG-FIXES

- Historical EvmCoinInfo don't exists ([\#503](https://github.com/MANTRA-Chain/mantrachain/issues/503)) ([\#505](https://github.com/MANTRA-Chain/mantrachain/issues/505))

### CI

- *(release)* Add gpg signing on the binaries checksum ([\#457](https://github.com/MANTRA-Chain/mantrachain/issues/457)) ([\#507](https://github.com/MANTRA-Chain/mantrachain/issues/507))

## v6.1.0

*October 27, 2025*

### BUG-FIXES

- Evm fix not merged ([\#497](https://github.com/MANTRA-Chain/mantrachain/issues/497)) ([\#499](https://github.com/MANTRA-Chain/mantrachain/issues/499))
- V6.1.0 upgrade handler not added ([\#500](https://github.com/MANTRA-Chain/mantrachain/issues/500))

## v6.0.1

*October 30, 2025*

### BUG-FIXES

- Historical EvmCoinInfo don't exists ([\#503](https://github.com/MANTRA-Chain/mantrachain/issues/503)) ([\#504](https://github.com/MANTRA-Chain/mantrachain/issues/504))
- Remove duplicate backport rule for release/v6.0.x branch

### CI

- *(release)* Add gpg signing on the binaries checksum ([\#457](https://github.com/MANTRA-Chain/mantrachain/issues/457)) ([\#506](https://github.com/MANTRA-Chain/mantrachain/issues/506))

## v6.0.0

*October 22, 2025*

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

## v5.0.2

*October 15, 2025*

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

## v5.0.0

*September 13, 2025*

### FEATURES

- Ibc-v10 ([\#321](https://github.com/MANTRA-Chain/mantrachain/issues/321))
- Evm module with v5 upgrade ([\#345](https://github.com/MANTRA-Chain/mantrachain/issues/345))
- Block mint burn for disabled tokenfactory coins ([\#359](https://github.com/MANTRA-Chain/mantrachain/issues/359))
- V5.0.0 evm enhancements + bug fixes ([\#364](https://github.com/MANTRA-Chain/mantrachain/issues/364))

## v4.0.2

*September 8, 2025*

### FEATURES

- Added x/circuit api ([\#305](https://github.com/MANTRA-Chain/mantrachain/issues/305)) ([\#306](https://github.com/MANTRA-Chain/mantrachain/issues/306))

### BUG-FIXES

- None ([\#437](https://github.com/MANTRA-Chain/mantrachain/issues/437))
- None ([\#433](https://github.com/MANTRA-Chain/mantrachain/issues/433)) ([\#438](https://github.com/MANTRA-Chain/mantrachain/issues/438))
- Behavioral change between proxygasmeter and standalone gas meter ([\#439](https://github.com/MANTRA-Chain/mantrachain/issues/439))

### CI

- Bump cosmos-sdk fork from v0.50.13-v3-mantra-1 to v0.50.14-v4-mantra-1 ([\#368](https://github.com/MANTRA-Chain/mantrachain/issues/368))

## v4.0.1

*April 3, 2025*

### FEATURES

- Pump go in release.yml ([\#295](https://github.com/MANTRA-Chain/mantrachain/issues/295)) ([\#296](https://github.com/MANTRA-Chain/mantrachain/issues/296))

### BUG-FIXES

- Add openapi route ([\#299](https://github.com/MANTRA-Chain/mantrachain/issues/299)) ([\#300](https://github.com/MANTRA-Chain/mantrachain/issues/300))
- Remove broken links ([\#304](https://github.com/MANTRA-Chain/mantrachain/issues/304))
- Change okx websocket url ([\#310](https://github.com/MANTRA-Chain/mantrachain/issues/310))
- Update heighliner tag to v1.7.5 ([\#325](https://github.com/MANTRA-Chain/mantrachain/issues/325))
- Update workflow configurations for CodeQL and connect-test ([\#327](https://github.com/MANTRA-Chain/mantrachain/issues/327))

## v4.0.0

*March 25, 2025*

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

## v3.0.3

*March 21, 2025*

### CI

- Bump cosmos-sdk fork from v0.50.12-v3-mantra-1 to v0.50.13-v3-mantra-1 ([\#286](https://github.com/MANTRA-Chain/mantrachain/issues/286))

## v3.0.2

*March 18, 2025*

### CI

- Bump cosmos-sdk fork from v0.50.12-v2-mantra-1 to v0.50.12-v3-mantra-1 ([\#281](https://github.com/MANTRA-Chain/mantrachain/issues/281))

## v3.0.1

*March 14, 2025*

### DEPENDENCIES

- *(deps)* Bump [github.com/cosmos/ibc-go/v8](https://www.google.com/search?q=https://github.com/cosmos/ibc-go/v8) from v8.6.1 to v8.7.0 ([\#279](https://github.com/MANTRA-Chain/mantrachain/issues/279))

## v3.0.0

*March 4, 2025*

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

## v2.0.2

*February 21, 2025*

### DEPENDENCIES

- *(deps)* Bump cosmos-sdk MANTRA fork to v0.50.12 ([\#268](https://github.com/MANTRA-Chain/mantrachain/issues/268))

## v2.0.0

*February 13, 2025*

### FEATURES

- Add e2e tests ([\#252](https://github.com/MANTRA-Chain/mantrachain/issues/252))

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

## v1.0.3

*December 11, 2024*

### CI

- Upgrade cosmossdk.io/math to v1.4.0 ([\#239](https://github.com/MANTRA-Chain/mantrachain/issues/239))

## v1.0.2

*November 8, 2024*

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

## v1.0.0

*October 23, 2024*

Initial Release!
