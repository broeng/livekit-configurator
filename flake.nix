{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {

        formatter = pkgs.nixfmt-tree;

        checks = {
          test-livekit = pkgs.testers.runNixOSTest {
            imports = [ ./tests/livekit.nix ];
          };
        };

        devShell = pkgs.mkShell {
          buildInputs = [ pkgs.go_1_26 ];
        };
      }
    );
}
