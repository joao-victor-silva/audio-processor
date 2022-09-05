{pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/refs/tags/22.05.tar.gz") {}}:
pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.gopls
    pkgs.golangci-lint
    pkgs.SDL2
    pkgs.gitlint
    pkgs.codespell
    pkgs.nodePackages.cspell
    pkgs.gh
  ];

  shellHook = ''
  '';

  DIAGNOSTICS = "golangci_lint:gitlint:codespell";
  FORMATTING = "";
}
