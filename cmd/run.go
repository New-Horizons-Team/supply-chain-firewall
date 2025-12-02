package cmd

import (
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/interactive"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/printer"
	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/run"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a package manager command through Supply-Chain Firewall.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = len(args) > 0

		if len(args) == 0 {
			cmd.SilenceErrors = true
			return ErrHalt
		}

		return setup(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		automation, err := cmd.Flags().GetBool("automation")
		if err != nil {
			return err
		}
		executable, err := cmd.Flags().GetString("executable")
		if err != nil {
			return err
		}

		return run.RunFirewall(ctx, dryRun, automation, executable, args, interactive.Confirm, printer.DefaultPrinter{}, logger.DefaultLogger{})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().Bool("dry-run", false, "Verify any installation targets but do not run the package manager command")
	runCmd.PersistentFlags().String("executable", "", "Go, Python, or npm executable to use for running commands (default: environmentally determined)")
	runCmd.PersistentFlags().Bool("automation", false, "Change behaviour to be called by automations. Namely, exit with status code 1 if a verifier blocks the command, and automatically proceed regardless of report warnings")
}
