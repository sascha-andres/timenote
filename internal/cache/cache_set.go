package cache

// Set stores data for key and resets its expiry.
func Set[T any](c *Cache, key string, data T) error {
	return writeJSON(c.path(key), record[T]{
		Meta: newMetaData(c.maxAge),
		Data: data,
	})
}
