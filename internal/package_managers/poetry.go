package package_managers

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	int_logger "github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
	"golang.org/x/mod/semver"
)

const minPoetryVersion = "1.17.0"

var reGetPoetryVersion = regexp.MustCompile(`Poetry \(version (.*)\)`)

var reGetPoetryPackageVersion = regexp.MustCompile(`(?:Installing|Updating|Downgrading) (?:the current project: )?(.*) \((.*)\)`)

var inspectedPoetrySubCommands = []string{"add", "install", "sync", "update"}

type poetryPackageManager struct {
	basicManager
}

// init registers a getter for poetry.
func init() {
	packageManagers["poetry"] = func(ctx context.Context) (PackageManager, error) {
		pkgMngr := &poetryPackageManager{}

		err := pkgMngr.init(ctx)
		if err != nil {
			return nil, err
		}

		err = pkgMngr.validateVersion(ctx)
		if err != nil {
			return nil, err
		}

		return pkgMngr, nil
	}
}

// init inializes this manager with values from the context.
func (p *poetryPackageManager) init(ctx context.Context) error {
	p.basicManager.init(ctx)

	if p.executable == "" {
		var err error

		p.executable, err = exec.LookPath("poetry")
		if err != nil {
			return errors.Join(ErrMissingPoetry, err)
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
func (p *poetryPackageManager) validateVersion(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, p.executable, "--version")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.Join(ErrGetVersion, err)
	}

	// All supported versions adhere to this format.
	found := reGetPoetryVersion.FindStringSubmatch(string(stdout))
	if len(found) != 2 || !semver.IsValid("v"+found[1]) {
		return ErrGetVersion
	} else if semver.Compare(found[1], minPoetryVersion) < 0 {
		return ErrInvalidVersion
	}

	return nil
}

// Name implements PackageManager.
func (*poetryPackageManager) Name() string {
	return "poetry"
}

// Ecosystem implements PackageManager.
func (*poetryPackageManager) Ecosystem() target.Ecosystem {
	return target.EcosystemPypi
}

// normalizeCommand normalizes a pip command, so it may be executed.
func (p *poetryPackageManager) normalizeCommand(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, ErrEmptyCommand
	} else if args[0] != p.Name() {
		return nil, ErrInvalidCommand
	} else {
		args = append([]string{p.Executable()}, args[1:]...)
		return args, nil
	}
}

// RunCommand implements PackageManager.
func (p *poetryPackageManager) RunCommand(ctx context.Context, command []string) error {
	command, err := p.normalizeCommand(command)
	if err != nil {
		return err
	}

	return p.basicManager.RunCommand(ctx, command)
}

// ResolveInstallTargets implements PackageManager.
func (p *poetryPackageManager) ResolveInstallTargets(ctx context.Context, command []string) ([]target.Package, error) {
	// Check that the command should be inspected.
	command, err := p.normalizeCommand(command)
	if err != nil {
		return nil, err
	}

	var subCommand string
	for _, arg := range command[1:] {
		if !strings.HasPrefix(arg, "-") {
			subCommand = arg
			break
		}
	}

	if subCommand == "" || !slices.Contains(inspectedPoetrySubCommands, subCommand) {
		return nil, nil
	}

	// The presence of these options prevent the any command from running.
	shouldSkip := slices.ContainsFunc(command, func(s string) bool {
		switch s {
		case "-V",
			"--version",
			"-h",
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

	// Compute installation targets: new dependencies and updates/downgrades of existing ones
	return listPoetryPackages(ctx, command)
}

// listPoetryPackages list every package that would be installed by the command.
func listPoetryPackages(ctx context.Context, command []string) ([]target.Package, error) {
	command = append(command, "--dry-run")
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		logger := int_logger.GetLogger(ctx)
		logger.Info(ctx, "Encountered an error while resolving poetry installation targets")
		return nil, errors.Join(ErrDryRun, err)
	}

	var packages []target.Package

	for _, pkg := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if strings.Contains(pkg, "Skipped") {
			continue
		}

		match := reGetPoetryPackageVersion.FindStringSubmatch(pkg)
		if len(match) == 3 {
			version := match[2]
			for {
				_, newVersion, hasArrow := strings.Cut(version, " -> ")
				if hasArrow {
					version = newVersion
				} else {
					version, _, _ = strings.Cut(version, " ")
					break
				}
			}

			pkg := target.Package{
				Ecosystem: target.EcosystemPypi,
				Name:      match[1],
				Version:   version,
			}

			packages = append(packages, pkg)
		}
	}

	return packages, nil
}
