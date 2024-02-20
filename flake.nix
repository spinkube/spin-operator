{
  description = "The Spin Operator for Kubernetes";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    nixpkgs-format.url = "github:nix-community/nixpkgs-fmt";
    nixpkgs-format.inputs.nixpkgs.follows = "nixpkgs";
    nixpkgs-format.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, nixpkgs-format }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        buildDeps = with pkgs; [
          go_1_22
          gnumake
          git
        ];

        devDeps = with pkgs; buildDeps ++ [
          gopls
          kubectl
          kubernetes-helm
          gotestsum
          golangci-lint
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = devDeps ++ [
            nixpkgs-format.defaultPackage.${system}
          ];
        };
      });
}
