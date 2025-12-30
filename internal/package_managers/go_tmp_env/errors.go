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
	// Failed to verify whether the cache dir exists.
	ErrCheckCacheDir
	// The provided path already exists but is not a directory.
	ErrInvalidCacheDir
	// Failed to inspect the directory to be copied to.
	ErrGetCacheDir
	// The cache source exists but is not a directory.
	ErrInvalidCacheSource
	// Failed to inspect the directory to be copied.
	ErrCheckSourceCache
	// Failed to copy the directory being cached.
	ErrCopyCache
	// Failed to copy the cache directory.
	ErrCopyCacheDir
	// Failed to copy the mod directory.
	ErrCopyModDir
	// Failed to purge the cache.
	ErrPurgeCache
	// Failed to create the global cache.
	ErrCreateCache
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
	case ErrCheckCacheDir:
		return "failed to verify whether the cache dir exists"
	case ErrInvalidCacheDir:
		return "the provided path already exists but is not a directory"
	case ErrGetCacheDir:
		return "failed to inspect the directory to be copied to"
	case ErrInvalidCacheSource:
		return "the cache source exists but is not a directory"
	case ErrCheckSourceCache:
		return "failed to inspect the directory to be copied"
	case ErrCopyCache:
		return "failed to copy the directory being cached"
	case ErrCopyCacheDir:
		return "failed to copy the cache directory"
	case ErrCopyModDir:
		return "failed to copy the mod directory"
	case ErrPurgeCache:
		return "failed to purge the cache"
	case ErrCreateCache:
		return "failed to create the global cache"

	default:
		return "unknown error"
	}
}
