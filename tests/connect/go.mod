module github.com/MANTRA-Chain/mantrachain/tests/connect

go 1.23

toolchain go1.23.1

require (
	github.com/icza/dyno v0.0.0-20230330125955-09f820a8d9c0
	github.com/stretchr/testify v1.9.0
)

replace (
	github.com/ChainSafe/go-schnorrkel => github.com/ChainSafe/go-schnorrkel v0.0.0-20200405005733-88cbf1b4c40d
	github.com/ChainSafe/go-schnorrkel/1 => github.com/ChainSafe/go-schnorrkel v1.0.0
	github.com/cosmos/cosmos-sdk => github.com/MANTRA-Chain/cosmos-sdk v0.50.8-0.20240920090924-e76f4a426195
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/vedhavyas/go-subkey => github.com/strangelove-ventures/go-subkey v1.0.7
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
