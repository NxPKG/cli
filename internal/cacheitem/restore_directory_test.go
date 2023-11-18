package cacheitem

import (
	"reflect"
	"testing"

	"github.com/nxpkg/cli/internal/nxpkgpath"
)

func Test_cachedDirTree_getStartingPoint(t *testing.T) {
	testDir := nxpkgpath.AbsoluteSystemPath("")
	tests := []struct {
		name string

		// STATE
		cachedDirTree cachedDirTree

		// INPUT
		path nxpkgpath.AnchoredSystemPath

		// OUTPUT
		calculatedAnchor nxpkgpath.AbsoluteSystemPath
		pathSegments     []nxpkgpath.RelativeSystemPath
	}{
		{
			name: "hello world",
			cachedDirTree: cachedDirTree{
				anchorAtDepth: []nxpkgpath.AbsoluteSystemPath{testDir},
				prefix:        []nxpkgpath.RelativeSystemPath{},
			},
			path:             nxpkgpath.AnchoredUnixPath("hello/world").ToSystemPath(),
			calculatedAnchor: testDir,
			pathSegments:     []nxpkgpath.RelativeSystemPath{"hello", "world"},
		},
		{
			name: "has a cache",
			cachedDirTree: cachedDirTree{
				anchorAtDepth: []nxpkgpath.AbsoluteSystemPath{
					testDir,
					testDir.UntypedJoin("hello"),
				},
				prefix: []nxpkgpath.RelativeSystemPath{"hello"},
			},
			path:             nxpkgpath.AnchoredUnixPath("hello/world").ToSystemPath(),
			calculatedAnchor: testDir.UntypedJoin("hello"),
			pathSegments:     []nxpkgpath.RelativeSystemPath{"world"},
		},
		{
			name: "ask for yourself",
			cachedDirTree: cachedDirTree{
				anchorAtDepth: []nxpkgpath.AbsoluteSystemPath{
					testDir,
					testDir.UntypedJoin("hello"),
					testDir.UntypedJoin("hello", "world"),
				},
				prefix: []nxpkgpath.RelativeSystemPath{"hello", "world"},
			},
			path:             nxpkgpath.AnchoredUnixPath("hello/world").ToSystemPath(),
			calculatedAnchor: testDir.UntypedJoin("hello", "world"),
			pathSegments:     []nxpkgpath.RelativeSystemPath{},
		},
		{
			name: "three layer cake",
			cachedDirTree: cachedDirTree{
				anchorAtDepth: []nxpkgpath.AbsoluteSystemPath{
					testDir,
					testDir.UntypedJoin("hello"),
					testDir.UntypedJoin("hello", "world"),
				},
				prefix: []nxpkgpath.RelativeSystemPath{"hello", "world"},
			},
			path:             nxpkgpath.AnchoredUnixPath("hello/world/again").ToSystemPath(),
			calculatedAnchor: testDir.UntypedJoin("hello", "world"),
			pathSegments:     []nxpkgpath.RelativeSystemPath{"again"},
		},
		{
			name: "outside of cache hierarchy",
			cachedDirTree: cachedDirTree{
				anchorAtDepth: []nxpkgpath.AbsoluteSystemPath{
					testDir,
					testDir.UntypedJoin("hello"),
					testDir.UntypedJoin("hello", "world"),
				},
				prefix: []nxpkgpath.RelativeSystemPath{"hello", "world"},
			},
			path:             nxpkgpath.AnchoredUnixPath("somewhere/else").ToSystemPath(),
			calculatedAnchor: testDir,
			pathSegments:     []nxpkgpath.RelativeSystemPath{"somewhere", "else"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := tt.cachedDirTree
			calculatedAnchor, pathSegments := cr.getStartingPoint(tt.path)
			if !reflect.DeepEqual(calculatedAnchor, tt.calculatedAnchor) {
				t.Errorf("cachedDirTree.getStartingPoint() calculatedAnchor = %v, want %v", calculatedAnchor, tt.calculatedAnchor)
			}
			if !reflect.DeepEqual(pathSegments, tt.pathSegments) {
				t.Errorf("cachedDirTree.getStartingPoint() pathSegments = %v, want %v", pathSegments, tt.pathSegments)
			}
		})
	}
}
