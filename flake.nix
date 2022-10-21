{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-22.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      inherit (nixpkgs) lib;
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      defaultPackage = pkgs.buildGoModule {
        name = "kubenetbench";

        src = self;

        vendorSha256 = "sha256-pEnpCMhutntoIM2U2fZP5wnjsH6v6moSdAK+5mrJDEM=";

        subPackages = [ "kubenetbench" ];
      };
    });
}