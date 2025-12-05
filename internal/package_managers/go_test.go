package package_managers

import (
	"fmt"
	"os"
	"testing"

	"github.com/New-Horizons-Team/supply-chain-firewall/internal/temp_directory"
	"github.com/stretchr/testify/require"
)

func TestExtractGoTargets(t *testing.T) {
	type testCase struct {
		command            []string
		wantLocalPackages  []string
		wantRemotePackages []string
		wantGetFlags       []string
		createPath         string
	}

	testCases := []testCase{{
		command:           []string{"go", "build"},
		wantLocalPackages: []string{"."},
	}, {
		command:           []string{"go", "build", "."},
		wantLocalPackages: []string{"."},
	}, {
		command:           []string{"go", "build", "-o", "/tmp/app.bin", "."},
		wantLocalPackages: []string{"."},
	}, {
		command:           []string{"go", "build", "-o=/tmp/app.bin", "."},
		wantLocalPackages: []string{"."},
	}, {
		command:           []string{"go", "build", "-o", "/tmp/app.bin", "./cmd/run"},
		wantLocalPackages: []string{"./cmd/run"},
		createPath:        "./cmd/run",
	}, {
		command:            []string{"go", "install", "github.com/go-delve/delve/cmd/dlv@latest"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv@latest"},
	}, {
		command:            []string{"go", "install", "-buildmode", "exe", "github.com/go-delve/delve/cmd/dlv@latest"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv@latest"},
	}, {
		command:            []string{"go", "install", "-buildmode=exe", "github.com/go-delve/delve/cmd/dlv@latest"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv@latest"},
	}, {
		command:            []string{"go", "get", "github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
	}, {
		command:            []string{"go", "get", "-t", "github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantGetFlags:       []string{"-t"},
	}, {
		command:            []string{"go", "get", "-u", "github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantGetFlags:       []string{"-u"},
	}, {
		command:            []string{"go", "get", "-u=patch", "github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv", "github.com/stretchr/testify/require"},
		wantGetFlags:       []string{"-u=patch"},
	}, {
		command:           []string{"go", "mod", "tidy"},
		wantLocalPackages: []string{"."},
	}, {
		command:           []string{"go", "mod", "download"},
		wantLocalPackages: []string{"."},
	}, {
		// This command is invalid, but checks whether the flag parser is working as expected.
		command:            []string{"go", "mod", "download", "-buildmode", "exe", "github.com/go-delve/delve/cmd/dlv@latest"},
		wantLocalPackages:  []string{"exe"},
		wantRemotePackages: []string{"github.com/go-delve/delve/cmd/dlv@latest"},
	}}

	cwd, err := os.Getwd()
	require.NoError(t, err, "failed to retrieved the current directory")

	for _, tc := range testCases {
		name := fmt.Sprintf("%s", tc.command)
		name = name[1 : len(name)-1]

		t.Run(name, func(t *testing.T) {
			if tc.createPath != "" {
				tmp, err := temp_directory.NewTempDir("")
				require.NoError(t, err, "failed to create the temporary directory")
				defer func() {
					_ = tmp.Close()
				}()

				err = os.Chdir(tmp.GetPath())
				require.NoError(t, err, "failed to move into the temporary directory")
				defer func() {
					_ = os.Chdir(cwd)
				}()

				err = os.MkdirAll(tc.createPath, 0755)
				require.NoError(t, err, "failed to create the dummy local path")
			}

			gotLocalPackages, gotRemotePackages, gotGetFlags := extractGoTargets(tc.command)
			require.Equal(t, tc.wantLocalPackages, gotLocalPackages, "local packages doesn't match")
			require.Equal(t, tc.wantRemotePackages, gotRemotePackages, "remote packages doesn't match")
			require.Equal(t, tc.wantGetFlags, gotGetFlags, "get flags doesn't match")
		})
	}
}
