/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"testing"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/testutil"
)

func TestMain(m *testing.M) {
	testutil.MarkUpdateGoldenFlag()

	testutil.New(
		m,
		testutil.OptLog(logger.All()),
	).Run()
}

const buildTags1 = `//go:build tag1
// +build tag1

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}
`
const buildTags2 = `// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3
// +build tag6

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

// +bulid tag9000
`

const buildTags3 = `//go:build tag1 & tag2

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

//go:build tag3
`

const goldenBuildTags1 = `//go:build tag1
// +build tag1

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}
`

const goldenBuildTags2 = `// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3
// +build tag6

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

// +bulid tag9000
`

const goldenBuildTags3 = `//go:build tag1 & tag2

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

//go:build tag3
`
