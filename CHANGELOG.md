## [unreleased]

### 🐛 Bug Fixes

- Update version extraction regex for SDK version detection (#487)

### 💼 Other

- Miss coin info in connect test (#484)
- Add template for v7 upgrade (#486)
- Upstream evm v0.5.0 not used (#491)
- Evm fix not merged (#497)
- Historical EvmCoinInfo don't exists (#503)

### ⚙️ Dependencies

- *(deps)* Update cosmos-sdk fork from v0.53.3 to v0.53.4 (#488)

### ⚙️ Miscellaneous Tasks

- *(release)* Add gpg signing on the binaries checksum (#457)
## [6.0.0-rc0] - 2025-10-17

### 🚀 Features

- Deprecate x/group (#459)
- Additional blacklist for wasm msgserver (#464)
- Add comet query commands (#475)

### 🐛 Bug Fixes

- Error ignored by scoped error variables (#452)
- Change chain id back to default after sig verify (#453)

### 💼 Other

- Add template for v6 upgrade (#461)
- Owner is not aligned in tokenfactory token_pairs for ibc transfer (#460)
- Authz module is risky (#462)
- Static precompiles are not enabled (#465)
- Not set disable list of wasm breaker in upgrade (#466)
- Cometbft is outdated (#477)
- Evm is outdated (#478)
- Upstream evm is not integrated (#468)

### ⚙️ Dependencies

- *(deps)* Bump github.com/ulikunitz/xz from v0.5.11 to v0.5.15 (#450)
- *(deps)* Bump github.com/Zondax/ledger-go from 0.14.3 to 1.0.1 (#458)

### ⚙️ Miscellaneous Tasks

- Read chain ID from app.toml before genesis (#454)
- *(release)* Add github attestation on the binaries (#456)
## [5.0.0] - 2025-09-13

### 🚀 Features

- Add grpc querier for wasm query plugins (#289)
- Add mergify.yml (#294)
- Pump go in release.yml (#295)
- Added x/circuit api (#305)
- Ibc-v10 (#321)
- Evm module with v5 upgrade (#345)
- Block mint burn for disabled tokenfactory coins (#359)
- V5.0.0 evm enhancements + bug fixes (#364)

### 🐛 Bug Fixes

- Add openapi route (#299)
- Remove broken links (#304)
- Change okx websocket url (#310)
- Update heighliner tag to v1.7.5 (#325)
- Update workflow configurations for CodeQL and connect-test (#327)

### 💼 Other

- Add template for v5 upgrade (#283)
- *(golangci-lint)* Migrate to v2 config and apply lint to repo (#358)
- Allow Unprotected Transactions (#372)

### ⚙️ Miscellaneous Tasks

- Update connect to v2.3.0-v4-mantra-1 (#293)
- Add license badge (#317)
- Remove networks folder (#323)
- Remove registerdenoms (#349)
- Update backport rules for release/v5.0.x branch (#445)
## [4.0.0-rc1] - 2025-03-21

### 🚀 Features

- Interchain accounts (#278)
- Whitelist getprice query for wasm (#284)

### 💼 Other

- Add template for v4 upgrade (#277)

### ⚙️ Dependencies

- *(deps)* Bump github.com/cosmos/ibc-go/v8 from v8.6.1 to v8.7.0 (#280)

### ⚙️ Miscellaneous Tasks

- Bump cosmos-sdk fork from v0.50.12-v2-mantra-1 to v0.50.12-v3-mantra-1 (#282)
- Bump cosmos-sdk fork from v0.50.12-v3-mantra-1 to v0.50.13-v3-mantra-1 (#287)
## [3.0.0] - 2025-03-04

### 🚀 Features

- Increase max-wasm-size to 3mb (#266)
- Sanction module (#267)

### ⚙️ Dependencies

- *(deps)* Bump cosmos-sdk MANTRA fork to v0.50.12 (#269)
- *(deps)* Bump lycheeverse/lychee-action from 2.1.0 to 2.3.0 (#258)
- *(deps)* Bump github.com/cosmos/ibc-go/v8 from v8.5.2 to v8.6.1 (#272)

### ⚙️ Miscellaneous Tasks

- Upgrade cometbft to v0.38.17 + wasmvm to v2.2.2 (#265)
- Use new connect fork that is based off v2.3.0 (#270)
## [2.0.0] - 2025-02-13

### 🚀 Features

- Upgrade to v2
## [2.0.0-rc2] - 2025-01-06

### 🚀 Features

- Add e2e tests (#252)

### 🐛 Bug Fixes

- Register uom as DefaultBondDenom (#195)

### ⚙️ Dependencies

- *(deps)* Bump elliptic (#198)
- *(deps)* Bump github.com/bufbuild/buf from 1.45.0 to 1.46.0 (#197)
- *(deps)* Bump github.com/grpc-ecosystem/grpc-gateway/v2 (#206)
- *(deps)* Bump github.com/cometbft/cometbft from 0.38.12 to 0.38.13 (#196)
- *(deps)* Bump github.com/prometheus/client_golang (#193)
- *(deps)* Bump github.com/skip-mev/connect/v2 in /tests/connect (#209)
- *(deps)* Bump github.com/skip-mev/connect/v2 from 2.1.0 to 2.1.2 (#208)
- *(deps)* Bump github.com/cometbft/cometbft in the go_modules group (#215)
- *(deps)* Bump cosmossdk.io/log from 1.4.1 to 1.5.0 (#220)
- *(deps)* Bump github.com/cosmos/ibc-go/v8 from 8.5.1 to 8.5.2 (#212)
- *(deps)* Bump github.com/cometbft/cometbft (#214)
- *(deps)* Bump cosmossdk.io/math (#227)
- *(deps)* Bump cosmossdk.io/math in the go_modules group (#228)
- *(deps)* Bump github.com/stretchr/testify in /tests/connect (#229)
- *(deps)* Bump google.golang.org/grpc from 1.67.1 to 1.68.1 (#235)
- *(deps)* Bump golang.org/x/tools from 0.26.0 to 0.28.0 (#234)
- *(deps)* Bump codecov/codecov-action from 4 to 5 (#223)
- *(deps)* Bump lycheeverse/lychee-action from 2.0.2 to 2.1.0 (#218)
- *(deps)* Bump cosmossdk.io/depinject from 1.0.0 to 1.1.0 (#216)
- *(deps)* Bump github.com/cosmos/cosmos-db from 1.0.2 to 1.1.0 (#237)

### ⚙️ Miscellaneous Tasks

- Remove use of osv scanner (#204)
- Remove unnecessary ci jobs (#205)
- Remove legacy ADRs (#207)
- Adjust cron schedule for CodeQL (#230)
- Add sonar project properties (#232)
- Use new cosmos-sdk fork that is based off v0.50.10 (#238)
- Bump cosmos-sdk MANTRA fork to v0.50.11 (#250)
## [1.0.0-rc3] - 2024-10-15

### 🚀 Features

- Wasm api (#191)

### 🐛 Bug Fixes

- Simulate gas calculations (#187)
- Downgrade cometbft-db due to forward compatibility issue (#189)

### 💼 Other

- No license (#192)

### ⚙️ Dependencies

- *(deps)* Bump github.com/bufbuild/buf from 1.44.0 to 1.45.0 (#183)
- *(deps)* Bump lycheeverse/lychee-action from 2.0.0 to 2.0.2 (#188)

### ⚙️ Miscellaneous Tasks

- Add genesis for testnet and mainnet (#190)
## [1.0.0-rc2] - 2024-10-09

### ⚙️ Dependencies

- *(deps)* Bump github.com/bufbuild/buf from 1.43.0 to 1.44.0 (#172)
- *(deps)* Bump lycheeverse/lychee-action from 1.10.0 to 2.0.0 (#179)
- *(deps)* Bump google.golang.org/protobuf (#177)

### ⚙️ Miscellaneous Tasks

- Update genesis (#178)
## [1.0.0-rc1] - 2024-10-07

### 🚀 Features

- Add ability to query the tokenfactory module parameters and denom info (#168)
- Update Genesis for RC1 and Tokenomics Gov Proposal alignment (#173)
- Add default denom resolver (#175)

### ⚙️ Dependencies

- *(deps)* Bump google/osv-scanner-action from 1.8.5 to 1.9.0 (#170)
- *(deps)* Bump golang.org/x/tools from 0.25.0 to 0.26.0 (#174)
## [1.0.0-alpha.10] - 2024-10-02

### 🐛 Bug Fixes

- Discontinue the use of depinject (#160)

### ⚙️ Dependencies

- *(deps)* Bump github.com/skip-mev/connect/v2 in /tests/connect (#163)
- *(deps)* Bump github.com/bufbuild/buf from 1.42.0 to 1.43.0 (#165)
- *(deps)* Bump google.golang.org/grpc from 1.67.0 to 1.67.1 (#166)
- *(deps)* Bump github.com/skip-mev/connect/v2 from 2.0.1 to 2.1.0 (#164)
## [1.0.0-alpha.8] - 2024-09-27

### 🐛 Bug Fixes

- Remove the use of app_config authority for tax (#135)
- Error log in connect test (#143)

### ⚙️ Dependencies

- *(deps)* Bump github.com/strangelove-ventures/interchaintest/v8 (#128)
- *(deps)* Bump google.golang.org/grpc from 1.66.2 to 1.67.0 (#127)
- *(deps)* Bump github.com/btcsuite/btcd from 0.22.0-beta to 0.24.0 in /tests/connect (#8)
- *(deps)* Bump cosmossdk.io/api from 0.7.5 to 0.7.6 (#134)
- *(deps)* Bump github.com/CosmWasm/wasmvm/v2 from 2.1.2 to 2.1.3 (#133)

### ⚙️ Miscellaneous Tasks

- Adjust voting and inflation parameters in genesis.json (#145)
## [1.0.0-alpha.7] - 2024-09-20

### 🐛 Bug Fixes

- Align min_base_gas_price for fee market (#123)
## [1.0.0-alpha.6] - 2024-09-20

### ⚙️ Dependencies

- *(deps)* Bump github.com/bufbuild/buf from 1.41.0 to 1.42.0 (#119)
- *(deps)* Bump github.com/skip-mev/connect/v2 (#120)
- *(deps)* Bump github.com/skip-mev/connect/v2 in /tests/connect (#121)
## [1.0.0-alpha.5] - 2024-09-20

### 🐛 Bug Fixes

- Use github.ref_name instead (#95)

### ⚙️ Dependencies

- *(deps)* Bump cosmossdk.io/client/v2 (#109)

### ⚙️ Miscellaneous Tasks

- Update genesis configuration. increment app version, adjust transaction parameters, and modify voting settings (#122)
## [1.0.0-alpha.4] - 2024-09-17

### 🚀 Features

- Add bridge module
- Add update for buy nft to works with cw20
- Add query bridge
- Add init sample state script
- [**breaking**] Add update nft stake
- Add update vault
- Add get last epoch block - vault query
- Add marketplace nft royalties
- Add evm
- Add guard module
- Add create coin factory module
- Scaffold new chain
- Add custom modules mantra
- Add upgrade cosmos-sdk v0.47-rc2
- Add crescent modules
- Update account privileges
- Add required priveleges
- Add lock and seize to coins
- Add guard ante authz decorator
- Add authz grant revoke generic batch
- Add upgrade cosmos-sdk v0.47.0
- Add update dex privileges check
- Update account privileges
- Add required priveleges
- Add lock and seize to coins
- Add guard ante authz decorator
- Add authz grant revoke generic batch
- Add upgrade cosmos-sdk v0.47.0
- Add update dex privileges check
- Add downgrade sdk v0.45.10 because of ccv
- Add interchain ccv v1.1.0 cosmos-sdk v0.45.10
- Add coinfactory bank query methods
- Add did module
- Add update guard soul-bond nft image by index
- Add/remove dids for guard soul-bond nfts
- Add liquidity module query pairs-by-denoms
- Scripts update
- Add marketmaker dex module
- Add liquidfarming dex module
- Add guard checks to liquidfarming and marketmaker module
- Scaffold new chain with ignite cli v0.27.1, cosmos-sdkv 0.47.3, ibc-go v7.1.0
- Add dex modules v4.2.0
- Add update cosmos-sdk and ibc-go with mantra forks
- Adjust dex + custom modules
- Replaced mantrachain with github.com/MANTRA-Finance/mantrachain as name of the project
- Added testdata_pulsar
- Removed testdata_pulsar
- Add liquidity module query pairs-by-denoms
- Updated tests e2e, updated ignite config
- Add update ibc-go to support guard
- Add include ante handler
- Add pay fees to admin acc
- Pay with alternate token
- Add remove validator rewards
- Add impl. cosmwasm module
- Add coinfactory cosmwasm bindings
- Bring back lost file
- Bring back lost change
- Remove liquid farming module as not needed
- Add authz required privileges
- Add collect swap fees to an pair swap fees escrow address
- Add liquidity module hooks
- Add rewards module
- Add snapshot to rewards module
- Add provider to rewards module
- Updated comment
- Add collect swap fees to an pair swap fees escrow address
- Add rewards module
- Update e2e tests
- Fix github pipelines
- Add batch update restricted collection nft image
- Add swap fee rate and pair creator swap fee ration params on pair creation
- Add txfees gas estimation query
- Add airdrop module
- Update native denom to uom
- Coinfactory fix issues and add mintTo and burnFrom flags
- Add bridge module
- Upgrade to cosmos-sdk v0.50
- Update ante handler
- Add the v1 - v2 migration logic
- Add cosmwasm
- Update guard to allow transfers when required coin privileges is not being set (#53)
- Allow to delete unused image tag
- Integrate slinky (#66)
- Feemarket with dependency injection (#177)
- Added oracle genesis (#47)
- Ibc rate limit (#67)
- Visualize genesis (#75)

### 🐛 Bug Fixes

- Change proto queries
- Update bridge txs namespace
- Add update staked amount param on buy nft
- Add burn when stake on remote chain
- Fix query all nft collections
- Update buy nft handler
- Fix marketplace/vault staking
- Add set staked fix
- Add fix ustake delegation reward
- Vault refactoring
- Add nftId to buyNft event
- Add update set nft stake
- Add fix epoch shares
- Add update modules ValidateBasic handlers
- Add update nft stake query amounts to string
- Add update nft stake staked
- Add fix start epoch expects init epoch
- Add update cr20 amounts types from int64 to string
- Update init test script
- Add update init script
- Update scripts
- Update init scripts
- Add update coin factory msgs
- Add fix coinfactory module query cli
- Add guard update
- Upgrade ibc-go; fix token soul bonded nfts
- Fix some todos in code
- Add guard module fixes and update the docs
- Fix guard module check transactions
- Fix token not retrieve nft id
- Remove default denom creation free from coinfactory module
- Add fix query nft approved response
- Add update guard
- Fix some todos in code
- Add guard module fixes and update the docs
- Fix guard module check transactions
- Fix token not retrieve nft id
- Remove default denom creation free from coinfactory module
- Add fix query nft approved response
- Add update guard
- Add fix coinfactory seize coins
- Add fix ccv
- Add fix tranfer coins guard check
- Add fix guard error
- Add coinfactory seize fix and ccv init scripts
- Add fix transfer coins guard whitelist addresses
- Add guard grant fix and remove unused dex modules
- Update makefile
- Swagger ui
- Add update init scripts and guard module e2e tests
- Add fix guard transfer coins
- Add fix guard transfer coins + tests
- Add fix guard auth restrictions
- Add coinfactory force transfer fix + tests
- Add update guard locks and whitelist addresses for coins transfers
- Add fix whitelist pair escrow address on order
- Remove token did authentication and fix e2e tests
- Add fix compare privileges
- Add update make localnet script
- Add fix start consumer
- Add check second coin privileges on swap
- Fig ignite chain serve ignite ver0.27.0
- Add fix throw when gas fee is zero
- Add update coinfactory txs permissions
- Add fixes coinfactory and cosmwasm settings
- Update coinfactory mint and add tutorial docs
- Add did module import/export genesis
- Add update create denom fee receiver
- Fix did import genesis
- Add token module fixes + unit tests
- Add token module fixes + unit tests
- Add token module fixes + unit tests
- Add fix guard module check privileges
- Add fix guard module check privileges
- Add fix guard module check privileges
- Add fix liquidity module
- Fix liquidity module issue and fix e2e tests
- Fix guard allow transfer pool coins and fix tests
- Fix guard
- Added version ldflag
- Update tx fees module to have version and add test
- Update authz for coinfactory
- Add fix account/required privileges cli
- Tmp remove rewards module
- Fix guard swap privileges issue
- Add update guard swap privileges issue
- Update rewards module to not keep balances
- Rename MANTRA-Finance to AumegaChain for the rewards module
- Add proto fix
- Add fix x/rewards err message
- Add rewards fix
- Fix proto queries path
- Add fix x/rewards err message
- Add rewards fix
- Fix proto queries path
- Add extend finish swap events
- Fix swap refund coin nil
- Add update liquidity swap fee collecting logic
- Fix
- Fiw txfees swap amount
- Fix swap
- Add swap fees fix
- Add swap fees fix
- Fix swap fees
- Fix swap fees snapshot remaining coins
- Add fix snapshot pool unused var
- Fix swap fee refund issue
- Fix txfees swap amount
- Fix swap fees events params
- Add batch update restricted nft image validation checks
- Replace swap finish order returning errors with logging
- Update logger messages
- Restore swap finish order return errors on refund
- Fix update swap fees for a liquidity pair
- Fix swap fee calculation on swap
- Add fix gas estimation url
- Fix txfees gas estimation
- Fix swap fee rate when query for liqudity pairs
- Fix swap fees issue
- Update denom
- Fix getting balances on airdrop campaign delete
- Update airdrop type urls
- Add fix swap fees
- Add fix swap fees
- Add fix swap fees
- Fix swap fees rounding error issue
- Update rewards purge/distribution cycle
- Refactor guard coin transfer updates and fix e2e tests
- Split WhitelistTransferAccAddresses to two functions
- Fix guard multi coin transfer issue
- Fix guard unit tests
- Fix guard multi coin transfer issue
- Fix genesis proto
- Update v1 - v2 migration logic
- Fix dublicate types of upgrate proto
- Add fix proto files types
- Add fix migrate modules params
- Fix comments and binary version
- Try fix consensus params issue
- Fix the upgrade
- Fix ibc transfer keeper init
- Fix migrating params
- Fix rewards module upgrade pools coin supply issue
- Fix issue with migrate or purge not distributed snapshots
- Fix snapshots purge
- Fix rewards module migration
- Revert the rewards snapshots purge fix
- Small update on farming genesis validate
- Remove duplicate code
- Fix make cosmwasm version and add mock deps
- Fix coinfactory types tests
- Fix swap fees empty coins
- Fix swap fees empty coins
- Add guard fix causing chain halt
- Update ci for build and docker
- Disable codeql and gosec until repo is public and bump go version for unit test ci
- Update build and release ci
- Fix musl image build and use larger runner
- *(lint)* Commits from previous main (#63)
- E2e tests jest global setup (#65)
- Set local setup default swap fee rate param to zero (#72)
- *(lint)* Revert lint-install to previous script
- The name of static lib changes
- Cli get params (#133)
- Generated swagger + openapi
- Lint
- Simplify TokenFactory depinject (#36)
- Use goreleaser for release (#41)
- Enable vote extension (#48)
- Update denoms and timings in genesis file (#52)
- Missing rest api (#64)
- Harmonize block gas (#91)
- Allow docker and static linked binary built for PR ci and fix no… (#93)

### 💼 Other

- Upgrade sdk
- Upgrade sdk - fix issues
- Upgrade sdk
- Update privileges
- Add guard module token txs authz restrictions
- Update privileges
- Add guard module token txs authz restrictions
- Add guard module tests
- Add guard module tests
- Add guard module tests
- Add token module tests
- Add token module tests
- Add guard module tests
- Add guard module tests
- Add guard module tests
- Add e2e guard tests
- Add e2e guard tests
- Add e2e guard tests
- Restore e2e tests
- Add and adjust unit tests
- Add txfees module
- Params.md
- Add swap fees fix
- Add tests and fix guard
- Migrating unit tests
- Fixing tests
- Fix unit tests
- Add and fix liquidity module unit tests
- Fix liquidity tests and add lpfarm tests
- Patching iavl app hash mismatch glitch (#95)
- Use osmosis sdk (#37)

### 🚜 Refactor

- Some guard coin transfers updates
- Refactor WhitelistTransferAccAddresses execute
- Remove unnecessary whitelisting code (#73)
- Did store code to allow lint forcetypeassert (#93)
- Remove tools (#109)

### 📚 Documentation

- Add guard readme
- Add guard readme
- Add guard module README.md
- Update docs
- Update docs
- Update params docs
- Add coinfactory txs docs
- Update docs
- Update docs
- Add update txs flows docs
- Fix typo
- Update docs
- Update token and tx fees modules docs
- Update coinfactory tx docs
- Add update txs docs
- Update docs
- Update openapi spec
- Update openapi.yml
- Added adr-006 adr-007 (#74)
- *(adr)* Update adr-006-standardise-coinfactory.md (#154)
- *(adr)* Create adr-008-use-neutron-tokenfactory.md (#153)

### 🧪 Testing

- Add dex liquidity module tests
- Add dex lpfarm module tests
- Add coinfactory module tests
- Add nft module test
- Add e2e tests setup
- Add fix guard module e2e tests
- Add guard module e2e tests
- Fix guard module e2e tests
- Add token module e2e tests
- Add guard module e2e tests
- Add guard module e2e tests
- Fix tests
- Fix unit tests
- Update liquidity module e2e tests
- Add e2e tests
- Fix dex tests
- Fix e2e tests
- Fix unit tests
- Add test cover command
- Add token module tests
- Add e2e cosm wasm test
- Add token module unit tests
- Add token module unit tests
- Add token module unit tests
- Add guard module unit tests
- Add guard module unit tests
- Update mantrachain-sdk for the e2e tests
- Add e2e tests
- Add e2e test for create denom account privileges
- Fix liquidity test
- Fix utnit tests
- Fix tests

### ⚙️ Dependencies

- *(deps)* Bump github.com/bufbuild/buf from 1.30.0 to 1.36.0 (#56)
- *(deps)* Bump actions/setup-go from 3 to 5 (#35)
- *(deps)* Bump github.com/spf13/viper from 1.18.2 to 1.19.0 (#48)
- *(deps)* Bump cosmossdk.io/x/evidence from 0.1.0 to 0.1.1 (#51)
- *(deps)* Bump github.com/cosmos/cosmos-sdk from 0.50.7 to 0.50.9 (#62)
- *(deps)* Bump github.com/cosmos/gogoproto from 1.5.0 to 1.6.0 (#78)
- *(deps)* Bump cosmossdk.io/x/nft from 0.1.0 to 0.1.1 (#77)
- *(deps)* Bump cosmossdk.io/x/upgrade from 0.1.3 to 0.1.4 (#76)
- *(deps)* Bump cosmossdk.io/log from 1.3.1 to 1.4.0 (#75)
- *(deps)* Bump github.com/CosmWasm/wasmvm/v2 from 2.0.0 to 2.0.3
- *(deps)* Bump github.com/cometbft/cometbft from 0.38.10 to 0.38.11 (#86)
- *(deps)* Bump github.com/spf13/cast from 1.6.0 to 1.7.0 (#85)
- *(deps)* Bump golang.org/x/tools from 0.22.0 to 0.24.0 (#84)
- *(deps)* Bump github.com/cosmos/ibc-go/modules/capability (#83)
- *(deps)* Bump github.com/cosmos/gogoproto from 1.6.0 to 1.7.0 (#88)
- *(deps)* Bump google.golang.org/grpc/cmd/protoc-gen-go-grpc (#91)
- *(deps)* Bump cosmossdk.io/x/feegrant from 0.1.0 to 0.1.1 (#89)
- *(deps)* Bump github.com/grpc-ecosystem/grpc-gateway/v2 (#98)
- *(deps)* Bump cosmossdk.io/x/circuit from 0.1.0 to 0.1.1 (#99)
- *(deps)* Bump cosmossdk.io/tools/confix from 0.1.1 to 0.1.2 (#101)
- *(deps)* Bump github.com/prometheus/client_golang (#97)
- *(deps)* Bump cosmossdk.io/log from 1.4.0 to 1.4.1 (#104)
- *(deps)* Bump github.com/grpc-ecosystem/grpc-gateway/v2 (#106)
- *(deps)* Bump github.com/bufbuild/buf from 1.36.0 to 1.37.0 (#105)
- *(deps)* Bump cosmossdk.io/client/v2
- *(deps)* Bump github.com/cometbft/cometbft-db from 0.11.0 to 0.14.0
- *(deps)* Bump github.com/prometheus/client_golang
- *(deps)* Bump axios
- *(deps)* Bump github.com/skip-mev/slinky from 1.0.8 to 1.0.10 (#132)
- *(deps)* Bump elliptic (#126)
- *(deps)* Bump the npm_and_yarn group across 1 directory with 3 updates (#136)
- *(deps)* Bump ws (#125)
- *(deps)* Bump github.com/btcsuite/btcd (#127)
- *(deps)* Bump github.com/CosmWasm/wasmd (#135)
- *(deps)* Bump github.com/prometheus/client_golang (#148)
- *(deps)* Bump google.golang.org/grpc from 1.65.0 to 1.66.0 (#155)
- *(deps)* Bump github.com/docker/docker (#169)
- *(deps)* Bump axios (#170)
- *(deps)* Bump github.com/cosmos/ibc-go/v8 from 8.4.0 to 8.5.0 (#173)
- *(deps)* Bump google/osv-scanner-action from 1.7.1 to 1.8.4 (#180)
- *(deps)* Bump dev-drprasad/delete-tag-and-release (#181)
- *(deps)* Bump softprops/action-gh-release from 1 to 2 (#182)
- *(deps)* Bump github.com/cometbft/cometbft from 0.38.11 to 0.38.12 (#188)
- *(deps)* Bump github.com/bufbuild/buf from 1.37.0 to 1.39.0 (#187)
- *(deps)* Bump dev-drprasad/delete-tag-and-release from 1.0.1 to 1.1 (#189)
- *(deps)* Bump the npm_and_yarn group across 1 directory with 7 updates (#195)
- *(deps)* Bump github.com/bufbuild/buf from 1.39.0 to 1.40.0 (#199)
- *(deps)* Bump actions/checkout from 3 to 4
- *(deps)* Bump technote-space/get-diff-action from 3 to 6
- *(deps)* Bump cosmossdk.io/store from 1.1.0 to 1.1.1
- *(deps)* Bump github.com/cosmos/interchain-security/v5 (#30)
- *(deps)* Bump github.com/bufbuild/buf from 1.40.0 to 1.40.1 (#31)
- *(deps)* Bump golang.org/x/tools from 0.24.0 to 0.25.0 (#39)
- *(deps)* Bump google.golang.org/grpc from 1.66.0 to 1.66.1 (#44)
- *(deps)* Bump github.com/skip-mev/slinky in /tests/connect (#62)
- *(deps)* Bump google/osv-scanner-action from 1.8.4 to 1.8.5 (#69)
- *(deps)* Bump github.com/bufbuild/buf from 1.40.1 to 1.41.0 (#71)
- *(deps)* Bump github.com/cosmos/ibc-go/v8 from 8.5.0 to 8.5.1 (#77)
- *(deps)* Bump github.com/skip-mev/slinky in /tests/connect (#76)
- *(deps)* Bump google.golang.org/grpc from 1.66.1 to 1.66.2 (#78)
- *(deps)* Bump actions/setup-python from 4 to 5 (#85)
- *(deps)* Bump actions/checkout from 3 to 4 (#86)
- *(deps)* Bump actions/upload-artifact from 3 to 4 (#87)
- *(deps)* Bump dompurify (#88)

### ⚙️ Miscellaneous Tasks

- Add update LimeChain to MANTRA-FINANCE in code
- - add DEX scripts
- Update mantra deps
- Add update LimeChain to MANTRA-FINANCE in code
- - add DEX scripts
- Add devcontainers support and Arm build profile in Makefile
- Add devcontainers support and Arm build profile in Makefile
- Remove commented dep
- Some formatting of the rewards snapshots purge
- Fix ci as no more private gomodule
- Fix mdlint
- Don't lint on windows
- Tidy go.mod
- Update openapi.yml
- Update openapi.yml
- Update buf.lock
- Lint all tests (#113)
- Add script for go-proto generation (#174)
- Update makefile + generate docs (#4)
- Simplify proto-download-deps
- Rename all MANTRA-Finance to MANTRA-Chain
- Update makefile + generate-docs + generate proto
- Update genesis.json (#50)
- Assign uom denom (#51)
- Validate genesis (#57)

### 🛡️ Security

- Bump deps (#28)
