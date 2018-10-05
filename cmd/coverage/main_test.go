package main

import "testing"

type coverProfileTestCase struct {
	BaseDir  string
	FileName string
	Expected string
}

func TestJoinCoverPath(t *testing.T) {
	testCases := []coverProfileTestCase{
		{
			"/",
			"foo/bar/baz.go",
			"/foo/bar/baz.go",
		},
		{
			"/users/foo/bar",
			"foo/bar/baz.go",
			"/users/foo/bar/baz.go",
		},
		{
			"/users/bailey/workspace/go/src/github.com/blend/go-sdk/",
			"github.com/blend/go-sdk/assert/assert.go",
			"/users/bailey/workspace/go/src/github.com/blend/go-sdk/assert/assert.go",
		},
		{
			"/go/src/git.blendlabs.com/blend/fees",
			"git.blendlabs.com/blend/fees/pkg/fees/fees.go",
			"/go/src/git.blendlabs.com/blend/fees/pkg/fees/fees.go",
		},
	}

	var actual string
	for _, testCase := range testCases {
		actual = joinCoverPath(testCase.BaseDir, testCase.FileName)
		if actual != testCase.Expected {
			t.Errorf("%s does not match %s", actual, testCase.Expected)
		}
	}
}
