package graph

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_CommandsInvokingNxpkg(t *testing.T) {
	type testCase struct {
		command string
		match   bool
	}
	testCases := []testCase{
		{
			"nxpkg run foo",
			true,
		},
		{
			"rm -rf ~/Library/Caches/pnpm && nxpkg run foo && rm -rf ~/.npm",
			true,
		},
		{
			"FLAG=true nxpkg run foo",
			true,
		},
		{
			"npx nxpkg run foo",
			true,
		},
		{
			"echo starting; nxpkg foo; echo done",
			true,
		},
		// We don't catch this as if people are going to try to invoke the nxpkg
		// binary directly, they'll always be able to work around us.
		{
			"./node_modules/.bin/nxpkg foo",
			false,
		},
		{
			"rm -rf ~/Library/Caches/pnpm && rm -rf ~/Library/Caches/nxpkg && rm -rf ~/.npm && rm -rf ~/.pnpm-store && rm -rf ~/.nxpkg",
			false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, commandLooksLikeNxpkg(tc.command), tc.match, tc.command)
	}
}
