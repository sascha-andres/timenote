package cache

import "time"

// NeedUpdate returns true if key is not cached or has expired.
func (c *Cache) NeedUpdate(key string) bool {
	md, err := c.MetaData(key)
	if err != nil || md.NextUpdate.IsZero() {
		return true
	}
	return time.Now().After(md.NextUpdate)
}
