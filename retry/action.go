/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package retry

import "context"

// Action is a function you can retry.
type Action func(ctx context.Context) (interface{}, error)
