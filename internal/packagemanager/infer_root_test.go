package packagemanager

import (
	"reflect"
	"testing"

	"github.com/nxpkg/cli/internal/nxpkgpath"
	"gotest.tools/v3/assert"
)

func TestInferRoot(t *testing.T) {
	type file struct {
		path    nxpkgpath.AnchoredSystemPath
		content []byte
	}

	tests := []struct {
		name               string
		fs                 []file
		executionDirectory nxpkgpath.AnchoredSystemPath
		rootPath           nxpkgpath.AnchoredSystemPath
		packageMode        PackageType
	}{
		// Scenario 0
		{
			name: "nxpkg.json at current dir, no package.json",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "nxpkg.json at parent dir, no package.json",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			// This is "no inference"
			rootPath:    nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			packageMode: Multi,
		},
		// Scenario 1A
		{
			name: "nxpkg.json at current dir, has package.json, has workspaces key",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "nxpkg.json at parent dir, has package.json, has workspaces key",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "nxpkg.json at parent dir, has package.json, has pnpm workspaces",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("pnpm-workspace.yaml").ToSystemPath(),
					content: []byte("packages:\n  - docs"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		// Scenario 1A aware of the weird thing we do for packages.
		{
			name: "nxpkg.json at current dir, has package.json, has packages key",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Single,
		},
		{
			name: "nxpkg.json at parent dir, has package.json, has packages key",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Single,
		},
		// Scenario 1A aware of the the weird thing we do for packages when both methods of specification exist.
		{
			name: "nxpkg.json at current dir, has package.json, has workspace and packages key",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"clobbered\" ], \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "nxpkg.json at parent dir, has package.json, has workspace and packages key",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"clobbered\" ], \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		// Scenario 1B
		{
			name: "nxpkg.json at current dir, has package.json, no workspaces",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{}"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Single,
		},
		{
			name: "nxpkg.json at parent dir, has package.json, no workspaces",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{}"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Single,
		},
		{
			name: "nxpkg.json at parent dir, has package.json, no workspaces, includes pnpm",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{path: nxpkgpath.AnchoredUnixPath("nxpkg.json").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("pnpm-workspace.yaml").ToSystemPath(),
					content: []byte(""),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Single,
		},
		// Scenario 2A
		{
			name:               "no nxpkg.json, no package.json at current",
			fs:                 []file{},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "no nxpkg.json, no package.json at parent",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			packageMode:        Multi,
		},
		// Scenario 2B
		{
			name: "no nxpkg.json, has package.json with workspaces at current",
			fs: []file{
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "no nxpkg.json, has package.json with workspaces at parent",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			packageMode:        Multi,
		},
		{
			name: "no nxpkg.json, has package.json with pnpm workspaces at parent",
			fs: []file{
				{path: nxpkgpath.AnchoredUnixPath("execution/path/subdir/.file").ToSystemPath()},
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("pnpm-workspace.yaml").ToSystemPath(),
					content: []byte("packages:\n  - docs"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("execution/path/subdir").ToSystemPath(),
			packageMode:        Multi,
		},
		// Scenario 3A
		{
			name: "no nxpkg.json, lots of package.json files but no workspaces",
			fs: []file{
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/two/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/two/three/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("one/two/three").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("one/two/three").ToSystemPath(),
			packageMode:        Single,
		},
		// Scenario 3BI
		{
			name: "no nxpkg.json, lots of package.json files, and a workspace at the root that matches execution directory",
			fs: []file{
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"one/two/three\" ] }"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/two/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/two/three/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("one/two/three").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("one/two/three").ToSystemPath(),
			packageMode:        Multi,
		},
		// Scenario 3BII
		{
			name: "no nxpkg.json, lots of package.json files, and a workspace at the root that matches execution directory",
			fs: []file{
				{
					path:    nxpkgpath.AnchoredUnixPath("package.json").ToSystemPath(),
					content: []byte("{ \"workspaces\": [ \"does-not-exist\" ] }"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/two/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
				{
					path:    nxpkgpath.AnchoredUnixPath("one/two/three/package.json").ToSystemPath(),
					content: []byte("{}"),
				},
			},
			executionDirectory: nxpkgpath.AnchoredUnixPath("one/two/three").ToSystemPath(),
			rootPath:           nxpkgpath.AnchoredUnixPath("one/two/three").ToSystemPath(),
			packageMode:        Single,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsRoot := nxpkgpath.AbsoluteSystemPath(t.TempDir())
			for _, file := range tt.fs {
				path := file.path.RestoreAnchor(fsRoot)
				assert.NilError(t, path.Dir().MkdirAll(0777))
				assert.NilError(t, path.WriteFile(file.content, 0777))
			}

			nxpkgRoot, packageMode := InferRoot(tt.executionDirectory.RestoreAnchor(fsRoot))
			if !reflect.DeepEqual(nxpkgRoot, tt.rootPath.RestoreAnchor(fsRoot)) {
				t.Errorf("InferRoot() nxpkgRoot = %v, want %v", nxpkgRoot, tt.rootPath.RestoreAnchor(fsRoot))
			}
			if packageMode != tt.packageMode {
				t.Errorf("InferRoot() packageMode = %v, want %v", packageMode, tt.packageMode)
			}
		})
	}
}
