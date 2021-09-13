// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package internal

import (
	"context"
	"os/exec"
)

func setupKeyring(ctx context.Context, dir string, pubKeyFiles []string) error {
	for _, file := range pubKeyFiles {
		args := []string{
			"--no-auto-key-locate",
			"--no-auto-key-retrieve",
			"--no-autostart",
			"--import",
			file,
		}
		cmd := exec.CommandContext(ctx, "gpg", args...)
		cmd.Dir = dir
		cmd.Env = []string{
			"GNUPGHOME=" + dir,
		}
		_, _, err := run(cmd)

		if err != nil {
			return err
		}
	}

	return nil
}
