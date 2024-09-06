#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto

proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package $file &>/dev/null; then
      echo "Generating gogo proto code for $file"
      buf generate $file --template buf.gen.gogo.yaml
    fi
  done
done

# move proto files to the right places
#
# Note: Proto files are suffixed with the current binary version.
echo "Copying proto files to the right places"
cp -r github.com/MANTRA-Chain/mantrachain/x/* ../x

# cleanup
echo "Cleaning up"
rm -rf github.com

cd ..

# TODO: Uncomment once ORM/Pulsar support is needed.
#
# Ref: https://github.com/osmosis-labs/osmosis/pull/1589
set -eo pipefail

protoc_install_gopulsar() {
  go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
}

protoc_install_gopulsar

echo "Cleaning API directory"
(
  cd api
  find ./ -type f \( -iname \*.pulsar.go -o -iname \*.pb.go -o -iname \*.cosmos_orm.go -o -iname \*.pb.gw.go \) -delete
  find . -empty -type d -delete
  cd ..
)

echo "Generating API module"
(
  cd proto
  buf generate --template buf.gen.pulsar.yaml
)

# echo "Generate Pulsar Test Data"
# (
#   cd testutil/testdata
#   buf generate --template buf.gen.pulsar.yaml
# )
