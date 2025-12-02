package temp_directory

import (
	"errors"
	"os"
)

type tempDir struct {
	// The temporary directory's path.
	path string
	// Whether the temporary directory was removed.
	removed bool
}

type TempDir interface {
	// GetPath retrieves the temporary directory
	// represented by this object.
	GetPath() string

	// Close implements io.Closer,
	// removing the temporary directory.
	Close() error
}

// NewTempDir creates a new temporary directory
// with the provided prefix.
func NewTempDir(prefix string) (TempDir, error) {
	name, err := os.MkdirTemp("/tmp", prefix)

	if err != nil {
		return nil, errors.Join(ErrCreateTempDir, err)
	}

	if name[len(name)-1] != '/' {
		name += "/"
	}

	tempDir := &tempDir{
		path: name,
	}

	return tempDir, nil
}

// GetPath implements TempDir
func (t *tempDir) GetPath() string {
	if t.removed {
		return ""
	}
	return t.path
}

// Close implements TempDir and io.Closer.
func (t *tempDir) Close() error {
	if t.removed {
		return nil
	}

	err := os.RemoveAll(t.path)
	if err != nil {
		return errors.Join(ErrRemoveTempDir, err)
	}
	t.removed = true

	return nil
}
