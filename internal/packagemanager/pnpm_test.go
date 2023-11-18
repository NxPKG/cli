package packagemanager

import (
	"os"
	"testing"

	"github.com/nxpkg/cli/internal/fs"
	"github.com/nxpkg/cli/internal/nxpkgpath"
	"gotest.tools/v3/assert"
)

func pnpmPatchesSection(t *testing.T, pkgJSON *fs.PackageJSON) map[string]interface{} {
	t.Helper()
	pnpmSection, ok := pkgJSON.RawJSON["pnpm"].(map[string]interface{})
	assert.Assert(t, ok)
	patchesSection, ok := pnpmSection["patchedDependencies"].(map[string]interface{})
	assert.Assert(t, ok)
	return patchesSection
}

func getPnpmPackageJSON(t *testing.T) *fs.PackageJSON {
	t.Helper()
	rawCwd, err := os.Getwd()
	assert.NilError(t, err)
	cwd, err := fs.CheckedToAbsoluteSystemPath(rawCwd)
	assert.NilError(t, err)
	pkgJSONPath := cwd.Join("fixtures", "pnpm-patches.json")
	pkgJSON, err := fs.ReadPackageJSON(pkgJSONPath)
	assert.NilError(t, err)
	return pkgJSON
}

func Test_PnpmPrunePatches_KeepsNecessary(t *testing.T) {
	pkgJSON := getPnpmPackageJSON(t)
	initialPatches := pnpmPatchesSection(t, pkgJSON)

	assert.DeepEqual(t, initialPatches, map[string]interface{}{"is-odd@3.0.1": "patches/is-odd@3.0.1.patch"})

	err := pnpmPrunePatches(pkgJSON, []nxpkgpath.AnchoredUnixPath{nxpkgpath.AnchoredUnixPath("patches/is-odd@3.0.1.patch")})
	assert.NilError(t, err)

	newPatches := pnpmPatchesSection(t, pkgJSON)
	assert.DeepEqual(t, newPatches, map[string]interface{}{"is-odd@3.0.1": "patches/is-odd@3.0.1.patch"})
}

func Test_PnpmPrunePatches_RemovesExtra(t *testing.T) {
	pkgJSON := getPnpmPackageJSON(t)
	initialPatches := pnpmPatchesSection(t, pkgJSON)

	assert.DeepEqual(t, initialPatches, map[string]interface{}{"is-odd@3.0.1": "patches/is-odd@3.0.1.patch"})

	err := pnpmPrunePatches(pkgJSON, nil)
	assert.NilError(t, err)

	newPatches := pnpmPatchesSection(t, pkgJSON)
	assert.DeepEqual(t, newPatches, map[string]interface{}{})
}
