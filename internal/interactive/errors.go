package interactive

type errCode int

const (
	// Failed to retrieve the configurations interactively.
	ErrQueryConfig errCode = iota
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrQueryConfig:
		return "failed to retrieve the configurations interactively"

	default:
		return "unknown error"
	}
}
