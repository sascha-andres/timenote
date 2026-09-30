package cache

import "github.com/jason0x43/go-toggl"

// Projects returns the cached projects for workspace.
func (c *Cache) Projects(workspace int) ([]toggl.Project, error) {
	var rec projectsRecord
	if err := readJSON(projectsPath(c.dir, workspace), &rec); err != nil {
		return nil, err
	}
	return rec.Projects, nil
}
