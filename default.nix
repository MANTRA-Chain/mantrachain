{
  lib,
  stdenv,
  buildGo123Module,
  nix-gitignore,
  darwin,
  network ? "mainnet", # mainnet|testnet
  rev ? "dirty",
  static ? stdenv.hostPlatform.isStatic,
  nativeByteOrder ? true, # nativeByteOrder mode will panic on big endian machines
  wget,
}:
let
  version = "v5.0.0";
  pname = "mantrachaind";
  wasmvmVersion = "v3.0.0-ibc2.1";
  tags = [
    "ledger"
    "netgo"
    network
  ] ++ lib.optionals nativeByteOrder [ "nativebyteorder" ];
  ldflags = [
    "-X github.com/cosmos/cosmos-sdk/version.Name=mantrachain"
    "-X github.com/cosmos/cosmos-sdk/version.AppName=${pname}"
    "-X github.com/cosmos/cosmos-sdk/version.Version=${version}"
    "-X github.com/cosmos/cosmos-sdk/version.BuildTags=${lib.concatStringsSep "," tags}"
    "-X github.com/cosmos/cosmos-sdk/version.Commit=${rev}"
  ];
  buildInputs = [ wget ];
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
  vendorHash = "sha256-AJRbAMOf7IkkZ43wHUH1PxJGu0RwQi7cIIp7kdBV0/E=";
  proxyVendor = true;
  modules = ./gomod2nix.toml;
  pwd = src; # needed to support replace
  subPackages = [ "cmd/mantrachaind" ];
  CGO_ENABLED = "1";

  postFixup = lib.optionalString (stdenv.isDarwin) ''
    mkdir -p $out/lib
    ${wget}/bin/wget --no-check-certificate https://github.com/CosmWasm/wasmvm/releases/download/${wasmvmVersion}/libwasmvm.dylib -O $out/lib/libwasmvm.dylib
    ${stdenv.cc.bintools.targetPrefix}install_name_tool -change "@rpath/libwasmvm.dylib" "$out/lib/libwasmvm.dylib" $out/bin/mantrachaind
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
