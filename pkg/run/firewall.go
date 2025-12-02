package run

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	int_logger "github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/package_managers"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/slice_utils"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/verifiers"
	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/printer"
)

// Function that receives a prompt for querying the user for a yes/no answer.
type Confirm = func(ctx context.Context, prompt string) bool

// RunFirewall executes a package manager command throught scfw.
func RunFirewall(ctx context.Context, dryRun, automation bool, executable string, command []string, confirm Confirm, printer printer.Printer, logger logger.Logger) error {
	ctx = int_logger.WithLogger(ctx, logger)

	logger.Info(ctx, "Command: %s", strings.Join(command, " "))

	if executable != "" {
		ctx = package_managers.WithExecutable(ctx, executable)
	}

	manager, err := package_managers.GetPackageManager(ctx, command)
	if err != nil {
		return err
	}

	targets, err := manager.ResolveInstallTargets(ctx, command)
	if err != nil {
		var exitErr *exec.ExitError

		printer.PrintWarningReports(ctx, "Failed to list install targets")

		if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
			stderr, _ := strings.CutSuffix(string(exitErr.Stderr), "\n")
			printer.Print(ctx, stderr)
		} else {
			printer.PrintWarningReports(ctx, err.Error())
		}

		return err
	}

	if len(targets) > 0 {
		logger.Info(ctx, "Command would install: %s", strings.Join(packagesToString(targets), ", "))

		logger.Info(
			ctx,
			"Using package verifiers: %s",
			strings.Join(verifiers.GetNames(), ", "),
		)

		reports := verifiers.VerifyPackages(ctx, targets)

		if criticalReport := reports[verifiers.SeverityCritical]; criticalReport != nil {
			printer.PrintCriticalReports(ctx, criticalReport.String())
			printer.Print(ctx, "The installation request was blocked. No changes have been made.")

			if automation {
				return ErrCriticalReport
			}
			return nil
		}

		if warningReport := reports[verifiers.SeverityWarning]; warningReport != nil {
			printer.PrintWarningReports(ctx, warningReport.String())

			if !automation && !confirm(ctx, "Proceed with installation?") {
				printer.Print(ctx, "The installation request was aborted. No changes have been made.")
				return nil
			}
		}
	} else {
		logger.Info(ctx, "No packages to be installed by command")
	}

	if dryRun {
		logger.Info(ctx, "Firewall dry-run mode enabled: command will not be run")
		printer.Print(ctx, "Dry-run: exiting without running command.")
	} else {
		return manager.RunCommand(ctx, command)
	}

	return nil
}

// packagesToString converts a list of packages to a list of strings.
func packagesToString(packages []target.Package) []string {
	fn := func(p target.Package) string {
		return p.String()
	}

	return slice_utils.Map(packages, fn, nil)
}
