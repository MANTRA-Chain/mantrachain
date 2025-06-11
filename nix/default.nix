{
  sources ? import ./sources.nix,
  system ? builtins.currentSystem,
  ...
}:
import sources.nixpkgs {
  overlays = [
    (import ./build_overlay.nix)
    (_: pkgs: {
      flake-compat = import sources.flake-compat;
    })
    (import "${sources.poetry2nix}/overlay.nix")
    (import "${sources.gomod2nix}/overlay.nix")
    (_: pkgs: { test-env = pkgs.callPackage ./testenv.nix { }; })
  ];
  config = { };
  inherit system;
}
