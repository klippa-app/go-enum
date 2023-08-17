{ pkgs, ... }:

{
  packages = with pkgs; [ git go_1_20 ];
}
