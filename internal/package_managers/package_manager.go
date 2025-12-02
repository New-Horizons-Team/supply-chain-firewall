package package_managers

import (
	"context"

	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
)

// Representation of supported package managers.
type PackageManager interface {
	// Name returns the name of the package manager, the standard fixed token by which
	// it is invoked on the command line.
	Name() string

	// Ecosystem returns the fixed package ecosystem the package manager is for.
	Ecosystem() target.Ecosystem

	// Executable returns the local filesystem path to the package manager executable.
	Executable() string

	// Run executes the given package manager command,
	// redirecting stdout and stderr to the caller's.
	RunCommand(ctx context.Context, command []string) error

	// ResolveInstallTargets resolves the package targets that would be installed if the given package
	// manager command were run (without running it).
	ResolveInstallTargets(ctx context.Context, command []string) ([]target.Package, error)
}

// getPackageManager initializes a new package manager,
// configured based on the values provided in the context.
type getPackageManager func(ctx context.Context) (PackageManager, error)

// Self-populated mapping of supported package managers.
var packageManagers = make(map[string]getPackageManager)

// GetPackageManager returns a PackageManager corresponding to the given command line provided to Supply-Chain Firewall.
func GetPackageManager(ctx context.Context, command []string) (PackageManager, error) {
	if len(command) == 0 {
		return nil, ErrMissingCommand
	}

	getter := packageManagers[command[0]]
	if getter == nil {
		return nil, ErrInvalidPackageManager
	}

	return getter(ctx)
}

type executableKey struct{}

// WithExecutable adds the executable to the context.
func WithExecutable(ctx context.Context, executable string) context.Context {
	return context.WithValue(ctx, executableKey{}, executable)
}

// getExecutable retrieves a executable from the context.
func getExecutable(ctx context.Context) string {
	executable, _ := ctx.Value(executableKey{}).(string)
	return executable
}
