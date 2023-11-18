package packagemanager

import (
	"fmt"

	"github.com/nxpkg/cli/internal/fs"
	"github.com/nxpkg/cli/internal/lockfile"
	"github.com/nxpkg/cli/internal/nxpkgpath"
)

const npmLockfile = "package-lock.json"

var nodejsNpm = PackageManager{
	Name:         "nodejs-npm",
	Slug:         "npm",
	Command:      "npm",
	Specfile:     "package.json",
	Lockfile:     npmLockfile,
	PackageDir:   "node_modules",
	ArgSeparator: func(_userArgs []string) []string { return []string{"--"} },

	getWorkspaceGlobs: func(rootpath nxpkgpath.AbsoluteSystemPath) ([]string, error) {
		pkg, err := fs.ReadPackageJSON(rootpath.UntypedJoin("package.json"))
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. Nxpkgrepo requires npm workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},

	getWorkspaceIgnores: func(pm PackageManager, rootpath nxpkgpath.AbsoluteSystemPath) ([]string, error) {
		// Matches upstream values:
		// function: https://github.com/npm/map-workspaces/blob/a46503543982cb35f51cc2d6253d4dcc6bca9b32/lib/index.js#L73
		// key code: https://github.com/npm/map-workspaces/blob/a46503543982cb35f51cc2d6253d4dcc6bca9b32/lib/index.js#L90-L96
		// call site: https://github.com/npm/cli/blob/7a858277171813b37d46a032e49db44c8624f78f/lib/workspaces/get-workspaces.js#L14
		return []string{
			"**/node_modules/**",
		}, nil
	},

	canPrune: func(cwd nxpkgpath.AbsoluteSystemPath) (bool, error) {
		return true, nil
	},

	GetLockfileName: func(_ nxpkgpath.AbsoluteSystemPath) string {
		return npmLockfile
	},

	GetLockfilePath: func(projectDirectory nxpkgpath.AbsoluteSystemPath) nxpkgpath.AbsoluteSystemPath {
		return projectDirectory.UntypedJoin(npmLockfile)
	},

	GetLockfileContents: func(projectDirectory nxpkgpath.AbsoluteSystemPath) ([]byte, error) {
		return projectDirectory.UntypedJoin(npmLockfile).ReadFile()
	},

	UnmarshalLockfile: func(_rootPackageJSON *fs.PackageJSON, contents []byte) (lockfile.Lockfile, error) {
		return lockfile.DecodeNpmLockfile(contents)
	},
}
