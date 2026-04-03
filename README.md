# SNX — Snippet Executor

SNX is a minimal CLI tool for quickly invoking personal scripts from a central directory. Instead of cluttering your `PATH` with dozens of one-off scripts, you keep them all in a single snippets folder and run them via `snx <name>`. SNX handles argument forwarding, stdin/stdout/stderr passthrough, and exit code propagation so your snippets behave exactly as if they were called directly.

---

## Installation

### Arch Linux (PKGBUILD)

A `PKGBUILD` is provided in the `packaging/` directory. You can build and install the package locally with:

```sh
cp packaging/PKGBUILD /tmp/snx-pkg/
cd /tmp/snx-pkg
makepkg -si
```

This installs the binary to `/usr/bin/snx` and places completion files in the system completion directories automatically.

### Other Linux (binary from GitHub Releases)

Download the latest binary from the [GitHub Releases page](https://github.com/LightJack05/snx/releases) and install it:

```sh
# Replace <version> and <arch> as appropriate (e.g. v0.1.0, amd64)
curl -L https://github.com/LightJack05/snx/releases/download/<version>/snx-linux-<arch> -o snx
chmod +x snx
sudo mv snx /usr/local/bin/snx
```

---

## First-Time Setup

Create the default snippets directory and add your scripts:

```sh
mkdir -p ~/.local/share/snx/snippets
```

Place any executable script in that directory:

```sh
cat > ~/.local/share/snx/snippets/hello <<'EOF'
#!/bin/sh
echo "Hello, ${1:-world}!"
EOF
chmod +x ~/.local/share/snx/snippets/hello
```

Run it:

```sh
snx hello
snx hello Alice
```

---

## Autocomplete Setup

### zsh

Add to `~/.zshrc`:

```sh
eval "$(snx --completions zsh)"
```

Or, if SNX was installed via the PKGBUILD, the completion file is already installed at `/usr/share/zsh/site-functions/_snx` and will be picked up automatically by zsh's completion system.

### bash

Add to `~/.bashrc`:

```sh
eval "$(snx --completions bash)"
```

Or, if installed via the PKGBUILD, the completion file is at `/usr/share/bash-completion/completions/snx` and is sourced automatically by bash-completion.

---

## Usage

```
snx <snippet> [args...]    Run a snippet, passing remaining args through
snx -l                     List available snippets (one per line)
snx --snippet-dir          Print the configured snippet directory path
snx --completions <shell>  Print shell completion script (bash or zsh)
snx --version              Print the version string
snx --help                 Print help and exit
```

### Examples

```sh
snx backup                        # run the 'backup' snippet
snx deploy production --dry-run   # run 'deploy' with extra args
snx -l                            # list all available snippets
snx --snippet-dir                 # show where snippets live
snx --completions zsh             # output the zsh completion script
```

---

## Configuration

SNX is configured via `~/.config/snx/config.toml`. This file is **optional** — all settings have sane defaults.

See [`config.example.toml`](config.example.toml) for a documented template with all available options.

| Setting       | Default                            | Description                          |
|---------------|------------------------------------|--------------------------------------|
| `snippet_dir` | `~/.local/share/snx/snippets`      | Directory containing your snippets   |

Example `~/.config/snx/config.toml`:

```toml
snippet_dir = "~/scripts/snippets"
```

---

## License

MIT — see [LICENSE](LICENSE).
