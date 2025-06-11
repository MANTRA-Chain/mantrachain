{
  system ? builtins.currentSystem,
  pkgs ? import ../nix { inherit system; },
}:
pkgs.mkShell {
  buildInputs = [
    (pkgs.callPackage ../. { }) # mantrachaind
    pkgs.test-env
    pkgs.poetry
  ];
  shellHook = ''
    export TMPDIR=/tmp
  '';
}
