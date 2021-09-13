// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// runError contains all details for better error output.
type runError struct {
	args   []string
	stdout []string
	stderr []string
	err    error
}

func (e *runError) Error() string {
	if e.err == nil {
		return ""
	}

	msg := e.err.Error()

	if e.args != nil {
		msg = strings.Join(e.args, " ") + ": " + msg
	}

	if e.stdout != nil {
		msg += "\n" + strings.Join(e.stdout, "\n")
	}

	if e.stderr != nil {
		msg += "\n" + strings.Join(e.stderr, "\n")
	}

	return msg
}

func run(cmd *exec.Cmd) (stdout, stderr []string, err error) {
	var outB, errB bytes.Buffer

	if cmd.Dir == "" {
		err = &runError{
			err: fmt.Errorf("exec.Cmd.Dir is empty"),
		}
	}

	if cmd.Env == nil {
		cmd.Env = []string{} // do not inherit from the parent process
	}

	cmd.Stdout = &outB
	cmd.Stderr = &errB

	err = cmd.Run()

	stdout, stderr = []string{}, []string{}

	outS := strings.TrimSpace(outB.String())
	if outS != "" {
		stdout = strings.Split(outS, "\n")
	}

	errS := strings.TrimSpace(errB.String())
	if errS != "" {
		stderr = strings.Split(errS, "\n")
	}

	if err != nil {
		err = &runError{
			args:   cmd.Args,
			stdout: stdout,
			stderr: stderr,
			err:    err,
		}
	}

	return
}
