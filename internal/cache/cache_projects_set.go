package cache

import "github.com/jason0x43/go-toggl"

// SetProjects stores the given projects in the cache.
func (c *Cache) SetProjects(workspace int, projects []toggl.Project) error {
	rec := projectsRecord{
		Projects: projects,
		Meta:     newMetaData(c.maxAge),
	}
	return writeJSON(projectsPath(c.dir, workspace), rec)
}
