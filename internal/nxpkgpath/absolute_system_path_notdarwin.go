//go:build !darwin
// +build !darwin

package nxpkgpath

import (
	"os"
)

// Lchmod changes the mode of a file not following symlinks.
func (p AbsoluteSystemPath) Lchmod(mode os.FileMode) error {
	return nil
}
