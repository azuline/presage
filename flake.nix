{
  description = "presage development environment";

  inputs = {
    nixpkgs = {
      url = github:nixos/nixpkgs/nixos-unstable;
    };
    flake-utils = {
      url = github:numtide/flake-utils;
    };
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in rec {
        defaultPackage = devShells.default;
        devShells.default = pkgs.mkShell {
          buildInputs = [
            (pkgs.buildEnv {
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
            })
          ];
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
        };
      }
    );
}
