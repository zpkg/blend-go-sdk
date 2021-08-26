/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package bindata

import "time"

// File is both the file metadata and the contents.
type File struct {
	Name		string
	Modtime		time.Time
	Contents	*FileCompressor
}
