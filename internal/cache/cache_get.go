package cache

// Get returns the cached data for key. If key is not cached, it returns the
// zero value of T.
func Get[T any](c *Cache, key string) (T, error) {
	var rec record[T]
	err := readJSON(c.path(key), &rec)
	return rec.Data, err
}
