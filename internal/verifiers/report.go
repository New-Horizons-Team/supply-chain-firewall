package verifiers

import (
	"fmt"
	"strings"

	"github.com/New-Horizons-Team/supply-chain-firewall/internal/slice_utils"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
)

// A structured report containing findings resulting from package verification.
type VerificationReport interface {
	// Count returns the number of entries in the report.
	Count() int

	// String implements fmt.Stringer.
	String() string

	// Get the findings for the given package.
	Get(pkg target.Package) []string

	// Insert the given package and finding into the report.
	Insert(pkg target.Package, finding string)

	// Packages returns the list of packages in the report.
	Packages() []target.Package
}

type verificationReport struct {
	// Map each package to the list of finding on it.
	reports map[target.Package][]string
}

// NewVerificationReport initializes a new VerificationReport.
func NewVerificationReport() VerificationReport {
	return &verificationReport{
		reports: make(map[target.Package][]string),
	}
}

// Count implements VerificationReport.
func (v *verificationReport) Count() int {
	return len(v.reports)
}

// String implements VerificationReport.
func (v *verificationReport) String() string {
	var reports []string

	fn := func(finding string) string {
		var lines []string

		for linenum, line := range strings.Split(finding, "\n") {
			if linenum == 0 {
				lines = append(lines, fmt.Sprintf("  - %s", line))
			} else {
				lines = append(lines, fmt.Sprintf("    %s", line))
			}
		}

		return strings.Join(lines, "\n")
	}

	for pkg, findings := range v.reports {
		reports = append(
			reports,
			fmt.Sprintf(
				"Package %s:\n%s",
				pkg,
				strings.Join(slice_utils.Map(findings, fn, nil), "\n"),
			),
		)
	}

	return strings.Join(reports, "\n")
}

// Get implements VerificationReport.
func (v *verificationReport) Get(pkg target.Package) []string {
	return v.reports[pkg]
}

// Insert implements VerificationReport.
func (v *verificationReport) Insert(pkg target.Package, finding string) {
	v.reports[pkg] = append(v.reports[pkg], finding)
}

// Packages implements VerificationReport.
func (v *verificationReport) Packages() []target.Package {
	var ret []target.Package

	for pkg := range v.reports {
		ret = append(ret, pkg)
	}

	return ret
}
