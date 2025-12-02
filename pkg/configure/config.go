package configure

import (
	"context"
	"fmt"
	"strings"
)

// Function that receives a prompt for querying the user for a yes/no answer.
type Confirm = func(ctx context.Context, prompt string) bool

type comparableConfig struct {
	command string
	pip     bool
	poetry  bool
	golang  bool
}

type config struct {
	comparableConfig

	diff    bool
	confirm Confirm
}

// String implements fmt.Stringer.
func (cfg config) String() string {
	var config []string

	if cfg.command == "" {
		cfg.command = "scfw"
	}

	if cfg.pip {
		config = append(config, fmt.Sprintf(`alias pip="%s run pip --"`, cfg.command))
	}
	if cfg.poetry {
		config = append(config, fmt.Sprintf(`alias poetry="%s run poetry --"`, cfg.command))
	}
	if cfg.golang {
		config = append(config, fmt.Sprintf(`alias go="%s run go --"`, cfg.command))
	}

	// Prepend a '\n' to the returned string.
	if len(config) > 0 {
		config = append([]string{""}, config...)
	}

	return strings.Join(config, "\n")
}

type WithOptions func(cfg *config)

// WithSetCommand changes the aliased command from scfw to whatever was provided.
func WithSetCommand(command string) WithOptions {
	return func(cfg *config) {
		cfg.command = command
	}
}
