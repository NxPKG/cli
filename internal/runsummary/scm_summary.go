package runsummary

import (
	"github.com/nxpkg/cli/internal/ci"
	"github.com/nxpkg/cli/internal/env"
	"github.com/nxpkg/cli/internal/scm"
	"github.com/nxpkg/cli/internal/nxpkgpath"
)

type scmState struct {
	Type   string `json:"type"`
	Sha    string `json:"sha"`
	Branch string `json:"branch"`
}

// getSCMState returns the sha and branch when in a git repo
// Otherwise it should return empty strings right now.
// We my add handling of other scms and non-git tracking in the future.
func getSCMState(envVars env.EnvironmentVariableMap, dir nxpkgpath.AbsoluteSystemPath) *scmState {

	state := &scmState{Type: "git"}

	// If we're in CI, try to get the values we need from environment variables
	if ci.IsCi() {
		vendor := ci.Info()
		state.Sha = envVars[vendor.ShaEnvVar]
		state.Branch = envVars[vendor.BranchEnvVar]
	}

	// Otherwise fallback to using `git`
	if state.Branch == "" {
		state.Branch = scm.GetCurrentBranch(dir)
	}

	if state.Sha == "" {
		state.Sha = scm.GetCurrentSha(dir)
	}

	return state
}
