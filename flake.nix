{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-24.11";
    flake-utils.url = "github:numtide/flake-utils";
    nix-bundle-exe = {
      url = "github:3noch/nix-bundle-exe";
      flake = false;
    };
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
    poetry2nix = {
      url = "github:nix-community/poetry2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      nix-bundle-exe,
      gomod2nix,
      flake-utils,
      poetry2nix,
    }:
    let
      rev = self.shortRev or "dirty";
      mkApp = drv: {
        type = "app";
        program = "${drv}/bin/${drv.meta.mainProgram}";
      };
    in
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = self.overlays.default;
          config = { };
        };
      in
      rec {
        packages = pkgs.mantra-matrix // {
        };
        apps = {
          mantrachaind = mkApp packages.mantrachaind;
        };
        defaultPackage = packages.mantrachaind;
        defaultApp = apps.mantrachaind;
        devShells = rec {
          default = pkgs.mkShell {
            buildInputs = [
              defaultPackage.go
              pkgs.nixfmt-rfc-style
              pkgs.gomod2nix
              pkgs.poetry
            ];
          };
          full = pkgs.mkShell { buildInputs =[ pkgs.test-env ]; };
          test = pkgs.mkShell {
            buildInputs = full.buildInputs ++ [ packages.mantrachaind ];
          };
        };
        legacyPackages = pkgs;
      }
    ))
    // {
      overlays.default = [
        poetry2nix.overlays.default
        gomod2nix.overlays.default
        (final: super: {
          go = super.go_1_23;
          test-env = final.callPackage ./nix/testenv.nix { };
          mantra-matrix = final.callPackage ./nix/mantra-matrix.nix {
            inherit rev;
            bundle-exe = final.pkgsBuildBuild.callPackage nix-bundle-exe { };
          };
        })
      ];
    };
}
