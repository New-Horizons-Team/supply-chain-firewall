package temp_directory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTempDir checks that the temporary directory works as expected.
func TestTempDir(t *testing.T) {
	dir, err := NewTempDir("tmp_dir_")
	require.NoError(t, err, "failed to create the temporary directory")
	defer func() {
		_ = dir.Close()
	}()

	stat, err := os.Stat(dir.GetPath())
	require.NoError(t, err, "failed to check the temporary directory")
	require.True(t, stat.IsDir(), "invalid temporary directory")

	filename := dir.GetPath() + "foo"
	err = os.WriteFile(filename, []byte("hello"), 0644)
	require.NoError(t, err, "failed to check the temporary directory")

	stat, err = os.Stat(filename)
	require.NoError(t, err, "failed to check the file within the temporary directory")
	require.False(t, stat.IsDir(), "invalid file within the temporary directory")

	err = dir.Close()
	require.NoError(t, err, "failed to remove the temporary directory")

	_, err = os.Stat(dir.GetPath())
	require.Error(t, err, "temporary directory still exists after removal")
}
