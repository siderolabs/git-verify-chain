// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package internal contains internal program logic.
package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// VerifyWithFiles checks that all git commits in dir in from..HEAD range are signed by *.gpg key files from pubKeysDir.
func VerifyWithFiles(ctx context.Context, dir, from, pubKeysDir string) error {
	gpgHomeDir, err := os.MkdirTemp("", "git-verify-chain-keyring-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(gpgHomeDir) //nolint:errcheck

	files, err := filepath.Glob(filepath.Join(pubKeysDir, "*.gpg"))
	if err != nil {
		return err
	}

	if err = setupKeyring(ctx, gpgHomeDir, files); err != nil {
		return err
	}

	return VerifyWithKeyring(ctx, dir, from, gpgHomeDir)
}

// Verify checks that all git commits in dir in from..HEAD range are signed by keys from the GnuPG keyring in gpgHomeDir.
func VerifyWithKeyring(ctx context.Context, dir, from, gpgHomeDir string) error {
	commits, err := listCommits(ctx, dir, from, "")
	if err != nil {
		return err
	}

	return verifyCommits(ctx, dir, commits, gpgHomeDir)
}

func listCommits(ctx context.Context, dir, from, to string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-list", from+".."+to)
	cmd.Dir = dir
	stdout, err := run(cmd)

	return stdout, err
}

func verifyCommits(ctx context.Context, dir string, commits []string, gpgHomeDir string) error {
	// we can pass --raw and parse output (https://github.com/gpg/gnupg/blob/master/doc/DETAILS#format-of-the-status-fd-output)
	// if we need more details like trust levels, etc.
	args := make([]string, 0, len(commits))
	args = append(args, "verify-commit")
	args = append(args, commits...)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = []string{
		"GNUPGHOME=" + gpgHomeDir,
	}

	_, err := run(cmd)
	if err == nil {
		return nil
	}

	// find the first bad commit
	for _, commit := range commits {
		args = []string{"verify-commit", "--verbose", commit}

		cmd = exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = []string{
			"GNUPGHOME=" + gpgHomeDir,
		}

		_, err = run(cmd)
		if err != nil {
			return fmt.Errorf("failed to verify commit %q:\n%w", commit, err)
		}
	}

	panic("not reached")
}
