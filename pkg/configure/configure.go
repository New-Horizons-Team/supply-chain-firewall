//go:build !windows

package configure

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/New-Horizons-Team/supply-chain-firewall/internal/diff"
	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/printer"
)

const (
	blockStart = "\n# BEGIN SCFW MANAGED BLOCK"
	blockEnd   = "\n# END SCFW MANAGED BLOCK\n"
)

var (
	// List of configuration files supported by scfw.
	configFiles = []string{".bashrc", ".zshrc"}

	// Regex for capturing the current configuration.
	// Flag 's' captures in multi-line mode
	oldConfig = regexp.MustCompile(blockStart + "(?s:.*?)" + blockEnd)
)

// RunConfigure configures the environment for use with the supply-chain firewall.
func RunConfigure(ctx context.Context, aliasPip, aliasPoetry, aliasGo, remove bool, printer printer.Printer, opts ...WithOptions) error {
	var cfg, emptyCfg config
	var skipped bool

	if !remove {
		cfg.pip = aliasPip
		cfg.poetry = aliasPoetry
		cfg.golang = aliasGo
	}

	isEmpty := emptyCfg.comparableConfig == cfg.comparableConfig

	for _, fn := range opts {
		if fn != nil {
			fn(&cfg)
		}
	}

	files, err := listConfigFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		skipped = false

		err = cfg.updateConfigFile(ctx, file, printer)
		if errors.Is(err, ErrNotApplied) {
			skipped = true
		} else if err != nil {
			return err
		}
	}

	if skipped {
		printer.Print(ctx, "No changes have been made.")
	} else if isEmpty {
		printer.Print(
			ctx,
			"All Supply-Chain Firewall-managed configuration has been removed from your environment."+
				"\n\nPost-removal tasks:"+
				"\n* Update your current shell environment by sourcing from your .bashrc/.zshrc file.",
		)
	} else {
		printer.Print(
			ctx,
			"The environment was successfully configured for Supply-Chain Firewall."+
				"\n\nPost-configuration tasks:"+
				"\n* Update your current shell environment by sourcing from your .bashrc/.zshrc file."+
				"\n\nGood luck!",
		)
	}

	return nil
}

// listConfigFiles tries to list every config file available in the user's home directory.
func listConfigFiles() ([]string, error) {
	var ret []string

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Join(ErrGetHome, err)
	}

	for _, file := range configFiles {
		file = filepath.Join(home, file)
		if _, err := os.Stat(file); err == nil {
			ret = append(ret, file)
		}
	}

	return ret, nil
}

// updateConfigFile updates the configuration file with the provided scfw config.
func (cfg *config) updateConfigFile(ctx context.Context, file string, printer printer.Printer) error {
	original, err := os.ReadFile(file)
	if err != nil {
		return errors.Join(ErrConfigFile, err)
	}

	var data []byte

	config := cfg.String()
	if config != "" {
		data = []byte(blockStart + config + blockEnd)
	}

	updated := oldConfig.ReplaceAll(original, data)
	if bytes.Equal(original, updated) && !bytes.Contains(original, data) {
		updated = append(updated, data...)
	}

	if !bytes.Equal(original, updated) {
		if cfg.diff && cfg.confirm != nil {
			d := diff.Unified(file, file, string(original), string(updated))

			// Ensure the unified diff doesn't end in a newline.
			d, _ = strings.CutSuffix(d, "\r\n")
			d, _ = strings.CutSuffix(d, "\n")
			printer.PrintDiff(ctx, d)

			if !cfg.confirm(ctx, "Would you like to apply these changes?") {
				return ErrNotApplied
			}
		}

		err = overwriteFileSafe(file, updated)
		if err != nil {
			return err
		}
	}

	return nil
}

// overwriteFileSafe overwrites file with data, ensuring that the original contents are kept on error.
func overwriteFileSafe(file string, data []byte) (err error) {
	var info os.FileInfo
	var out *os.File

	info, err = os.Stat(file)
	if err != nil {
		return errors.Join(ErrGetConfigFileMode, err)
	}
	mode := info.Mode()

	// Create a temporary file on the same directory
	// (to avoid issues with moving across volumes/devices/partitions).
	dir := filepath.Dir(file)
	out, err = os.CreateTemp(dir, ".tmp-config.")
	if err != nil {
		return errors.Join(ErrCreateTemporaryFile, err)
	}
	defer func() {
		_ = out.Close()
		if err != nil {
			_ = os.Remove(out.Name())
		}
	}()

	err = out.Chmod(mode)
	if err != nil {
		return errors.Join(ErrCreateTemporaryFile, err)
	}

	err = out.Close()
	if err != nil {
		return errors.Join(ErrCreateTemporaryFile, err)
	}

	// Replace the config file with the temporary file.
	err = os.WriteFile(out.Name(), data, mode)
	if err != nil {
		return errors.Join(ErrUpdateConfigFile, err)
	}

	err = os.Rename(out.Name(), file)
	if err != nil {
		return errors.Join(ErrUpdateConfigFile, err)
	}

	return nil
}
