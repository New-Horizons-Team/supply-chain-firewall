package osv_advisory

import "strings"

type Severity int

const (
	SeverityNone Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *Severity) UnmarshalText(text []byte) error {
	var err error

	switch strings.ToLower(string(text)) {
	case "none":
		*s = SeverityNone
	case "low":
		*s = SeverityLow
	case "medium":
		*s = SeverityMedium
	case "high":
		*s = SeverityHigh
	case "critical":
		*s = SeverityCritical
	default:
		err = ErrInvalidSeverity
	}

	return err
}

// String implements fmt.Stringer.
func (s *Severity) String() string {
	switch *s {
	case SeverityNone:
		return "None"
	case SeverityLow:
		return "Low"
	case SeverityMedium:
		return "Medium"
	case SeverityHigh:
		return "High"
	case SeverityCritical:
		return "Critical"
	}

	panic("internal/verifiers/osv_advisory: invalid severity")
}
