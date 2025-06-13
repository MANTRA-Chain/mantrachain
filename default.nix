{
  lib,
  stdenv,
  buildGo123Module,
  nix-gitignore,
  darwin,
  rev ? "dirty",
  static ? stdenv.hostPlatform.isStatic,
  nativeByteOrder ? true, # nativeByteOrder mode will panic on big endian machines
  fetchurl,
  pkgsStatic,
}:
let
  version = "v5.0.0";
  pname = "mantrachaind";
  wasmvmVersion = "v3.0.0-ibc2.1";

  # Use static packages for Linux to ensure musl compatibility
  buildPackages = if stdenv.isLinux then pkgsStatic else { inherit stdenv buildGo123Module; };
  buildStdenv = buildPackages.stdenv;
  buildGo123Module' = if stdenv.isLinux then buildPackages.buildGo123Module else buildGo123Module;

  # Download wasmvm libraries as fixed-output derivations
  wasmvmLibs = {
    darwin = fetchurl {
      url = "https://github.com/CosmWasm/wasmvm/releases/download/${wasmvmVersion}/libwasmvmstatic_darwin.a";
      sha256 = "sha256-c2O+xnYvOubIciFN1aphhZHxHW45DzJ7QAQDvVwj1Jk=";
    };
    linux-x86_64 = fetchurl {
      url = "https://github.com/CosmWasm/wasmvm/releases/download/${wasmvmVersion}/libwasmvm_muslc.x86_64.a";
      sha256 = "sha256-/Phan0mC/EludahLAeakxG7ruOOdXNNKLXw/rzAps8s=";
    };
    linux-aarch64 = fetchurl {
      url = "https://github.com/CosmWasm/wasmvm/releases/download/${wasmvmVersion}/libwasmvm_muslc.aarch64.a";
      sha256 = "sha256-fmG1Zp3S2sIkYFwFWlnBwjxS35jXqxocrKwjKti7f4c=";
    };
  };

  wasmvmLib = 
    if buildStdenv.isDarwin then wasmvmLibs.darwin
    else if buildStdenv.isLinux && buildStdenv.hostPlatform.isAarch64 then wasmvmLibs.linux-aarch64
    else if buildStdenv.isLinux then wasmvmLibs.linux-x86_64
    else throw "Unsupported platform for wasmvm";

  tags = [
    "ledger"
    "netgo"
  ] ++ lib.optionals nativeByteOrder [ "nativebyteorder" ]
    ++ lib.optionals buildStdenv.isDarwin [ "static_wasm" ]
    ++ lib.optionals buildStdenv.isLinux [ "muslc" "osusergo" ];

  ldflags = [
    "-X github.com/cosmos/cosmos-sdk/version.Name=mantrachain"
    "-X github.com/cosmos/cosmos-sdk/version.AppName=${pname}"
    "-X github.com/cosmos/cosmos-sdk/version.Version=${version}"
    "-X github.com/cosmos/cosmos-sdk/version.BuildTags=${lib.concatStringsSep "," tags}"
    "-X github.com/cosmos/cosmos-sdk/version.Commit=${rev}"
  ] ++ [
    "-w"
    "-s"
    "-linkmode=external"
  ] ++ lib.optionals buildStdenv.isLinux [
    "-extldflags '-static -lm'"
  ];

in
buildGo123Module' rec {
  inherit
    pname
    version
    tags
    ldflags
    ;
  stdenv = buildStdenv;
  src = (
    nix-gitignore.gitignoreSourcePure [
      "/*"
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

  preBuild = ''
    mkdir -p $TMPDIR/lib
    cp ${wasmvmLib} $TMPDIR/lib/$(basename ${wasmvmLib.name})
    export CGO_LDFLAGS="-L$TMPDIR/lib $CGO_LDFLAGS"
  '';

  doCheck = false;
  meta = with lib; {
    description = "Official implementation of the mantra protocol";
    homepage = "https://www.mantrachain.io/";
    license = licenses.asl20;
    mainProgram = "mantrachaind" + buildStdenv.hostPlatform.extensions.executable;
    platforms = platforms.all;
  };
}