package cache

func (c *Cache) Close() error {
	return c.db.Close()
}
