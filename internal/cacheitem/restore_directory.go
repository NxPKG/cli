package cacheitem

import (
	"archive/tar"
	"os"
	"path/filepath"
	"strings"

	"github.com/nxpkg/cli/internal/nxpkgpath"
)

// restoreDirectory restores a directory.
func restoreDirectory(dirCache *cachedDirTree, anchor nxpkgpath.AbsoluteSystemPath, header *tar.Header) (nxpkgpath.AnchoredSystemPath, error) {
	processedName, err := canonicalizeName(header.Name)
	if err != nil {
		return "", err
	}

	// We need to traverse `processedName` from base to root split at
	// `os.Separator` to make sure we don't end up following a symlink
	// outside of the restore path.

	// Create the directory.
	if err := safeMkdirAll(dirCache, anchor, processedName, header.Mode); err != nil {
		return "", err
	}

	return processedName, nil
}

type cachedDirTree struct {
	anchorAtDepth []nxpkgpath.AbsoluteSystemPath
	prefix        []nxpkgpath.RelativeSystemPath
}

func (cr *cachedDirTree) getStartingPoint(path nxpkgpath.AnchoredSystemPath) (nxpkgpath.AbsoluteSystemPath, []nxpkgpath.RelativeSystemPath) {
	pathSegmentStrings := strings.Split(path.ToString(), string(os.PathSeparator))
	pathSegments := make([]nxpkgpath.RelativeSystemPath, len(pathSegmentStrings))
	for index, pathSegmentString := range pathSegmentStrings {
		pathSegments[index] = nxpkgpath.RelativeSystemPathFromUpstream(pathSegmentString)
	}

	i := 0
	for i = 0; i < len(cr.prefix) && i < len(pathSegments); i++ {
		if pathSegments[i] != cr.prefix[i] {
			break
		}
	}

	// 0: root anchor, can't remove it.
	cr.anchorAtDepth = cr.anchorAtDepth[:i+1]

	// 0: first prefix.
	cr.prefix = cr.prefix[:i]

	return cr.anchorAtDepth[i], pathSegments[i:]
}

func (cr *cachedDirTree) Update(anchor nxpkgpath.AbsoluteSystemPath, newSegment nxpkgpath.RelativeSystemPath) {
	cr.anchorAtDepth = append(cr.anchorAtDepth, anchor)
	cr.prefix = append(cr.prefix, newSegment)
}

// safeMkdirAll creates all directories, assuming that the leaf node is a directory.
// FIXME: Recheck the symlink cache before creating a directory.
func safeMkdirAll(dirCache *cachedDirTree, anchor nxpkgpath.AbsoluteSystemPath, processedName nxpkgpath.AnchoredSystemPath, mode int64) error {
	// Iterate through path segments by os.Separator, appending them onto the anchor.
	// Check to see if that path segment is a symlink with a target outside of anchor.

	// Pull the iteration starting point from the directory cache.
	calculatedAnchor, pathSegments := dirCache.getStartingPoint(processedName)
	for _, segment := range pathSegments {
		calculatedAnchor, checkPathErr := checkPath(anchor, calculatedAnchor, segment)
		// We hit an existing directory or absolute path that was invalid.
		if checkPathErr != nil {
			return checkPathErr
		}

		// Otherwise we continue and check the next segment.
		dirCache.Update(calculatedAnchor, segment)
	}

	// If we have made it here we know that it is safe to call os.MkdirAll
	// on the Join of anchor and processedName.
	//
	// This could _still_ error, but we don't care.
	return processedName.RestoreAnchor(anchor).MkdirAll(os.FileMode(mode))
}

// checkPath ensures that the resolved path (if restoring symlinks).
// It makes sure to never traverse outside of the anchor.
func checkPath(originalAnchor nxpkgpath.AbsoluteSystemPath, accumulatedAnchor nxpkgpath.AbsoluteSystemPath, segment nxpkgpath.RelativeSystemPath) (nxpkgpath.AbsoluteSystemPath, error) {
	// Check if the segment itself is sneakily an absolute path...
	// (looking at you, Windows. CON, AUX...)
	if filepath.IsAbs(segment.ToString()) {
		return "", errTraversal
	}

	// Find out if this portion of the path is a symlink.
	combinedPath := accumulatedAnchor.Join(segment)
	fileInfo, err := combinedPath.Lstat()

	// Getting an error here means we failed to stat the path.
	// Assume that means we're safe and continue.
	if err != nil {
		return combinedPath, nil
	}

	// Find out if we have a symlink.
	isSymlink := fileInfo.Mode()&os.ModeSymlink != 0

	// If we don't have a symlink it's safe.
	if !isSymlink {
		return combinedPath, nil
	}

	// Check to see if the symlink targets outside of the originalAnchor.
	// We don't do eval symlinks because we could find ourself in a totally
	// different place.

	// 1. Get the target.
	linkTarget, readLinkErr := combinedPath.Readlink()
	if readLinkErr != nil {
		return "", readLinkErr
	}

	// 2. See if the target is absolute.
	if filepath.IsAbs(linkTarget) {
		absoluteLinkTarget := nxpkgpath.AbsoluteSystemPathFromUpstream(linkTarget)
		if originalAnchor.HasPrefix(absoluteLinkTarget) {
			return absoluteLinkTarget, nil
		}
		return "", errTraversal
	}

	// 3. Target is relative (or absolute Windows on a Unix device)
	relativeLinkTarget := nxpkgpath.RelativeSystemPathFromUpstream(linkTarget)
	computedTarget := accumulatedAnchor.UntypedJoin(linkTarget)
	if computedTarget.HasPrefix(originalAnchor) {
		// Need to recurse and make sure the target doesn't link out.
		return checkPath(originalAnchor, accumulatedAnchor, relativeLinkTarget)
	}
	return "", errTraversal
}
