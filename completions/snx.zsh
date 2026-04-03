#compdef snx

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
