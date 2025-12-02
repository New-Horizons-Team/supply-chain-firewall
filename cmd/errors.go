package cmd

type errCode int

const (
	// Signals the application that it should stop gracefully.
	ErrHalt errCode = iota
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrHalt:
		return "signals the application that it should stop gracefully"

	default:
		return "unknown error"
	}
}
