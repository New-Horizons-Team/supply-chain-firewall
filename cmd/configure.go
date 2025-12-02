package cmd

import (
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/interactive"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/printer"
	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/configure"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the environment for using Supply-Chain Firewall.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		remove, err := cmd.Flags().GetBool("remove")
		if err != nil {
			return err
		}
		aliasPip, err := cmd.Flags().GetBool("alias-pip")
		if err != nil {
			return err
		}
		aliasPoetry, err := cmd.Flags().GetBool("alias-poetry")
		if err != nil {
			return err
		}
		aliasGo, err := cmd.Flags().GetBool("alias-go")
		if err != nil {
			return err
		}

		if !remove && !aliasPip && !aliasPoetry && !aliasGo {
			cfg, err := interactive.GetConfig(ctx)
			if err != nil {
				return err
			}

			aliasPip = cfg.Pip
			aliasPoetry = cfg.Poetry
			aliasGo = cfg.Golang
		}

		return configure.RunConfigure(ctx, aliasPip, aliasPoetry, aliasGo, remove, printer.DefaultPrinter{}, configure.WithShowDiff(interactive.Confirm))
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.PersistentFlags().BoolP("remove", "r", false, "Remove all Supply-Chain Firewall-managed configuration")
	configureCmd.PersistentFlags().Bool("alias-pip", false, "Add shell aliases to always run pip commands through Supply-Chain Firewall")
	configureCmd.PersistentFlags().Bool("alias-poetry", false, "Add shell aliases to always run Poetry commands through Supply-Chain Firewall")
	configureCmd.PersistentFlags().Bool("alias-go", false, "Add shell aliases to always run go commands through Supply-Chain Firewall")
}
