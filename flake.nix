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
          shellHook = ''
            export LIVEKIT_CONF_SERVER_URL=http://localhost:8080
            export LIVEKIT_CONF_API_KEY=dummy
            export LIVEKIT_CONF_API_SECRET=dummy
            export LIVEKIT_CONF_CONFIG=config.yaml
          '';
        };
      }
    );
}
