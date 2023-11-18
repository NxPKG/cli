package cmdutil

import (
	"os"
	"testing"
	"time"

	"github.com/nxpkg/cli/internal/nxpkgstate"
	"gotest.tools/v3/assert"
)

func TestRemoteCacheTimeoutFlag(t *testing.T) {
	args := nxpkgstate.ParsedArgsFromRust{
		CWD: "",
	}

	executionState := nxpkgstate.ExecutionState{
		APIClientConfig: nxpkgstate.APIClientConfig{
			Timeout: 599,
		},
		CLIArgs: args,
	}

	h := NewHelper("test-version", &args)

	base, err := h.GetCmdBase(&executionState)
	if err != nil {
		t.Fatalf("failed to get command base %v", err)
	}

	assert.Equal(t, base.APIClient.HTTPClient.HTTPClient.Timeout, time.Duration(599)*time.Second)
}

func TestRemoteCacheTimeoutPrimacy(t *testing.T) {
	key := "NXPKG_REMOTE_CACHE_TIMEOUT"
	value := "2"

	t.Run(key, func(t *testing.T) {
		t.Cleanup(func() {
			_ = os.Unsetenv(key)
		})
		args := nxpkgstate.ParsedArgsFromRust{
			CWD: "",
		}
		executionState := nxpkgstate.ExecutionState{
			APIClientConfig: nxpkgstate.APIClientConfig{
				Timeout: 1,
			},
			CLIArgs: args,
		}
		h := NewHelper("test-version", &args)

		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("setenv %v", err)
		}

		base, err := h.GetCmdBase(&executionState)
		if err != nil {
			t.Fatalf("failed to get command base %v", err)
		}
		assert.Equal(t, base.APIClient.HTTPClient.HTTPClient.Timeout, time.Duration(1)*time.Second)
	})
}
