package go_tmp_env

type errCode int

const (
	// Failed to create the temporary directory.
	ErrCreateTempDir errCode = iota
	// Failed to detect the parent go project.
	ErrGetOriginalDir
	// Failed to get the absolute path to the parent go project.
	ErrExpandOriginalDir
	// Failed to create the temporary go packages directory.
	ErrCreateGoDir
	// Failed to create the temporary dummy project directory.
	ErrCreateDryRunDir
	// Failed to create the temporary cache directory.
	ErrCreateCache
	// Failed to create the temporary mod cache directory.
	ErrCreateModCache
	// Failed to detect the original GOPATH.
	ErrGetGoPath
	// Failed to open the file being copied.
	ErrCopyOpenSource
	// Failed to open the file that will have date copied into it.
	ErrCopyOpenDestination
	// Failed to copy the contents of the file.
	ErrCopyContent
	// Failed to restore the original 'go.mod'.
	ErrRestoreGoMod
	// Failed to restore the original 'go.sum'.
	ErrRestoreGoSum
	// Failed to remove the created 'go.sum'.
	ErrRemoveGoSum
	// Failed to execute the command in the temporary environment.
	ErrRunCommand
	// Cannot duplicate the original 'go.mod' was it wasn't found.
	ErrGoModNotFound
	// Failed to duplicate 'go.mod'.
	ErrCopyGoMod
	// Failed to duplicate 'go.sum'.
	ErrCopyGoSum
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrCreateTempDir:
		return "failed to create the temporary directory"
	case ErrExpandOriginalDir:
		return "failed to get the absolute path to the parent go project"
	case ErrGetOriginalDir:
		return "failed to detect the parent go project"
	case ErrCreateGoDir:
		return "failed to create the temporary go packages directory"
	case ErrCreateDryRunDir:
		return "failed to create the temporary dummy project directory"
	case ErrCreateCache:
		return "failed to create the temporary cache directory"
	case ErrCreateModCache:
		return "failed to create the temporary mod cache directory"
	case ErrGetGoPath:
		return "failed to detect the original GOPATH"
	case ErrCopyOpenSource:
		return "failed to open the file being copied"
	case ErrCopyOpenDestination:
		return "failed to open the file that will have date copied into it"
	case ErrCopyContent:
		return "failed to copy the contents of the file"
	case ErrRestoreGoMod:
		return "failed to restore the original 'go.mod'"
	case ErrRestoreGoSum:
		return "failed to restore the original 'go.sum'"
	case ErrRemoveGoSum:
		return "failed to remove the created 'go.sum'"
	case ErrRunCommand:
		return "failed to execute the command in the temporary environment"
	case ErrGoModNotFound:
		return "cannot duplicate the original 'go.mod' was it wasn't found"
	case ErrCopyGoMod:
		return "failed to duplicate 'go.mod'"
	case ErrCopyGoSum:
		return "failed to duplicate 'go.sum'"

	default:
		return "unknown error"
	}
}
