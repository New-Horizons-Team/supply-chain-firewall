package printer

import "context"

// Printer provides a unified interface for printing information for the user.
// Each call to a function within this interface is intended to be a new line,
// but the callee is responsible with adding the terminating newline ('\n')!
type Printer interface {
	// PrintDiff prints a unified diff to stdout.
	PrintDiff(ctx context.Context, str string)

	// PrintCriticalReports prints one or more critical reports to stdout.
	// Note that critical reports block the installation of packages.
	PrintCriticalReports(ctx context.Context, str string)

	// PrintWarningReports prints one or more warning reports to stdout.
	// Note that warning reports first ask for user confirmation before
	// blocking the installation of packages.
	PrintWarningReports(ctx context.Context, str string)

	// Print prints the provided string to stdout.
	Print(ctx context.Context, str string)
}
