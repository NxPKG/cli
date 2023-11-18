package packagemanager

import (
	"fmt"
	"strings"

	"github.com/nxpkg/cli/internal/fs"
	"github.com/nxpkg/cli/internal/lockfile"
	"github.com/nxpkg/cli/internal/nxpkgpath"
)

const berryLockfile = "yarn.lock"

var nodejsBerry = PackageManager{
	Name:         "nodejs-berry",
	Slug:         "yarn",
	Command:      "yarn",
	Specfile:     "package.json",
	Lockfile:     berryLockfile,
	PackageDir:   "node_modules",
	ArgSeparator: func(_userArgs []string) []string { return nil },

	getWorkspaceGlobs: func(rootpath nxpkgpath.AbsoluteSystemPath) ([]string, error) {
		pkg, err := fs.ReadPackageJSON(rootpath.UntypedJoin("package.json"))
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. Nxpkgrepo requires Yarn workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},

	getWorkspaceIgnores: func(pm PackageManager, rootpath nxpkgpath.AbsoluteSystemPath) ([]string, error) {
		// Matches upstream values:
		// Key code: https://github.com/yarnpkg/berry/blob/8e0c4b897b0881878a1f901230ea49b7c8113fbe/packages/yarnpkg-core/sources/Workspace.ts#L64-L70
		return []string{
			"**/node_modules",
			"**/.git",
			"**/.yarn",
		}, nil
	},

	canPrune: func(cwd nxpkgpath.AbsoluteSystemPath) (bool, error) {
		return true, nil
	},

	GetLockfileName: func(_ nxpkgpath.AbsoluteSystemPath) string {
		return berryLockfile
	},

	GetLockfilePath: func(projectDirectory nxpkgpath.AbsoluteSystemPath) nxpkgpath.AbsoluteSystemPath {
		return projectDirectory.UntypedJoin(berryLockfile)
	},

	GetLockfileContents: func(projectDirectory nxpkgpath.AbsoluteSystemPath) ([]byte, error) {
		return projectDirectory.UntypedJoin(berryLockfile).ReadFile()
	},

	UnmarshalLockfile: func(rootPackageJSON *fs.PackageJSON, contents []byte) (lockfile.Lockfile, error) {
		var resolutions map[string]string
		if untypedResolutions, ok := rootPackageJSON.RawJSON["resolutions"]; ok {
			if untypedResolutions, ok := untypedResolutions.(map[string]interface{}); ok {
				resolutions = make(map[string]string, len(untypedResolutions))
				for resolution, reference := range untypedResolutions {
					if reference, ok := reference.(string); ok {
						resolutions[resolution] = reference
					}
				}
			}
		}

		return lockfile.DecodeBerryLockfile(contents, resolutions)
	},

	prunePatches: func(pkgJSON *fs.PackageJSON, patches []nxpkgpath.AnchoredUnixPath) error {
		pkgJSON.Mu.Lock()
		defer pkgJSON.Mu.Unlock()

		keysToDelete := []string{}
		resolutions, ok := pkgJSON.RawJSON["resolutions"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Invalid structure for resolutions field in package.json")
		}

		for dependency, untypedPatch := range resolutions {
			inPatches := false
			patch, ok := untypedPatch.(string)
			if !ok {
				return fmt.Errorf("Expected value of %s in package.json to be a string, got %v", dependency, untypedPatch)
			}

			for _, wantedPatch := range patches {
				if strings.HasSuffix(patch, wantedPatch.ToString()) {
					inPatches = true
					break
				}
			}

			// We only want to delete unused patches as they are the only ones that throw if unused
			if !inPatches && strings.HasSuffix(patch, ".patch") {
				keysToDelete = append(keysToDelete, dependency)
			}
		}

		for _, key := range keysToDelete {
			delete(resolutions, key)
		}

		return nil
	},
}
