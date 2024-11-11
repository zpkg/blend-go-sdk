/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"fmt"
	"os"

	"github.com/zpkg/blend-go-sdk/certutil"
	"github.com/zpkg/blend-go-sdk/uuid"
)

func main() {
	ca, _ := certutil.CreateCertificateAuthority(certutil.OptSubjectCommonName("go-sdk certificate authority"))

	ca.WriteCertPem(os.Stdout)
	fmt.Println()

	ca.WriteKeyPem(os.Stdout)
	fmt.Println()

	certBundle, _ := certutil.CreateServer(uuid.V4().String(), ca)

	certBundle.WriteCertPem(os.Stdout)
	fmt.Println()
	certBundle.WriteKeyPem(os.Stdout)
	fmt.Println()
}
