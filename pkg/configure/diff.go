package configure

// WithShowDiff causes RunConfigure to show what was modified in the config file
// and to ask the user for confirmation before commiting it.
func WithShowDiff(confirm Confirm) WithOptions {
	return func(cfg *config) {
		cfg.diff = true
		cfg.confirm = confirm
	}
}
