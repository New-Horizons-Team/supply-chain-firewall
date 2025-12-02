package verifiers

type errCode int

const (
	// Invalid ecosystem.
	ErrInvalidEcosystem errCode = iota
	// Failed to download the ecosystem's DataDog manifest file.
	ErrDownloadDataDogManifest
	// Failed to decode the ecosystem's Datadog manifest file.
	ErrDecodeDataDogManifest
)

// Error implements error.
func (e errCode) Error() string {
	switch e {
	case ErrInvalidEcosystem:
		return "invalid ecosystem"
	case ErrDownloadDataDogManifest:
		return "failed to download the ecosystem's DataDog manifest file"
	case ErrDecodeDataDogManifest:
		return "failed to decode the ecosystem's Datadog manifest file"

	default:
		return "unknown error"
	}
}
