package package_managers

import (
	"context"
	"os"
	"os/exec"
)

type basicManager struct {
	// Path to the manager's executable binary, if set.
	executable string
}

// Executable partially implements PackageManager.
func (b *basicManager) Executable() string {
	return b.executable
}

// init inializes this manager with values from the context.
func (b *basicManager) init(ctx context.Context) {
	b.executable = getExecutable(ctx)
}

// RunCommand partially implements PackageManager.
func (b *basicManager) RunCommand(ctx context.Context, command []string) error {
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
