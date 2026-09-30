package cache

// ClientMetaData returns metadata about the clients cache.
func (c *Cache) ClientMetaData(workspace int) (*MetaData, error) {
	var rec clientsRecord
	if err := readJSON(clientsPath(c.dir, workspace), &rec); err != nil {
		return nil, err
	}
	return &rec.Meta, nil
}
