package cache

// AccountNeedUpdate returns true if the cached account needs a refresh.
func (c *Cache) AccountNeedUpdate() bool {
	return needsRefresh(accountPath(c.dir))
}
