package target

import (
	"fmt"
)

type Package struct {
	// The package's ecosystem.
	Ecosystem Ecosystem
	// The package's name.
	Name string
	// The package's version string.
	Version string
}

// String implements fmt.Stringer.
func (p Package) String() string {
	switch p.Ecosystem {
	case EcosystemNpm,
		EcosystemGo:

		return fmt.Sprintf("%s@%s", p.Name, p.Version)
	case EcosystemPypi:
		return fmt.Sprintf("%s-%s", p.Name, p.Version)
	}

	panic("internal/package: unsupported ecosystem")
}
