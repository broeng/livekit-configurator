{
  config,
  lib,
  pkgs,
  ...
}:
let
  livekitConfigurator = pkgs.callPackage ./packages/livekit-configurator.nix {};
in rec {
  name = "livekit-configurator";

  nodes = {
    machine = let
      lcOpts = lib.concatStringsSep " " [
        "--run-once"
        "--config" ./../config.yaml
        "--log-level" "debug"
        "--server-url" "http://127.0.0.1:${toString nodes.machine.services.livekit.settings.port}"
        "--api-key" "devkey"
        "--api-secret" nodes.machine.services.livekit.settings.keys.devkey
      ];
      livekitConfCmd = pkgs.writeShellScriptBin "livekit-conf" ''
        exec ${livekitConfigurator}/bin/livekit-configurator ${lcOpts} "$@"
      '';
    in {
      imports = [
        ./modules/livekit-sip.nix
      ];
      networking.enableIPv6 = false;
      environment.systemPackages = [
        livekitConfigurator
        livekitConfCmd
        pkgs.net-tools
      ];
      services.livekit = rec {
        enable = true;
        redis = {
          createLocally = true;
          port = 6888;
        };
        settings = {
          log_level = "debug";
          keys = {
            "devkey" = "N0inqPzWzHgrMCQuIJ1fiDZW6U6fzo1BGfF9HtKUfIqB";
          };
          port = 7880;
          rtc = {
            tcp_port = 7881;
            port_range_start = 50000;
            port_range_end = 60000;
            use_external_ip = false;
          };
          redis = {
            address = "127.0.0.1:${toString redis.port}";
          };
          turn = {
            enabled = false;
          };
        };
        keyFile = pkgs.writeText "keyfile" (
          lib.concatStringsSep "\n" (
            lib.mapAttrsToList (
              k: v: "${k}: ${v}")
              settings.keys));
      };
      services.livekit-sip = {
        enable = false;
        package = pkgs.callPackage ./packages/livekit-sip.nix {};
        settings = {
          api_key = "devkey";
          api_secret = nodes.machine.services.livekit.settings.keys.devkey;
          ws_url = "ws://127.0.0.1:${toString nodes.machine.services.livekit.settings.port}";
          sip_port = 5060;
          redis = {
            address = "127.0.0.1:${toString nodes.machine.services.livekit.redis.port}";
          };
          log_level = "debug";
        };
      };
    };
  };

  testScript = let
  in ''
    with subtest("livekit server starts and responds"):
      machine.wait_for_unit("livekit.service")
      machine.wait_for_unit("redis-livekit.service")
      #machine.wait_for_unit("livekit-sip.service")
      #print(machine.succeed("journalctl -u livekit-sip"))
      print(machine.succeed("journalctl -u redis-livekit"))
      machine.wait_until_succeeds("netstat -anp | grep LISTEN | grep -qF 6888")
      machine.wait_until_succeeds("netstat -anp | grep LISTEN | grep -qF 7880")
      #machine.wait_until_succeeds("netstat -anp | grep LISTEN | grep -qF 5060")
      print(machine.succeed("netstat -anp | grep LISTEN"))
      machine.succeed("livekit-conf")
  '';
}
