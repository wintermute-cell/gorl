{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { inherit system; };
          longgopher = pkgs.buildGoModule rec {
            pname = "longgopher";
            version = "0.0.3";
            src = pkgs.fetchFromGitHub {
              owner = "sheepla";
              repo = "longgopher";
              rev = "v${version}";
              sha256 = "sha256-q0I53lIrJBEx171iEDzrkemjQbDvhi0K8Snhwso4K5Y=";
            };
            vendorHash = "sha256-nzPHx+c369T4h9KETqMurxZK3LsJAhwBaunkcWIW3Ps=";
            subPackages = [ "." ];
          };
        in
        {
          devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [ 
              longgopher
              libGL
              xorg.libXi
              xorg.libXcursor
              xorg.libXrandr
              xorg.libXinerama
              wayland
              libxkbcommon
            ];
            shellHook = ''
              if [ -z "$IN_DEV_SHELL" ]; then
                echo -e "\033[1;32mEntering Nix shell...\033[0m"
                export IN_DEV_SHELL=1
                export PS1="[Nix] $PS1"
                longgopher -l 5
                echo ""
              else
                echo -e "\033[1;31mAlready in Nix shell!\033[0m"
                exit 1
              fi
              exec $SHELL
            '';
          };
        }
      );
}
