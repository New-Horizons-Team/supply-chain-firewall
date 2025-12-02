package verifiers

import (
	"context"
	"runtime"
	"strings"
	"sync"

	int_logger "github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/slice_utils"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
)

// A hierarchy of severity levels for package verifier findings.
type FindingSeverity string

const (
	SeverityCritical FindingSeverity = "CRITICAL"
	SeverityWarning  FindingSeverity = "WARNING"
)

type Finding struct {
	// A concise summary of this finding.
	Summary string
	// The finding severity.
	Severity FindingSeverity
}

type FirewallVerifier interface {
	// Name returns the verifier's name.
	Name() string

	// Verify verifies the given package and returns a list of all findings for the given package
	// reported by the backing data source.
	Verify(ctx context.Context, pkg target.Package) []Finding
}

// List of registered verifiers.
var verifiers []FirewallVerifier

// GetNames returns the names of discovered package verifiers.
func GetNames() []string {
	fn := func(value FirewallVerifier) string {
		return value.Name()
	}

	return slice_utils.Map(verifiers, fn, nil)
}

// VerifyPackages verifies a set of packages against all discovered verifiers.
func VerifyPackages(ctx context.Context, packages []target.Package) map[FindingSeverity]VerificationReport {
	type payload struct {
		verifier FirewallVerifier
		pkg      target.Package
	}

	type response struct {
		findings []Finding
		pkg      target.Package
	}

	logger := int_logger.GetLogger(ctx)

	// Spawn goroutines that analyze packages with the provided verifier.
	numJobs := len(verifiers) * len(packages)
	toWork := make(chan payload, numJobs)
	done := make(chan response, numJobs)

	// Maximum number of workers is based on python's
	// concurrent.futures.ThreadPoolExecutor's behaviour.
	numWorkers := max(runtime.NumCPU(), 1)
	numWorkers = min(32, numWorkers+4)

	var verifyWg sync.WaitGroup
	for range numWorkers {
		verifyWg.Add(1)

		go func(ctx context.Context, wg *sync.WaitGroup, recv <-chan payload, send chan<- response) {
			defer wg.Done()

			for data := range recv {
				findings := data.verifier.Verify(ctx, data.pkg)

				if len(findings) > 0 {
					send <- response{findings, data.pkg}
					logger.Info(ctx, "Verifier %s had findings for package %s", data.verifier.Name(), data.pkg)
				} else {
					logger.Info(ctx, "Verifier %s had no findings for package %s", data.verifier.Name(), data.pkg)
				}
			}
		}(ctx, &verifyWg, toWork, done)
	}

	// Spawn a goroutine to accumulate the generated reports.
	reports := make(map[FindingSeverity]VerificationReport)

	var accumulateWg sync.WaitGroup
	accumulateWg.Add(1)

	go func(wg *sync.WaitGroup, recv <-chan response, reports map[FindingSeverity]VerificationReport) {
		defer wg.Done()

		for data := range recv {
			for _, finding := range data.findings {
				report := reports[finding.Severity]
				if report == nil {
					report = NewVerificationReport()
					reports[finding.Severity] = report
				}

				report.Insert(data.pkg, finding.Summary)
			}
		}
	}(&accumulateWg, done, reports)

	// Normalize the package names.
	nomalizedPackages := make([]target.Package, len(packages))
	for i, pkg := range packages {
		nomalizedPackages[i] = pkg
		nomalizedPackages[i].Name = strings.ToLower(pkg.Name)
	}

	// Enqueue every job and wait until its done.
	for _, verifier := range verifiers {
		for _, pkg := range nomalizedPackages {
			toWork <- payload{verifier, pkg}
		}
	}

	close(toWork)
	verifyWg.Wait()

	close(done)
	accumulateWg.Wait()

	logger.Info(ctx, "Verification of packages complete")
	return reports
}
