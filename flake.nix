{
  description = "mfc";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    inputs@{
      self,
      nixpkgs,
    }:
    let
      inherit (import ./nix/lib.nix { inherit nixpkgs; }) forAllSystems;
    in
    {
      devShells = forAllSystems (args: import ./nix/devshell.nix (inputs // args));
      packages = forAllSystems (args: import ./nix/pkgs/default.nix (inputs // args));
    };
}
