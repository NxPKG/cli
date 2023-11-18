package nodes

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLogFilename(t *testing.T) {
	testCases := []struct{ input, want string }{
		{
			"build",
			"nxpkg.github.io.log",
		},
		{
			"build:prod",
			"nxpkg.github.io$colon$prod.log",
		},
		{
			"build:prod:extra",
			"nxpkg.github.io$colon$prod$colon$extra.log",
		},
	}

	for _, testCase := range testCases {
		got := logFilename(testCase.input)
		assert.Equal(t, got, testCase.want)
	}
}
