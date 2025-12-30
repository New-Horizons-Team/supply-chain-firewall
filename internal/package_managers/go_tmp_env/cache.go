package go_tmp_env

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetCacheDir retrieves the path to the temporary go cache.
func GetCacheDir() string {
	const cacheDir = "scfw-go-cache"

	tmp := os.TempDir()
	return filepath.Join(tmp, cacheDir)
}

// createCache create the go cache from the current environment,
// returning the environment variables that should be updated and their values.
// If the cache already exists, this is a noop.
func createCache(ctx context.Context, executable string) (map[string]string, error) {
	baseDir := GetCacheDir()
	cacheDir := filepath.Join(baseDir, "cache")
	modDir := filepath.Join(baseDir, "mod")

	envvars := map[string]string{
		"GOCACHE":    cacheDir,
		"GOMODCACHE": modDir,
	}

	// ~/.cache/go-build
	err := copyIfNotExists(ctx, executable, cacheDir, "GOCACHE")
	if err != nil {
		return envvars, errors.Join(ErrCopyCacheDir, err)
	}

	// ${GOPATH}/pkg/mod
	err = copyIfNotExists(ctx, executable, modDir, "GOMODCACHE")
	if err != nil {
		return envvars, errors.Join(ErrCopyModDir, err)
	}

	return envvars, nil
}

// createIfNotExists copies the directory in the provided go envvar
// if it doesn't exist.
func copyIfNotExists(ctx context.Context, executable, dir, sourceVar string) error {
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		// Assume that the cache directory is up to date.
		return nil
	} else if err == nil && !info.IsDir() {
		return ErrInvalidCacheDir
	} else if !os.IsNotExist(err) {
		return errors.Join(ErrCheckCacheDir, err)
	}

	cmd := exec.CommandContext(ctx, executable, "env", sourceVar)
	stdout, err := cmd.Output()
	if err != nil {
		return errors.Join(ErrGetCacheDir, err)
	}
	sourceDir := strings.TrimSpace(string(stdout))

	hasSourceDir := false
	info, err = os.Stat(sourceDir)
	if err == nil && info.IsDir() {
		hasSourceDir = true
	} else if err == nil && !info.IsDir() {
		return ErrInvalidCacheSource
	} else if !os.IsNotExist(err) {
		return errors.Join(ErrCheckSourceCache, err)
	}

	if hasSourceDir {
		err = os.CopyFS(dir, os.DirFS(sourceDir))
		if err != nil {
			return errors.Join(ErrCopyCache, err)
		}
	}

	return nil
}

// PurgeCache deletes the temporary cache used to install packages in an isolated environment.
func PurgeCache() error {
	err := os.RemoveAll(GetCacheDir())
	if err != nil {
		err = errors.Join(ErrPurgeCache, err)
	}
	return err
}
