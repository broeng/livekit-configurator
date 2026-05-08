{ config, pkgs, lib, ... }:
let
  cfg = config.services.livekit-sip;
  settingsFormat = pkgs.formats.json { };
  filterConfig = lib.converge (lib.filterAttrsRecursive (_: v: v != null));
in {

  options = {
    services.livekit-sip = {
      enable = lib.mkEnableOption "Enable livekit-sip";

      package = lib.mkOption {
        type = lib.types.package;
        description = "livekit-sip package to run";
      };

      settings = lib.mkOption {
        type = lib.types.attrs;
        description = "settings";
      };
    };
  };

  config = lib.mkIf cfg.enable {

    systemd.services.livekit-sip = {
      description = "livekit-sip service";
      after = [ "network-online.target" ];
      requires = [ "network-online.target" ];
      wantedBy = [ "multi-user.target" ];
      environment = {
        SIP_CONFIG_FILE = settingsFormat.generate "config.json" (filterConfig cfg.settings);
      };
      serviceConfig = {
        Type = "simple";
        User = "root"; # for testing...
        Group = "root";
        ExecStart = "${cfg.package}/bin/livekit-sip";
        Restart = "always";
        RestartSec = 5;
        NonBlocking = true;
      };
    };

  };

}
