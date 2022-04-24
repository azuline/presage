{ pkgs ? import <nixpkgs> { } }:

let
  env = pkgs.buildEnv {
    name = "env";
    paths = with pkgs; [
      golangci-lint
      go_1_18
      gopls
      gopkgs
      go-outline
      delve
      gotools
    ];
  };
in
pkgs.mkShell {
  buildInputs = [ env ];
}
