package cache

import "github.com/jason0x43/go-toggl"

// GetAccount returns the cached account.
func (c *Cache) GetAccount() (toggl.Account, error) {
	var rec accountRecord
	if err := readJSON(accountPath(c.dir), &rec); err != nil {
		return toggl.Account{}, err
	}
	return rec.Account, nil
}
