{ pkgs, lib, ... }:
pkgs.buildGo126Module rec {
  pname = "livekit-sip";
  version = "1.3.0";
  src = pkgs.fetchFromGitHub {
    owner = "livekit";
    repo = "sip";
    rev = "v${version}";
    sha256 = "sha256-JSZRtV9FAN6+4poVhCOe+6X6fzLDEDCDCmZX2BJjazc=";
  };
  vendorHash = "sha256-dI47gKRhE7MdssnK2UEvw+MqkaCII5ZLYRWsRp0eXEE=";
  subPackages = [
    "cmd/livekit-sip"
  ];
  nativeBuildInputs = [
    pkgs.pkg-config
  ];
  buildInputs = [
    pkgs.libopus
    pkgs.opusfile
    pkgs.soxr
    pkgs.libogg
  ];
}
