package go_tmp_env

import (
	"context"
	"testing"
	"os"
	"os/exec"
	"strings"

	"github.com/stretchr/testify/require"
)

func TestCreateCache(t *testing.T) {
	defer func() {
		_ = PurgeCache()
	}()

	err := PurgeCache()
	require.NoError(t, err, "failed to remove the cache (if it already existed)")

	goBin, err := exec.LookPath("go")
	require.NoError(t, err, "failed to find the go executable")

	envvars, err := createCache(context.Background(), goBin)
	require.NoError(t, err, "failed to create the cache")

	dir := GetCacheDir()
	info, err := os.Stat(dir)
	require.NoError(t, err, "failed to find the cache")
	require.True(t, info.IsDir(), "cache isn't a directory")

	testCompareDirs(t, envvars, goBin, "GOMODCACHE")
	testCompareDirs(t, envvars, goBin, "GOCACHE")

	err = PurgeCache()
	require.NoError(t, err, "failed to remove the cache")


	_, err = os.Stat(dir)
	require.Error(t, err, "cache still exists after purge")
	require.True(t, os.IsNotExist(err), "failed to check if the cache still exists after purge")
}

func testCompareDirs(t *testing.T, envvars map[string]string, executable, envvar string) {
	cmd := exec.Command(executable, "env", envvar)
	stdout, err := cmd.Output()
	require.NoError(t, err, "%s: couldn't read original path", envvar)
	src := strings.TrimSpace(string(stdout))

	dst := envvars[envvar]
	require.NotEmpty(t, dst, "%s: envvar is empty", envvar)

	// Ideally, we should compare the contents here.
	// One way to do it would be to hash the archive of the directories and compare that...
	// However, zip compresses the data (and thus is slow) and tar stores the file's timestamp.
	// So, calling 'diff' directly is both faster and more correct.
	cmd = exec.Command("diff", "-r", src, dst)
	err = cmd.Run()
	require.NoError(t, err, "%s: directories do not match", envvar)
}
