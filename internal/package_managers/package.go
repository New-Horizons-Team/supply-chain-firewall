package package_managers

import (
	"fmt"
)

type Package struct {
	// The package's ecosystem.
	ecosystem Ecosystem
	// The package's name.
	name string
	// The package's version string.
	version string
}

// String implements fmt.Stringer.
func (p *Package) String() string {
	switch p.ecosystem {
	case EcosystemNpm,
		EcosystemGo:

		return fmt.Sprintf("%s@%s", p.name, p.version)
	case EcosystemPypi:
		return fmt.Sprintf("%s-%s", p.name, p.version)
	}

	panic("internal/package_managers: unsupported ecosystem")
}
