{
  pkgs,
  ...
}:
{
  default = pkgs.mkShell {
    buildInputs = with pkgs; [
      go
      nixfmt
    ];
  };
}
