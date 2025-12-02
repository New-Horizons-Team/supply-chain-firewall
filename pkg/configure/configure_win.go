//go:build windows

package configure

import "context"

// RunConfigure configures the environment for use with the supply-chain firewall.
func RunConfigure(ctx context.Context, aliasPip, aliasPoetry, aliasGo, remove bool, opts ...WithOptions) error {
	return ErrWindows
}
