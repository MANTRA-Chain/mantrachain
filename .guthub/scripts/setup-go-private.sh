#!/bin/bash

cp .github/files/.gitconfig ~/.gitconfig

mkdir -p ~/.ssh

cp .github/files/config ~/.ssh/config

echo "$COSMOS_SDK" > ~/.ssh/cosmos-sdk
echo "$IBC_GO" > ~/.ssh/ibc-go

chmod 600 ~/.ssh/cosmos-sdk
chmod 600 ~/.ssh/ibc-go