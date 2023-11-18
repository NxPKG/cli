//go:build rust
// +build rust

package hashing

import (
	"github.com/nxpkg/cli/internal/ffi"
	"github.com/nxpkg/cli/internal/nxpkgpath"
)

func GetPackageFileHashes(rootPath nxpkgpath.AbsoluteSystemPath, packagePath nxpkgpath.AnchoredSystemPath, inputs []string) (map[nxpkgpath.AnchoredUnixPath]string, error) {
	rawHashes, err := ffi.GetPackageFileHashes(rootPath.ToString(), packagePath.ToString(), inputs)
	if err != nil {
		return nil, err
	}

	hashes := make(map[nxpkgpath.AnchoredUnixPath]string, len(rawHashes))
	for rawPath, hash := range rawHashes {
		hashes[nxpkgpath.AnchoredUnixPathFromUpstream(rawPath)] = hash
	}
	return hashes, nil
}

func GetHashesForFiles(rootPath nxpkgpath.AbsoluteSystemPath, files []nxpkgpath.AnchoredSystemPath) (map[nxpkgpath.AnchoredUnixPath]string, error) {
	rawFiles := make([]string, len(files))
	for i, file := range files {
		rawFiles[i] = file.ToString()
	}
	rawHashes, err := ffi.GetHashesForFiles(rootPath.ToString(), rawFiles, false)
	if err != nil {
		return nil, err
	}

	hashes := make(map[nxpkgpath.AnchoredUnixPath]string, len(rawHashes))
	for rawPath, hash := range rawHashes {
		hashes[nxpkgpath.AnchoredUnixPathFromUpstream(rawPath)] = hash
	}
	return hashes, nil
}

func GetHashesForExistingFiles(rootPath nxpkgpath.AbsoluteSystemPath, files []nxpkgpath.AnchoredSystemPath) (map[nxpkgpath.AnchoredUnixPath]string, error) {
	rawFiles := make([]string, len(files))
	for i, file := range files {
		rawFiles[i] = file.ToString()
	}
	rawHashes, err := ffi.GetHashesForFiles(rootPath.ToString(), rawFiles, true)
	if err != nil {
		return nil, err
	}

	hashes := make(map[nxpkgpath.AnchoredUnixPath]string, len(rawHashes))
	for rawPath, hash := range rawHashes {
		hashes[nxpkgpath.AnchoredUnixPathFromUpstream(rawPath)] = hash
	}
	return hashes, nil
}
