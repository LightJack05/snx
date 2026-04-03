package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LightJack05/snx/internal/completions"
	"github.com/LightJack05/snx/internal/config"
	"github.com/LightJack05/snx/internal/runner"
)

// version is set at build time via -ldflags "-X main.version=<ver>".
var version = "dev"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "--help", "-h":
		printHelp()
		os.Exit(0)

	case "--version":
		fmt.Printf("snx %s\n", version)
		os.Exit(0)

	case "--list-snippets", "-l":
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "snx: %v\n", err)
			os.Exit(1)
		}
		names, err := runner.List(cfg.SnippetDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "snx: %v\n", err)
			os.Exit(1)
		}
		for _, name := range names {
			fmt.Println(name)
		}
		os.Exit(0)

	case "--snippet-dir":
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "snx: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(cfg.SnippetDir)
		os.Exit(0)

	case "--completions":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "snx: --completions requires a shell argument (bash, zsh)")
			os.Exit(1)
		}
		shell := args[1]
		script, err := completions.Generate(shell)
		if err != nil {
			fmt.Fprintf(os.Stderr, "snx: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(script)
		os.Exit(0)

	default:
		if strings.HasPrefix(args[0], "-") {
			fmt.Fprintf(os.Stderr, "snx: unknown flag: %s\n", args[0])
			fmt.Fprintln(os.Stderr, "Run 'snx --help' for usage.")
			os.Exit(1)
		}

		// Treat args[0] as snippet name, remainder as passthrough args.
		snippetName := args[0]
		snippetArgs := args[1:]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "snx: %v\n", err)
			os.Exit(1)
		}

		if err := runner.Run(cfg.SnippetDir, snippetName, snippetArgs); err != nil {
			var exitErr *runner.ExitError
			switch e := err.(type) {
			case *runner.ExitError:
				exitErr = e
				os.Exit(exitErr.Code)
			default:
				fmt.Fprintf(os.Stderr, "snx: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func printHelp() {
	cfg, err := config.Load()
	snippetDir := "~/.local/share/snx/snippets"
	if err == nil {
		snippetDir = cfg.SnippetDir
	}

	fmt.Printf(`snx — Snippet Executor %s

USAGE:
    snx <snippet> [args...]    Run a snippet, passing remaining args through
    snx -l                     List available snippets
    snx --snippet-dir          Print the snippet directory path
    snx --completions <shell>  Print shell completion script (bash, zsh)
    snx --version              Print version
    snx --help                 Print this help message

SNIPPET DIRECTORY:
    %s

EXAMPLES:
    snx hello                  Run the 'hello' snippet
    snx greet Alice Bob        Run 'greet' with arguments 'Alice' and 'Bob'
    snx -l                     List all available snippets
    snx --completions zsh      Output zsh completion script

CONFIGURATION:
    ~/.config/snx/config.toml  (optional, all settings have sane defaults)
`, version, snippetDir)
}
