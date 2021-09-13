// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package internal

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternal(t *testing.T) {
	t.Parallel()

	gitDir := filepath.Join(os.TempDir(), "git-verify-chain")

	_, err := os.Stat(gitDir)
	if os.IsNotExist(err) {
		err = nil

		require.NoError(t, exec.Command("git", "clone", "https://github.com/talos-systems/talos.git", gitDir).Run())
		require.NoError(t, exec.Command("git", "-C", gitDir, "checkout", "v0.12.0").Run())
	}

	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("listCommits", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"ff9681a74e4927e147753aed6a283b2d9ad22c29",
		}
		commits, err := listCommits(ctx, gitDir, "v0.12.0-beta.2", "")
		require.NoError(t, err)
		assert.Equal(t, expected, commits)

		expected = []string{
			"75ce68d909309e019c4a37ce5e7bcb6eda8aa1e4",
			"87c258093a17dd5725b84377a3f42270c977df63",
			"eba00723d3df0d8be30c9d3931421117ea74deea",
			"3a38f0deda18906b785b1ad76872a78b689eb171",
			"2e220cb65876e5ee0153afe67e1929beacd27152",
			"b63a2ea0e230046fca327c420250b25d722febf5",
			"cd0532848f80625c67ff98947d206c7faf4d27c7",
			"e22301e762e6f112d64d9fb3a91cb9b4db877b4e",
		}
		commits, err = listCommits(ctx, gitDir, "v0.12.0-beta.1", "v0.12.0-beta.2")
		require.NoError(t, err)
		assert.Equal(t, expected, commits)

		expected = []string{}
		commits, err = listCommits(ctx, gitDir, "v0.12.0-beta.2", "v0.12.0-beta.2")
		require.NoError(t, err)
		assert.Equal(t, expected, commits)
	})

	t.Run("setupKeyring", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		files, err := filepath.Glob(filepath.Join(wd, "testdata", "*.gpg"))
		require.NoError(t, err)

		err = setupKeyring(ctx, dir, files)
		require.NoError(t, err)
	})

	t.Run("VerifyWithFiles", func(t *testing.T) {
		t.Parallel()

		// smira is the committer of all commits
		err := VerifyWithFiles(ctx, gitDir, "v0.12.0-beta.0", filepath.Join(wd, "testdata"))
		require.NoError(t, err)

		err = VerifyWithFiles(ctx, gitDir, "v0.12.0-alpha.0", filepath.Join(wd, "testdata"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), `failed to verify commit "1ed5e545385e160fe3b61e6dbbcaa8a701437b62"`)
	})
}
