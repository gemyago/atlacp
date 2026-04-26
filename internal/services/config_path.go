package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const fileParentDirPerm = 0o750

// AtlacpPathResolver resolves Atlacp-managed file paths and can create missing parent
// directories for any resolved file path.
// Use [NewAtlacpPathResolver] for production defaults; tests may set fields directly.
type AtlacpPathResolver struct {
	Home           string
	ExecutablePath string
}

// NewAtlacpPathResolver returns a resolver using $HOME.
func NewAtlacpPathResolver() *AtlacpPathResolver {
	executablePath, _ := os.Executable()
	return &AtlacpPathResolver{
		Home:           os.Getenv("HOME"),
		ExecutablePath: executablePath,
	}
}

func (r *AtlacpPathResolver) defaultBaseDir() string {
	return filepath.Join(r.Home, ".atlacp")
}

// DefaultAccountsFilePath returns the default filesystem path for the Atlassian accounts file.
func (r *AtlacpPathResolver) DefaultAccountsFilePath() (string, error) {
	return filepath.Join(r.defaultBaseDir(), "accounts.json"), nil
}

// DefaultLogPath returns the default log file path for the given filename.
// go run keeps logs in the current working directory; compiled binaries write under ~/.atlacp.
func (r *AtlacpPathResolver) DefaultLogPath(logFileName string) (string, error) {
	if logFileName == "" {
		return "", errors.New("log file name is empty")
	}

	if isGoRunExecutablePath(r.ExecutablePath) {
		return filepath.Clean(logFileName), nil
	}

	return filepath.Join(r.defaultBaseDir(), logFileName), nil
}

// EnsureParentDirsForFile creates the parent directory of filePath if needed (idempotent).
// It MUST be called before writing files when the directory may not exist yet.
func (r *AtlacpPathResolver) EnsureParentDirsForFile(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}

	dir := filepath.Dir(filepath.Clean(filePath))
	if err := os.MkdirAll(dir, fileParentDirPerm); err != nil {
		return fmt.Errorf("create file parent directory %q: %w", dir, err)
	}

	return nil
}

func isGoRunExecutablePath(executablePath string) bool {
	if executablePath == "" {
		return false
	}
	cleaned := filepath.ToSlash(filepath.Clean(executablePath))
	return strings.Contains(cleaned, "/go-build") && strings.Contains(cleaned, "/exe/")
}
