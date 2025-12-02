package go_tmp_env

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/New-Horizons-Team/supply-chain-firewall/internal/temp_directory"
)

type TempGoEnvironment interface {
	io.Closer

	// Run executes a go sub-command in the temporary environment and returns its stdout.
	// If local is true, the command is executed on the caller's CWD.
	Run(ctx context.Context, args []string, local bool) (string, error)

	// GetRoot retrieves the base directory for this temporary environment.
	GetRoot() string
}

type tempGoEnvironment struct {
	// The temporary directory where the dummy environment will be prepared
	dir temp_directory.TempDir
	// Path to the go's executable binary.
	executable string
	// The original dir, for which this environment was created.
	originalDir string
	// Temporary directory where the commands are executed.
	dryRunDir string
	// Environment variables set to every go command execute in the environment.
	env []string
	// Whether the original dir's go.mod was modified and should be restored.
	restoreModFile bool
	// Whether the original dir's go.sum was modified and should be restored.
	restoreSumFile bool
	// Whether the original dir didn't have a go.sum file when this environment was created.
	removeSumFile bool
}

// New creates a temporary go environment based on the current directory.
func New(ctx context.Context, executable string) (tmp TempGoEnvironment, err error) {
	var dir temp_directory.TempDir

	dir, err = temp_directory.NewTempDir("go-tmp-env-")
	if err != nil {
		err = errors.Join(ErrCreateTempDir, err)
		return
	}
	defer func() {
		if err != nil {
			_ = dir.Close()
		}
	}()

	var originalDir string
	originalDir, err = getOriginalDir(ctx, executable)
	if err != nil {
		return
	}

	g := &tempGoEnvironment{
		dir:         dir,
		executable:  executable,
		originalDir: originalDir,
	}
	tmp = g

	err = g.createTmpEnv(ctx)
	return
}

// getOriginalDir retrieves the directory with a go.mod
// based on which the temporary environment should be created.
// If no such directory is found, this returns the empty string and nil!
func getOriginalDir(ctx context.Context, executable string) (string, error) {
	cmd := exec.CommandContext(ctx, executable, "env", "GOMOD")
	stdout, err := cmd.Output()
	if err != nil {
		return "", errors.Join(ErrGetOriginalDir, err)
	}

	path := string(stdout)
	if path != "/dev/null" && path != "NUL" {
		path, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			return "", errors.Join(ErrExpandOriginalDir, err)
		}
		return path, nil
	}
	return "", nil
}

// Create the temporary environment and set every environment variable
// required to run `go` commands keeping the global environment clean.
func (tmp *tempGoEnvironment) createTmpEnv(ctx context.Context) error {
	goDir := filepath.Join(tmp.dir.GetPath(), "go")
	err := os.Mkdir(goDir, 0750)
	if err != nil {
		return errors.Join(ErrCreateGoDir, err)
	}

	tmp.dryRunDir = filepath.Join(tmp.dir.GetPath(), "dry_run")
	err = os.Mkdir(tmp.dryRunDir, 0750)
	if err != nil {
		return errors.Join(ErrCreateDryRunDir, err)
	}

	cache := filepath.Join(tmp.dir.GetPath(), "cache")
	err = os.Mkdir(cache, 0750)
	if err != nil {
		return errors.Join(ErrCreateCache, err)
	}

	modCache := filepath.Join(tmp.dir.GetPath(), "mod_cache")
	err = os.Mkdir(modCache, 0750)
	if err != nil {
		return errors.Join(ErrCreateModCache, err)
	}

	// Go searches each directory listed in GOPATH to find source code,
	// but new packages are always downloaded into the first directory
	// in the list.
	cmd := exec.CommandContext(ctx, tmp.executable, "env", "GOPATH")
	baseGoPath, err := cmd.Output()
	if err != nil {
		return errors.Join(ErrGetGoPath, err)
	}

	goPath := fmt.Sprintf("GOPATH=%s%c%s", goDir, os.PathListSeparator, strings.TrimSpace(string(baseGoPath)))
	goCache := fmt.Sprintf("GOCACHE=%s", cache)
	goModCache := fmt.Sprintf("GOMODCACHE=%s", modCache)

	tmp.env = os.Environ()
	tmp.env = append(tmp.env, goPath)
	tmp.env = append(tmp.env, goCache)
	tmp.env = append(tmp.env, goModCache)

	return nil
}

// copyFile copies the given file from dir from to dir to.
func copyFile(file, from, to string) error {
	from = filepath.Join(from, file)
	to = filepath.Join(to, file)

	src, err := os.Open(from)
	if err != nil {
		return errors.Join(ErrCopyOpenSource, err)
	}

	dst, err := os.Create(to)
	if err != nil {
		return errors.Join(ErrCopyOpenDestination, err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return errors.Join(ErrCopyContent, err)
	}

	return nil
}

// Close implements TempGoEnvironment.
func (tmp *tempGoEnvironment) Close() error {
	var err error

	if tmp.originalDir != "" {
		if tmp.restoreModFile {
			tmpErr := copyFile("go.mod", tmp.dir.GetPath(), tmp.originalDir)
			if tmpErr != nil {
				err = errors.Join(ErrRestoreGoMod, tmpErr, err)
			}
		}

		if tmp.restoreSumFile {
			tmpErr := copyFile("go.sum", tmp.dir.GetPath(), tmp.originalDir)
			if tmpErr != nil {
				err = errors.Join(ErrRestoreGoSum, tmpErr, err)
			}
		} else if tmp.removeSumFile {
			tmpErr := os.Remove(filepath.Join(tmp.dir.GetPath(), "go.sum"))
			if tmpErr != nil && !strings.Contains(tmpErr.Error(), "no such file or directory") {
				err = errors.Join(ErrRemoveGoSum, tmpErr, err)
			}
		}
	}

	tmpErr := tmp.dir.Close()
	return errors.Join(tmpErr, err)
}

// Run implements TempGoEnvironment.
func (tmp *tempGoEnvironment) Run(ctx context.Context, args []string, local bool) (string, error) {
	if local {
		err := tmp.duplicateGoMod()
		if err != nil {
			return "", err
		}
	}

	cmd := exec.CommandContext(ctx, tmp.executable, args...)
	cmd.Env = tmp.env

	if !local {
		cmd.Dir = tmp.dryRunDir
	}

	stdout, err := cmd.Output()
	if err != nil {
		err = errors.Join(ErrRunCommand, err)
	}

	return string(stdout), err
}

// GetRoot implements TempGoEnvironment.
func (tmp *tempGoEnvironment) GetRoot() string {
	return tmp.dir.GetPath()
}

// duplicateGoMod duplicates the go.mod and go.sum in the nearest ancestor directory.
// On clean up, these files are recovered, in case they were modified.
func (tmp *tempGoEnvironment) duplicateGoMod() error {
	if tmp.originalDir == "" {
		return ErrGoModNotFound
	}

	if !tmp.restoreModFile {
		err := copyFile("go.mod", tmp.originalDir, tmp.dir.GetPath())
		if err != nil {
			return errors.Join(ErrCopyGoMod, err)
		}
		tmp.restoreModFile = true
	}

	if !(tmp.restoreSumFile || tmp.removeSumFile) {
		if _, err := os.Stat(filepath.Join(tmp.originalDir, "go.sum")); err == nil {
			err = copyFile("go.sum", tmp.originalDir, tmp.dir.GetPath())
			if err != nil {
				return errors.Join(ErrCopyGoSum, err)
			}

			tmp.restoreSumFile = true
		} else {
			tmp.removeSumFile = true
		}
	}

	return nil
}
