package completions

import "fmt"

const zshCompletion = `#compdef snx

_snx() {
    local -a snippets
    snippets=(${(f)"$(snx --list-snippets 2>/dev/null)"})
    _arguments '1: :->snippet' '*: :->args'
    case $state in
        snippet)
            _describe 'snippet' snippets
            ;;
    esac
}

_snx "$@"
`

const bashCompletion = `_snx_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    if [ "$COMP_CWORD" -eq 1 ]; then
        local snippets
        snippets=$(snx --list-snippets 2>/dev/null)
        COMPREPLY=($(compgen -W "$snippets" -- "$cur"))
    fi
}

complete -F _snx_completions snx
`

// Generate returns the shell completion script for the given shell.
// Supported shells: bash, zsh.
func Generate(shell string) (string, error) {
	switch shell {
	case "zsh":
		return zshCompletion, nil
	case "bash":
		return bashCompletion, nil
	default:
		return "", fmt.Errorf("unsupported shell '%s': supported shells are bash, zsh", shell)
	}
}
