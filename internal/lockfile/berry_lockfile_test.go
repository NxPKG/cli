package lockfile

import (
	"testing"

	"github.com/nxpkg/cli/internal/nxpkgpath"
	"gotest.tools/v3/assert"
)

func Test_BerryTransitiveClosure(t *testing.T) {
	contents := getRustFixture(t, "berry.lock")
	lf, err := DecodeBerryLockfile(contents, map[string]string{"lodash@^4.17.21": "patch:lodash@npm%3A4.17.21#./.yarn/patches/lodash-npm-4.17.21-6382451519.patch"})
	assert.NilError(t, err)
	closures, err := AllTransitiveClosures(map[nxpkgpath.AnchoredUnixPath]map[string]string{
		nxpkgpath.AnchoredUnixPath(""):         {},
		nxpkgpath.AnchoredUnixPath("apps/web"): {},
		nxpkgpath.AnchoredUnixPath("apps/docs"): {
			"lodash": "^4.17.21",
		},
	}, lf)
	assert.NilError(t, err)
	assert.Equal(t, len(closures), 3)

	lodash := Package{
		Key:     "lodash@npm:4.17.21",
		Version: "4.17.21",
		Found:   true,
	}
	assert.Assert(t, !closures[nxpkgpath.AnchoredUnixPath("apps/web")].Contains(lodash))
	assert.Assert(t, closures[nxpkgpath.AnchoredUnixPath("apps/docs")].Contains(lodash))
}
