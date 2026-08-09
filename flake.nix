{
  description = "SNX — Snippet Executor";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    goShell.url = "git+https://gitea.lightjack.de/LightJack05/nix-library?dir=shells/go";
    generalLib.url = "git+https://gitea.lightjack.de/LightJack05/nix-library?dir=lib/general";
    # --- Optional libs (uncomment input + merge lines below to enable) ---
    # podmanLib.url = "git+https://gitea.lightjack.de/LightJack05/nix-library?dir=lib/podman";
    # qemuLib.url = "git+https://gitea.lightjack.de/LightJack05/nix-library?dir=lib/qemu";
    goLicenseCollectorLib.url = "git+https://gitea.lightjack.de/LightJack05/nix-library?dir=lib/go-license-collector";
  };

  outputs = { self, nixpkgs, goShell, generalLib, goLicenseCollectorLib, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          # Flakes cannot see git tags, so the version identifies the exact
          # commit instead; semver tags are handled by the release workflow.
          date = nixpkgs.lib.substring 0 8 (self.lastModifiedDate or "19700101");
          rev = self.shortRev or self.dirtyShortRev or "dirty";
          version = "unstable-${date}-${rev}";
        in
        rec {
          snx = pkgs.buildGoModule {
            pname = "snx";
            inherit version;

            src = self;

            vendorHash = "sha256-pbA/AlBz3cQYRTMnQ/qBPcinYOKokrBLNhkbRTq54gE=";

            env.CGO_ENABLED = 0;

            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
            ];

            nativeBuildInputs = [ pkgs.installShellFiles ];

            postInstall = ''
              installShellCompletion --cmd snx \
                --bash completions/snx.bash \
                --zsh completions/snx.zsh
              install -Dm644 config.example.toml $out/share/doc/snx/config.example.toml
              install -Dm644 README.md $out/share/doc/snx/README.md
            '';

            meta = {
              description = "Snippet Executor — quickly invoke personal scripts from a central directory";
              homepage = "https://github.com/LightJack05/snx";
              license = nixpkgs.lib.licenses.mit;
              mainProgram = "snx";
              platforms = systems;
            };
          };

          default = snx;
        }
      );

      apps = forAllSystems (system: rec {
        snx = {
          type = "app";
          program = nixpkgs.lib.getExe self.packages.${system}.snx;
        };
        default = snx;
      });

      overlays.default = final: prev: {
        snx = self.packages.${final.stdenv.hostPlatform.system}.snx;
      };

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

          # --- Add project-specific packages here ---
          extraPackages = [
          ];

          # --- Add project-specific shell hook here (env vars, startup messages, etc.) ---
          extraShellHook = ''
          '';

          # --- Optional lib packages (uncomment matching input above to enable) ---
          optionalPackages = []
          # ++ podmanLib.packages.${system}
          # ++ qemuLib.packages.${system}
          ;

          # --- Optional lib hooks (uncomment matching input above to enable) ---
          optionalHook = ""
          # + podmanLib.shellHook
          # + qemuLib.shellHook
          ;
        in
        {
          default = pkgs.mkShell {
            name = "go-dev-shell";
            packages = goShell.shellConfig.${system}.packages
              ++ generalLib.packages.${system}
              ++ goLicenseCollectorLib.packages.${system}
              ++ optionalPackages
              ++ extraPackages;
            shellHook = goShell.shellConfig.${system}.shellHook
              + generalLib.shellHook
              + optionalHook
              + extraShellHook;
          };
        }
      );
    };
}
