#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/nft/types/expected_keepers.go -package testutil -destination x/nft/testutil/expected_keepers_mocks.go
