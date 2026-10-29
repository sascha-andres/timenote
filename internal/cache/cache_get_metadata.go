package cache

// MetaData returns when key was updated and when it needs the next update.
func (c *Cache) MetaData(key string) (MetaData, error) {
	var rec struct{ Meta MetaData }
	err := readJSON(c.path(key), &rec)
	return rec.Meta, err
}
