package configure

type errCode int

const (
	// Configure isn't currently supported on Windows.
	ErrWindows errCode = iota
	// Failed to retrieve the configurations interactively.
	ErrQueryConfig
	// Failed to list the user's home directory.
	ErrGetHome
	// Failed to retrieve the configuration file's mode.
	ErrGetConfigFileMode
	// Failed to read the configuration file.
	ErrConfigFile
	// Failed to create the temporary file to overwriting the configuration file.
	ErrCreateTemporaryFile
	// Failed to update the configuration file.
	ErrUpdateConfigFile
	// The modifications were refused by the user.
	ErrNotApplied
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrWindows:
		return "configure isn't currently supported on Windows."
	case ErrQueryConfig:
		return "failed to retrieve the configurations interactively"
	case ErrGetHome:
		return "failed to list the user's home directory"
	case ErrGetConfigFileMode:
		return "failed to retrieve the configuration file's mode"
	case ErrConfigFile:
		return "failed to read the configuration file"
	case ErrCreateTemporaryFile:
		return "failed to create the temporary file to overwriting the configuration file"
	case ErrUpdateConfigFile:
		return "failed to update the configuration file"
	case ErrNotApplied:
		return "the modifications were refused by the user"

	default:
		return "unknown error"
	}
}
