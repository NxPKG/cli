//go:build windows
// +build windows

package cacheitem

import (
	"testing"

	"github.com/nxpkg/cli/internal/nxpkgpath"
)

func createFifo(t *testing.T, anchor nxpkgpath.AbsoluteSystemPath, fileDefinition createFileDefinition) error {
	return errUnsupportedFileType
}
