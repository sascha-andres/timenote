package cache

import "github.com/jason0x43/go-toggl"

// AccountSet updates the cached account.
func (c *Cache) AccountSet(account *toggl.Account) error {
	rec := accountRecord{
		Account: *account,
		Meta:    newMetaData(c.maxAge),
	}
	return writeJSON(accountPath(c.dir), rec)
}
