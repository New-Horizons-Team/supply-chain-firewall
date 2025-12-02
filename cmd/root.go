package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/New-Horizons-Team/supply-chain-firewall/internal/choice_flag"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/root"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "scfw",
	Short:             "A tool for preventing the installation of malicious Go, PyPI and npm packages.",
	Version:           root.GetVersion().Version,
	PersistentPreRunE: setup,
	SilenceErrors:     true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var exit *exec.ExitError

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := rootCmd.ExecuteContext(ctx)
	if errors.As(err, &exit) {
		os.Exit(exit.ExitCode())
	} else if err != nil && !errors.Is(err, ErrHalt) {
		logger := logger.DefaultLogger{}
		logger.Error(ctx, "%s", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Var(
		choice_flag.StringChoice(logger.GetLevels()),
		"log-level",
		fmt.Sprintf(
			"Desired logging level (default: %s, options: %s)",
			logger.DefaultLevel,
			strings.Join(logger.GetLevels(), ", "),
		),
	)
}

// setup configures the application before any command is executed.
func setup(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	if level, err := cmd.Flags().GetString("log-level"); err != nil {
		return err
	} else {
		logger.Configure(logger.LogLevel(level))
	}

	logger := logger.DefaultLogger{}
	logger.Info(ctx, "Starting Supply-Chain Firewall on %s", time.Now().Format("Mon Jan _2 15:04:05 2006"))

	return nil
}
