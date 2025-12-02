package run

type errCode int

const (
	// A critical report has been found; execution will halt.
	ErrCriticalReport errCode = iota
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrCriticalReport:
		return "a critical report has been found; execution will halt"

	default:
		return "unknown error"
	}
}
