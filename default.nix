{
  lib,
  stdenv,
  buildGo123Module,
  nix-gitignore,
  rocksdb,
  darwin,
  network ? "mainnet", # mainnet|testnet
  rev ? "dirty",
  static ? stdenv.hostPlatform.isStatic,
  nativeByteOrder ? true, # nativeByteOrder mode will panic on big endian machines
}:
let
  version = "v4.0.1";
  pname = "mantrachaind";
  tags = [
    "ledger"
    "netgo"
    network
    "rocksdb"
    "grocksdb_no_link"
  ] ++ lib.optionals nativeByteOrder [ "nativebyteorder" ];
  ldflags = [
    "-X github.com/cosmos/cosmos-sdk/version.Name=mantrachain"
    "-X github.com/cosmos/cosmos-sdk/version.AppName=${pname}"
    "-X github.com/cosmos/cosmos-sdk/version.Version=${version}"
    "-X github.com/cosmos/cosmos-sdk/version.BuildTags=${lib.concatStringsSep "," tags}"
    "-X github.com/cosmos/cosmos-sdk/version.Commit=${rev}"
  ];
  buildInputs = [ rocksdb ];
in
buildGo123Module rec {
  inherit
    pname
    version
    buildInputs
    tags
    ldflags
    ;
  src = (
    nix-gitignore.gitignoreSourcePure [
      "/*" # ignore all, then add whitelists
      "!/app/"
      "!/cmd/"
      "!/api/"
      "!/client/"
      "!/docs/"
      "!/testutil/"
      "!/x/"
      "!go.mod"
      "!go.sum"
      "!gomod2nix.toml"
    ] ./.
  );
  vendorHash = "sha256-JfIc+dYBCwvrqpEbbLbpPfWJONgfdB9tYzigZA+J/4M=";
  proxyVendor = true;
  modules = ./gomod2nix.toml;
  pwd = src; # needed to support replace
  subPackages = [ "cmd/mantrachaind" ];
  CGO_ENABLED = "1";
  CGO_LDFLAGS = lib.optionalString (rocksdb != null) (
    if static then
      "-lrocksdb -pthread -lstdc++ -ldl -lzstd -lsnappy -llz4 -lbz2 -lz"
    else if stdenv.hostPlatform.isWindows then
      "-lrocksdb-shared"
    else
      "-lrocksdb -pthread -lstdc++ -ldl"
  );

  postFixup = lib.optionalString (stdenv.isDarwin && rocksdb != null) ''
    ${stdenv.cc.bintools.targetPrefix}install_name_tool -change "@rpath/librocksdb.9.dylib" "${rocksdb}/lib/librocksdb.dylib" $out/bin/mantrachaind
  '';

  doCheck = false;
  meta = with lib; {
    description = "Official implementation of the mantra protocol";
    homepage = "https://www.mantrachain.io/";
    license = licenses.asl20;
    mainProgram = "mantrachaind" + stdenv.hostPlatform.extensions.executable;
    platforms = platforms.all;
  };
}
