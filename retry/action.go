/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package retry

import "context"

// Action is a function you can retry.
type Action func(ctx context.Context) (interface{}, error)
