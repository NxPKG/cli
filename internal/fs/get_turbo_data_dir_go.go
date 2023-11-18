//go:build go || !rust
// +build go !rust

package fs

import (
	"github.com/adrg/xdg"
	"github.com/nxpkg/cli/internal/nxpkgpath"
)

// GetNxpkgDataDir returns a directory outside of the repo
// where nxpkg can store data files related to nxpkg.
func GetNxpkgDataDir() nxpkgpath.AbsoluteSystemPath {
	dataHome := AbsoluteSystemPathFromUpstream(xdg.DataHome)
	return dataHome.UntypedJoin("nxpkgrepo")
}
