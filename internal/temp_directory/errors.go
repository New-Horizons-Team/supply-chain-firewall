package temp_directory

type errCode int

const (
	// Failed to create the temporary directory
	ErrCreateTempDir errCode = iota
	// Failed to remove the temporary directory
	ErrRemoveTempDir
)

// Error implements the error interface for errCode.
func (err errCode) Error() string {
	switch err {
	case ErrCreateTempDir:
		return "failed to create the temporary directory"
	case ErrRemoveTempDir:
		return "failed to remove the temporary directory"
	default:
		return "Unknown error"
	}
}
