_snx_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    if [ "$COMP_CWORD" -eq 1 ]; then
        local snippets
        snippets=$(snx --list-snippets 2>/dev/null)
        COMPREPLY=($(compgen -W "$snippets" -- "$cur"))
    fi
}

complete -F _snx_completions snx
