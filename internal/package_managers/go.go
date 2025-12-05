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
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/package_managers/go_tmp_env"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
	"golang.org/x/mod/semver"
)

const minGoVersion = "1.17.0"

var reGetGoVersion = regexp.MustCompile(`.*go(\d*(?:\.\d+)*).*`)

var inspectedGoSubCommands = []string{"build", "generate", "get", "install", "mod", "run"}

var inspectedGoModCommands = []string{"download", "graph", "tidy", "verify", "why"}

var inspectedGoNoPackageCommands = []string{"build", "get", "install"}

const dryRunProject = "localhost/dry_run"

type goPackageManager struct {
	basicManager
}

// init registers a getter for go.
func init() {
	packageManagers["go"] = func(ctx context.Context) (PackageManager, error) {
		pkgMngr := &goPackageManager{}

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
func (g *goPackageManager) init(ctx context.Context) error {
	g.basicManager.init(ctx)

	if g.executable == "" {
		var err error

		g.executable, err = exec.LookPath("go")
		if err != nil {
			return errors.Join(ErrMissingGo, err)
		}
	}

	if g.executable == "" {
		return ErrMissingExecutable
	}
	if stat, err := os.Stat(g.executable); err != nil {
		return errors.Join(ErrInvalidExecutable, err)
	} else if !stat.Mode().IsRegular() {
		return ErrInvalidExecutable
	}

	return nil
}

// validateVersion checks if the package manager version is supported.
func (g *goPackageManager) validateVersion(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, g.executable, "version")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.Join(ErrGetVersion, err)
	}

	// All supported versions adhere to this format.
	found := reGetGoVersion.FindStringSubmatch(string(stdout))
	if len(found) != 2 || !semver.IsValid("v"+found[1]) {
		return ErrGetVersion
	} else if semver.Compare(found[1], minGoVersion) < 0 {
		return ErrInvalidVersion
	}

	return nil
}

// Name implements PackageManager.
func (*goPackageManager) Name() string {
	return "go"
}

// Ecosystem implements PackageManager.
func (*goPackageManager) Ecosystem() target.Ecosystem {
	return target.EcosystemGo
}

// normalizeCommand normalizes a pip command, so it may be executed.
func (g *goPackageManager) normalizeCommand(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, ErrEmptyCommand
	} else if args[0] != g.Name() {
		return nil, ErrInvalidCommand
	} else {
		args = append([]string{g.Executable()}, args[1:]...)
		return args, nil
	}
}

// RunCommand implements PackageManager.
func (g *goPackageManager) RunCommand(ctx context.Context, command []string) error {
	command, err := g.normalizeCommand(command)
	if err != nil {
		return err
	}

	return g.basicManager.RunCommand(ctx, command)
}

// ResolveInstallTargets implements PackageManager.
func (g *goPackageManager) ResolveInstallTargets(ctx context.Context, command []string) ([]target.Package, error) {
	// Check that the command should be inspected.
	command, err := g.normalizeCommand(command)
	if err != nil {
		return nil, err
	}

	if len(command) < 2 || !slices.Contains(inspectedGoSubCommands, command[1]) ||
		(len(command) > 2 && command[1] == "mod" && !slices.Contains(inspectedGoModCommands, command[2])) {

		return nil, nil
	}

	// The presence of these options prevent any command from running.
	shouldSkip := slices.ContainsFunc(command, func(s string) bool {
		switch s {
		case "-h",
			"-help":

			return true
		default:
			return false
		}
	})
	if shouldSkip {
		return nil, nil
	}

	// Compute installation targets: new dependencies and updates/downgrades of existing ones
	localPackages, remotePackages, getFlags := extractGoTargets(command)

	logger := int_logger.GetLogger(ctx)
	logger.Debug(ctx, "Local packages to be inspected: %v", localPackages)
	logger.Debug(ctx, "Remote packages to be inspected: %v", remotePackages)

	tmp, err := go_tmp_env.New(ctx, g.executable)
	if err != nil {
		return nil, errors.Join(ErrGoCreateTempEnv, err)
	}
	defer func() {
		tmpErr := tmp.Close()
		if tmpErr != nil {
			logger.Warn(
				ctx,
				"Failed to remove the temporary go directory '%s': %s.\n"+
					"Please remove it manually.",
				tmp.GetRoot(),
				tmpErr.Error(),
			)
		}
	}()

	foundPackages := make(map[target.Package]struct{})

	if len(remotePackages) > 0 {
		// Create a temporary project and retrieve what would be installed in it.
		_, err := tmp.Run(ctx, []string{"mod", "init", dryRunProject}, false)
		if err != nil {
			return nil, errors.Join(ErrGoPrepareTempProject, err)
		}

		getArgs := append([]string{"get"}, getFlags...)
		getArgs = append(getArgs, remotePackages...)
		_, err = tmp.Run(ctx, getArgs, false)
		if err != nil {
			return nil, errors.Join(ErrGoGetPackage, err)
		}

		err = listGoPackages(ctx, tmp, false, foundPackages)
		if err != nil {
			return nil, err
		}
	}

	isGet := command[1] == "get"
	isTidy := len(command) > 2 && command[1] == "mod" && command[2] == "tidy"

	if isTidy || len(localPackages) > 0 {
		if isTidy {
			_, err = tmp.Run(ctx, []string{"mod", "tidy"}, true)
			if err != nil {
				return nil, errors.Join(ErrGoModTidy, err)
			}
		} else if isGet && len(getFlags) > 0 {
			getArgs := append([]string{"get"}, getFlags...)
			_, err = tmp.Run(ctx, getArgs, true)
			if err != nil {
				return nil, errors.Join(ErrGoGetPackage, err)
			}
		}

		err = listGoPackages(ctx, tmp, true, foundPackages)
		if err != nil {
			return nil, err
		}
	}

	var packages []target.Package
	for pkg := range foundPackages {
		packages = append(packages, pkg)
	}
	return packages, nil
}

// extractGoTargets list which local/remote packages are to be installed by the provided command.
func extractGoTargets(command []string) (localPackages []string, remotePackages []string, getFlags []string) {
	if len(command) < 2 {
		return
	}

	isGet := command[1] == "get"
	isTidy := len(command) > 2 && command[1] == "mod" && command[2] == "tidy"
	targetPackages := command[2:]
	if isTidy {
		targetPackages = targetPackages[1:]
	}

	var nonFlagCount int

	for _, pkg := range targetPackages {
		if strings.HasPrefix(pkg, "-") {
			if isGet && pkg == "-t" || pkg == "-u" || strings.HasPrefix(pkg, "-u=") {
				getFlags = append(getFlags, pkg)
			}

			continue
		}
		nonFlagCount += 1

		_, err := os.Stat(pkg)
		isLocal := err == nil

		if isLocal || len(strings.Split(pkg, "/")) == 1 {
			localPackages = append(localPackages, pkg)
		} else {
			remotePackages = append(remotePackages, pkg)
		}
	}
	if nonFlagCount == 0 && slices.Contains(inspectedGoNoPackageCommands, command[1]) {
		localPackages = append(localPackages, ".")
	}

	return
}

// listGoPackages list every package either in the temporary environment or in the local directory
// and add them to foundPackages.
func listGoPackages(ctx context.Context, tmp go_tmp_env.TempGoEnvironment, local bool, foundPackages map[target.Package]struct{}) error {
	stdout, err := tmp.Run(ctx, []string{"list", "-m", "all"}, local)
	if err != nil {
		return errors.Join(ErrGoListPackages, err)
	}

	for _, pkg := range strings.Split(strings.TrimSpace(stdout), "\n") {
		components := strings.Split(strings.TrimSpace(pkg), " ")

		if len(components) == 2 && components[0] != dryRunProject {
			pkg := target.Package{
				Ecosystem: target.EcosystemGo,
				Name:      components[0],
				Version:   components[1],
			}

			foundPackages[pkg] = struct{}{}
		}
	}

	return nil
}
