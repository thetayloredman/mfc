{ nixpkgs }:
{
  forAllSystems =
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
    in
    f:
    nixpkgs.lib.genAttrs systems (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      f { inherit system pkgs; }
    );
}
