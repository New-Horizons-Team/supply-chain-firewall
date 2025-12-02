package interactive

import (
	"context"

	"github.com/charmbracelet/huh"
)

// Confirm queries whether the user want to accept something.
func Confirm(ctx context.Context, prompt string) bool {
	var value bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Value(&value),
		),
	)

	err := form.
		WithAccessible(true).
		RunWithContext(ctx)

	return err == nil && value
}
