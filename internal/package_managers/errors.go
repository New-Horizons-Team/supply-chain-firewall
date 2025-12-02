package package_managers

type errCode int

const (
	// Missing package manager command.
	ErrMissingCommand errCode = iota
	// Unsupported package manager.
	ErrInvalidPackageManager
	// Failed to resolve the global python executable.
	ErrMissingPython
	// Failed to resolve local executable.
	ErrMissingExecutable
	// Executable does not correspond to a regular file.
	ErrInvalidExecutable
	// Failed to get the package manager's version.
	ErrGetVersion
	// Package manager doesn't match the minimally supported version.
	ErrInvalidVersion
	// Received empty command line.
	ErrEmptyCommand
	// Received invalid command line.
	ErrInvalidCommand
	// Failed to run the command in dry-run mode.
	ErrDryRun
	// Failed to parse the provided installation targets.
	ErrParseTargets
	// Missing name for pip installation target.
	ErrPipMissingName
	// Missing version for pip installation target.
	ErrPipMissingVersion
	// Failed to resolve the global go executable.
	ErrMissingGo
	// Failed to create the temporary go environment.
	ErrGoCreateTempEnv
	// Failed to initialize the temporary go project.
	ErrGoPrepareTempProject
	// Failed to get the packages being installed.
	ErrGoGetPackage
	// Failed to list the installed packages.
	ErrGoListPackages
	// Failed to update the list of packages required by the project.
	ErrGoModTidy
	// Failed to resolve the global poetry executable.
	ErrMissingPoetry
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrMissingCommand:
		return "missing package manager command"
	case ErrInvalidPackageManager:
		return "unsupported package manager"
	case ErrMissingPython:
		return "failed to resolve the global python executable"
	case ErrMissingExecutable:
		return "failed to resolve local executable"
	case ErrInvalidExecutable:
		return "executable does not correspond to a regular file"
	case ErrGetVersion:
		return "failed to get the package manager's version"
	case ErrInvalidVersion:
		return "package manager doesn't match the minimally supported version"
	case ErrEmptyCommand:
		return "received empty command line"
	case ErrInvalidCommand:
		return "received invalid command line"
	case ErrDryRun:
		return "failed to run the command in dry-run mode"
	case ErrParseTargets:
		return "failed to parse the provided installation targets"
	case ErrPipMissingName:
		return "missing name for pip installation target"
	case ErrPipMissingVersion:
		return "missing version for pip installation target"
	case ErrMissingGo:
		return "failed to resolve the global go executable"
	case ErrGoCreateTempEnv:
		return "failed to create the temporary go environment"
	case ErrGoPrepareTempProject:
		return "failed to initialize the temporary go project"
	case ErrGoGetPackage:
		return "failed to get the packages being installed"
	case ErrGoListPackages:
		return "failed to list the installed packages"
	case ErrGoModTidy:
		return "failed to update the list of packages required by the project"
	case ErrMissingPoetry:
		return "failed to resolve the global poetry executable"

	default:
		return "unknown error"
	}
}
