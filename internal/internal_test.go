// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package internal

import (
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

		require.NoError(t, exec.CommandContext(t.Context(), "git", "clone", "https://github.com/siderolabs/talos.git", gitDir).Run())
		require.NoError(t, exec.CommandContext(t.Context(), "git", "-C", gitDir, "checkout", "v1.12.0").Run())
	}

	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	ctx := t.Context()

	t.Run("listCommits", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"ac91ade2c7e435e63ed2546244d428a81abd22ad",
			"82553b2a1a713836f496b822e86e5e6788c5ebd1",
			"33f6e22ecb3b393d1488730c67d6f973a46b0b39",
			"d5be50ac55cac1c1c1deff4971fd991f364696a1",
			"70d3ab9ac090095c2fc8cbbfaa9c5c472d76c794",
			"101814d889924afe7c049106c638a32ae107a139",
			"ce286825a7f969f847ea7ad17bd2a31fa85d301c",
			"96f724adccbc6fac844f9a341e36eede331b3947",
			"e195427c17a004b5bcaa6f1870ce6c855ae61f1d",
			"e025355b759bb110925631f5f84230e99b9069df",
			"21a914a1d1ca48d6bb4d47ddc8be0d0fdf74800d",
			"ca645777dae5ad07501501dafc4717e7383045b0",
			"6dd0558a314af9a0dfda77b4f58a7115ef86b6fc",
			"c931847ccaadf84f84e5f2befadaffb55740b592",
			"a2a77004deac3efe6ac14f906a8bd0a3b0f926ca",
			"47198780bfc084347b9ae675aaeb27a1c1d58d38",
			"03a424bdf1b8a270dd694fc2738d81a3261d80cf",
			"688fb789beb979544e16447e512419629ea61b21",
			"66e67fd1394946b3425543a1aac52d4a8338e375",
			"d8403498c92e0f9c37b04ad6786b2c84df5e7c95",
			"5ced4258c18f5649590a50c2927ab8e16db298ec",
			"fabf3f0e73918b650b33ef0f009cacb9a15ecbc0",
			"93cec4b9dfdef0566152ef80c28439a7dbb0c320",
			"964098d9696a804de5d27284cd79dccffa7c81b9",
			"bce04084d6f5a9c703c7d63d1558d7d43c54dfbf",
			"d1abc0f8473c1a562e37a712624f803ce0f60fec",
		}
		commits, err := listCommits(ctx, gitDir, "v1.12.0-rc.0", "")
		require.NoError(t, err)
		assert.Equal(t, expected, commits)

		expected = []string{
			"0613076873bbd2d763da30ae2e9e1903486f7cb8",
			"bc4de5b7926a9a2e7a7af9da4763effb5c33693e",
			"4a15763a962cad0c020e01f66948ba1f326c9201",
			"29733654902be5cb72b71a9a64ea0ed3c0a0f011",
			"0ac58929db6960ef91c1bcfbc891264e18e1e930",
			"184a45c405530c73c31d5b6c642cda4ddd1772ca",
			"8eac9f37d9dddc507c988cfb187b939a5624f563",
			"e79a94d57781d6ede61e6205f6f5d0f0708a8ddb",
			"7a1bb4c26a99c7f4e37196b40aced6334eeda731",
			"5c6ee6aceeb87785c08a05f2ddc6b7cbcad0bc9a",
			"2e6fe4684b98ca4432284b7b51dfcd1a8b91a03c",
			"473bc17c199165dd0f925981753dec431cc5613b",
			"6dc8e82b31d095a357b9f6d99420bb860e51261c",
			"a7dbbbd4d87feeace427e4c63f67880c72f7cd22",
			"3ca342c0979ffcfe7bee95a4e56c98ddece8abb5",
			"364ebb6baf3c77a1e2dd28d83b6af7cfe821e1e8",
			"aa286d3f6eb28a813c982a9cc1230c138e56b33a",
			"f4891eebb192d2895f27f85502fd223290217d90",
		}
		commits, err = listCommits(ctx, gitDir, "v1.12.0-beta.1", "v1.12.0-rc.0")
		require.NoError(t, err)
		assert.Equal(t, expected, commits)

		expected = []string{}
		commits, err = listCommits(ctx, gitDir, "v1.12.0-beta.1", "v1.12.0-beta.1")
		require.NoError(t, err)
		assert.Equal(t, expected, commits)
	})

	t.Run("setupKeyring", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		files, err := filepath.Glob(filepath.Join(wd, "testdata", "*.pgp"))
		require.NoError(t, err)

		err = setupKeyring(ctx, dir, files)
		require.NoError(t, err)
	})

	t.Run("VerifyWithFiles", func(t *testing.T) {
		t.Parallel()

		err := VerifyWithFiles(ctx, gitDir, "v1.12.0-rc.0", filepath.Join(wd, "testdata"))
		require.NoError(t, err)

		err = VerifyWithFiles(ctx, gitDir, "v1.12.0-alpha.0", filepath.Join(wd, "testdata"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), `failed to verify commit "0613076873bbd2d763da30ae2e9e1903486f7cb8"`)
	})
}
