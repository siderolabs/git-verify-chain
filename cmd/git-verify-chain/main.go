// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"log"
	"os/user"
	"path/filepath"

	"github.com/talos-systems/git-verify-chain/internal"
)

func main() {
	log.SetFlags(0)

	var gpgHome string
	if u, _ := user.Current(); u != nil { //nolint:errcheck
		gpgHome = filepath.Join(u.HomeDir, ".gnupg")
	}

	fromF := flag.String("from", "", "commit or tag name to verify from")
	keysDirF := flag.String("keys-dir", "", "directory with *.gpg files")
	gpgHomeF := flag.String("gpg-home", gpgHome, "GnuPG home directory")
	flag.Parse()

	if *fromF == "" {
		log.Fatal("-from flag is required.")
	}

	ctx := context.Background()

	var err error

	if *keysDirF != "" {
		log.Printf("Verifying using public keys from %s ...", *keysDirF)
		err = internal.VerifyWithFiles(ctx, ".", *fromF, *keysDirF)
	} else {
		log.Printf("Verifying using GnuPG keyring from %s ...", *gpgHomeF)
		err = internal.VerifyWithKeyring(ctx, ".", *fromF, *gpgHomeF)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
}
