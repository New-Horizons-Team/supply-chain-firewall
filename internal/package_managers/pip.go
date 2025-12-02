package package_managers

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	int_logger "github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/slice_utils"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
	"golang.org/x/mod/semver"
)

const minPipVersion = "22.2"

type pipPackageManager struct {
	basicManager
}

// init registers a getter for pip.
func init() {
	packageManagers["pip"] = func(ctx context.Context) (PackageManager, error) {
		pip := &pipPackageManager{}

		err := pip.init(ctx)
		if err != nil {
			return nil, err
		}

		err = pip.validateVersion(ctx)
		if err != nil {
			return nil, err
		}

		return pip, nil
	}
}

// init inializes this manager with values from the context.
func (p *pipPackageManager) init(ctx context.Context) error {
	p.basicManager.init(ctx)

	if p.executable == "" {
		venv := os.Getenv("VIRTUAL_ENV")

		if venv != "" {
			p.executable = filepath.Join(venv, "bin", "python")
		} else {
			var err error

			p.executable, err = exec.LookPath("python")
			if err != nil {
				return errors.Join(ErrMissingPython, err)
			}
		}
	}

	if p.executable == "" {
		return ErrMissingExecutable
	}
	if stat, err := os.Stat(p.executable); err != nil {
		return errors.Join(ErrInvalidExecutable, err)
	} else if !stat.Mode().IsRegular() {
		return ErrInvalidExecutable
	}

	return nil
}

// validateVersion checks if the package manager version is supported.
func (p *pipPackageManager) validateVersion(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, p.executable, "-m", "pip", "--version")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.Join(ErrGetVersion, err)
	}

	// All supported versions adhere to this format.
	versionOut := strings.Split(string(stdout), " ")
	if len(versionOut) < 2 || !semver.IsValid("v"+versionOut[1]) {
		return ErrGetVersion
	} else if semver.Compare(versionOut[1], minPipVersion) < 0 {
		return ErrInvalidVersion
	}

	return nil
}

// Name implements PackageManager.
func (*pipPackageManager) Name() string {
	return "pip"
}

// Ecosystem implements PackageManager.
func (*pipPackageManager) Ecosystem() target.Ecosystem {
	return target.EcosystemPypi
}

// normalizeCommand normalizes a pip command, so it may be executed.
func (p *pipPackageManager) normalizeCommand(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, ErrEmptyCommand
	} else if args[0] != p.Name() {
		return nil, ErrInvalidCommand
	} else {
		args = append([]string{p.Executable(), "-m"}, args...)
		return args, nil
	}
}

// RunCommand implements PackageManager.
func (p *pipPackageManager) RunCommand(ctx context.Context, command []string) error {
	command, err := p.normalizeCommand(command)
	if err != nil {
		return err
	}

	return p.basicManager.RunCommand(ctx, command)
}

// ResolveInstallTargets implements PackageManager.
func (p *pipPackageManager) ResolveInstallTargets(ctx context.Context, command []string) ([]target.Package, error) {
	// pip only installs or upgrades packages via the `pip install` subcommand
	// If `install` is not present, the command is automatically safe to run.
	if !slices.Contains(command, "install") {
		return nil, nil
	}

	// If `install` is present with any of the below options, a usage or error
	// message is printed or a dry-run install occurs: nothing will be installed.
	shouldSkip := slices.ContainsFunc(command, func(s string) bool {
		switch s {
		case "-h",
			"--help",
			"--dry-run":

			return true
		default:
			return false
		}
	})
	if shouldSkip {
		return nil, nil
	}

	// Otherwise, this is probably a live `pip install` command
	// To be certain, we would need to write a full parser for pip.
	command, err := p.normalizeCommand(command)
	if err != nil {
		return nil, err
	}

	command = append(command, "--dry-run", "--quiet", "--report", "-")
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		logger := int_logger.GetLogger(ctx)
		logger.Info(ctx, "Encountered an error while resolving pip installation targets")
		return nil, errors.Join(ErrDryRun, err)
	}

	var report pipReport
	err = json.Unmarshal(stdout, &report)
	if err != nil {
		return nil, errors.Join(ErrParseTargets, err)
	}

	fn := func(v pipInstallMetadata) target.Package {
		var tmp target.Package

		if err != nil {
			return tmp
		}

		if v.Metadata.Name == "" {
			err = ErrPipMissingName
			return tmp
		} else if v.Metadata.Version == "" {
			err = ErrPipMissingVersion
			return tmp
		}

		return target.Package{
			Ecosystem: p.Ecosystem(),
			Name:      v.Metadata.Name,
			Version:   v.Metadata.Version,
		}
	}

	packages := slice_utils.Map(report.Install, fn, nil)

	return packages, err
}

type pipReport struct {
	Install []pipInstallMetadata `json:"install"`
}

type pipInstallMetadata struct {
	Metadata pipInstallTarget `json:"metadata"`
}

type pipInstallTarget struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
