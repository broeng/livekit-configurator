{
  lib,
  pkgs,
  ...
}:
pkgs.buildGo126Module {
  pname = "livekit-configurator";
  version = "0.1-dev";
  src = ./../..;
  vendorHash = "sha256-GNwQWH6fI0rm3BlT83IaDK/h0gRhegPyAPMI7fh5gtM=";
}
