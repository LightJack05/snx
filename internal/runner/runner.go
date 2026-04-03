package runner

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

// SnippetNotFoundError is returned when the requested snippet does not exist.
type SnippetNotFoundError struct {
	Name string
	Dir  string
	Msg  string
}

func (e *SnippetNotFoundError) Error() string {
	return e.Msg
}

// SnippetNotExecutableError is returned when the snippet exists but lacks the execute bit.
type SnippetNotExecutableError struct {
	Name string
	Path string
}

func (e *SnippetNotExecutableError) Error() string {
	return fmt.Sprintf("snippet '%s' exists but is not executable — run: chmod +x %s", e.Name, e.Path)
}

// ExitError is returned when the snippet process exits with a non-zero status.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("snippet exited with code %d", e.Code)
}

// Run executes the named snippet from snippetDir, forwarding args and all stdio.
func Run(snippetDir, snippetName string, args []string) error {
	fullPath := filepath.Join(snippetDir, snippetName)

	info, err := os.Stat(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Check whether the directory itself exists to give a better message.
			if _, dirErr := os.Stat(snippetDir); errors.Is(dirErr, os.ErrNotExist) {
				return &SnippetNotFoundError{
					Name: snippetName,
					Dir:  snippetDir,
					Msg:  fmt.Sprintf("snippet '%s' not found in %s (directory does not exist)", snippetName, snippetDir),
				}
			}
			return &SnippetNotFoundError{
				Name: snippetName,
				Dir:  snippetDir,
				Msg:  fmt.Sprintf("snippet '%s' not found in %s", snippetName, snippetDir),
			}
		}
		return fmt.Errorf("could not stat snippet '%s': %w", snippetName, err)
	}

	if info.Mode()&0111 == 0 {
		return &SnippetNotExecutableError{Name: snippetName, Path: fullPath}
	}

	cmd := exec.Command(fullPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &ExitError{Code: exitErr.ExitCode()}
		}
		return fmt.Errorf("failed to run snippet '%s': %w", snippetName, err)
	}

	return nil
}

// List returns a sorted slice of executable snippet names found in snippetDir.
// If the directory does not exist, an empty slice is returned without error.
func List(snippetDir string) ([]string, error) {
	entries, err := os.ReadDir(snippetDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("could not read snippet directory %s: %w", snippetDir, err)
	}

	var names []string
	for _, entry := range entries {
		// Accept regular files and symlinks; skip directories.
		if entry.IsDir() {
			continue
		}

		var mode fs.FileMode
		if entry.Type()&fs.ModeSymlink != 0 {
			// Resolve the symlink to get the real file info.
			info, err := os.Stat(filepath.Join(snippetDir, entry.Name()))
			if err != nil {
				continue
			}
			mode = info.Mode()
		} else if entry.Type().IsRegular() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			mode = info.Mode()
		} else {
			continue
		}

		if mode&0111 != 0 {
			names = append(names, entry.Name())
		}
	}

	sort.Strings(names)
	return names, nil
}
