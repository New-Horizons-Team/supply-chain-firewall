package osv_advisory

type errCode int

const (
	// Invalid severity.
	ErrInvalidSeverity errCode = iota
	// Encountered OSV advisory with missing ID field.
	ErrMissingAdvisoryID
	// Encountered an invalid score for a Ubuntu severity.
	ErrInvalidUbuntuScore
	// Encountered an invalid CVSS type.
	ErrInvalidCvssType
	// Failed to parse the CVSS score.
	ErrParseScore
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrInvalidSeverity:
		return "invalid severity"
	case ErrMissingAdvisoryID:
		return "encountered OSV advisory with missing ID field"
	case ErrInvalidUbuntuScore:
		return "encountered an invalid score for a Ubuntu severity"
	case ErrInvalidCvssType:
		return "encountered an invalid CVSS type"
	case ErrParseScore:
		return "failed to parse the CVSS score"

	default:
		return "unknown error"
	}
}
