/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bindata

import "time"

// File is both the file metadata and the contents.
type File struct {
	Name     string
	Modtime  time.Time
	Contents *FileCompressor
}
