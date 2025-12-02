package root

import _ "embed"

//go:embed VERSION
var version string

type ScfwVersion struct {
	// The application's version.
	Version string
}

// GetVersion retrieves version information about scfw's.
func GetVersion() ScfwVersion {
	return ScfwVersion{
		Version: version,
	}
}
