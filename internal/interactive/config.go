package interactive

import (
	"context"
	"errors"

	"github.com/charmbracelet/huh"
)

type Config struct {
	// Whether pip shall be aliased.
	Pip bool
	// Whether poetry shall be aliased.
	Poetry bool
	// Whether go shall be aliased.
	Golang bool
}

// GetConfig queries the configurations to the user.
func GetConfig(ctx context.Context) (Config, error) {
	var cfg Config

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("iFood Supply-Chain Firewall").
				Description(
					"Thank you for using scfw, the iFood Supply-Chain Firewall!\n\n"+
						"scfw is a tool for preventing the installation of malicious PyPI and go packages.\n\n"+
						"This script will walk you through setting up your environment to get the most out\n"+
						"of scfw. You can rerun this script at any time.",
				).
				Next(true),

			huh.NewConfirm().
				Title("Would you like to set a shell alias to run all pip commands through the firewall?").
				Value(&cfg.Pip),

			huh.NewConfirm().
				Title("Would you like to set a shell alias to run all Poetry commands through the firewall?").
				Value(&cfg.Poetry),

			huh.NewConfirm().
				Title("Would you like to set a shell alias to run all go commands through the firewall?").
				Value(&cfg.Golang),
		),
	)

	err := form.
		WithAccessible(true).
		RunWithContext(ctx)
	if err != nil {
		return cfg, errors.Join(ErrQueryConfig, err)
	}

	return cfg, nil
}
