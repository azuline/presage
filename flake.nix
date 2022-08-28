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
        defaultApp = {
          type = "app";
          program = "${defaultPackage}/bin/presage";
        };
        defaultPackage = pkgs.buildGoModule {
          pname = "presage";
          version = "0.1.0";
          src = ./.;
          vendorSha256 = "sha256-O9u8ThzGOcfWrDjA87RaOPez8pfqUo+AcciSSAw2sfk=";
          meta = {
            description = "scrape rss feeds and send emails for new articles";
            homepage = "https://github.com/azuline/presage";
            license = nixpkgs.lib.licenses.agpl3Plus;
          };
        };
        devShell = pkgs.mkShell {
          buildInputs = [
            (pkgs.buildEnv {
              name = "env";
              paths = with pkgs; [
                gnumake
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
