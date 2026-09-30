package cache

import "github.com/jason0x43/go-toggl"

// ProjectByID returns the cached project with the given id.
func (c *Cache) ProjectByID(projectID, workspace int) (*toggl.Project, error) {
	projects, err := c.Projects(workspace)
	if err != nil {
		return nil, err
	}
	for _, p := range projects {
		if p.ID == projectID {
			return &p, nil
		}
	}
	return &toggl.Project{}, nil
}
