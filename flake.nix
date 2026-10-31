{
    description = "Compiler for the KrabbaScript language";

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        flake-utils.url = "github:numtide/flake-utils";
    };

    outputs =
        {
            self,
            nixpkgs,
            flake-utils,
        }:
        flake-utils.lib.eachDefaultSystem (
            system:
            let
                pkgs = import nixpkgs { inherit system; };
                krabbascript = pkgs.buildGo127Module {
                    pname = "krabbascript";
                    version = "0.1.0";
                    src = ./.;
                    vendorHash = "sha256-ApEbeTZrLDERd/P9u/j0r/FNRiDpJZuDEMOY/Aquik8=";
                };
            in
            {
                packages.default = krabbascript;
                devShells.default = pkgs.mkShell {
                    packages = with pkgs; [
                        go_1_27
                        gopls
                        golangci-lint
                        gotools
                        llvmPackages_latest.llvm
                    ];

                    shellHook = ''
                        printf '\033[96mKrabbascript development environment loaded!\033[0m\n'
                        go version
                    '';
                };
            }
        );
}
