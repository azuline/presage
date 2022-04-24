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
  shellHook = ''
    # Isolate build stuff to this repo's directory.

    export PRESAGE_ROOT="$(pwd)"
    export PRESAGE_CACHE_ROOT="$(pwd)/.cache"

    export GOCACHE="$PRESAGE_CACHE_ROOT/go/cache"
    export GOENV="$PRESAGE_CACHE_ROOT/go/env"
    export GOPATH="$PRESAGE_CACHE_ROOT/go/path"
    export GOMODCACHE="$GOPATH/pkg/mod"
    export GOROOT=
    export PATH=$(go env GOPATH)/bin:$PATH
  '';
}
