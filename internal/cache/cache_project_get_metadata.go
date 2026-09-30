package cache

// ProjectMetaData returns metadata about the projects cache.
func (c *Cache) ProjectMetaData(workspace int) (*MetaData, error) {
	var rec projectsRecord
	if err := readJSON(projectsPath(c.dir, workspace), &rec); err != nil {
		return nil, err
	}
	return &rec.Meta, nil
}
