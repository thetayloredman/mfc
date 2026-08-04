{
  self,
  pkgs,
  system,
  ...
}:
pkgs.buildGoModule {
  pname = "mfc";
  version = "0.0.1";
  src = ../../.;
  proxyVendor = true;
  vendorHash = "sha256-+ofOeiAc3c99TdCn5XBEm3kltkRM0/h1UBtQwGJm3lA=";
  subPackages = [ "." ];
}
