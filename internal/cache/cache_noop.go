package cache

import "github.com/nxpkg/cli/internal/nxpkgpath"

type noopCache struct{}

func newNoopCache() *noopCache {
	return &noopCache{}
}

func (c *noopCache) Put(_ nxpkgpath.AbsoluteSystemPath, _ string, _ int, _ []nxpkgpath.AnchoredSystemPath) error {
	return nil
}
func (c *noopCache) Fetch(_ nxpkgpath.AbsoluteSystemPath, _ string, _ []string) (ItemStatus, []nxpkgpath.AnchoredSystemPath, error) {
	return NewCacheMiss(), nil, nil
}

func (c *noopCache) Exists(_ string) ItemStatus {
	return NewCacheMiss()
}

func (c *noopCache) Clean(_ nxpkgpath.AbsoluteSystemPath) {}
func (c *noopCache) CleanAll()                            {}
func (c *noopCache) Shutdown()                            {}
