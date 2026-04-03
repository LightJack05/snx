# **Project: SNX (Snippet Executor)**

## **Idea**
- This project aims to simplify calling personal scripts that allow you to do something quickly
    - think: Create a kind cluster, recreate it, initialize a project with defaults, do a system update, etc.
- Tools should be places as executables in a central directory (which should be configurable via a ~/.config/snx/config file)
    - Tools may be bash scripts, ELF binaries, etc.
- The tool should receive any following parameters given to the snx command (so `snx foo bar baz` should call `foo bar baz`)

## **Autocomplete**
- The tool should provide comprehensive autocomplete for the tools to make it easier to execute them, so typing `snx f<tab>` should autocomplete `snx foo` if foo is in the bin directory
- Autocomplete should be handled via the shell. Add support for bash and zsh for now. It should be generatable via `snx --completions zsh` for example.

## **Packaging**
- The tool should be available as a PKGBUILD for Arch Linux, as well as a regular binary file for other platforms.

## **Documentation**
- There should be a readme, as well as a documented sample config file to go with the program.

## **Sane defaults**
- After installing, the program should already be ready to use, make sure the default config works already. Configuration should only be necessary if the user desires to do so.

## **Meaningful error handling**
- Give meaningful error messages. Also pass through errors from the tools below.

## **Additional parameters**
- Add at least the following:
    - A --help / -h parameter
    - A --list-snippets / -l parameter
    - A --snippet-dir parameter to quickly get the current snippet directory

