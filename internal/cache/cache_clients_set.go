package cache

import "github.com/jason0x43/go-toggl"

// SetClients stores the given clients in the cache.
func (c *Cache) SetClients(workspace int, clients []toggl.Client) error {
	rec := clientsRecord{
		Clients: clients,
		Meta:    newMetaData(c.maxAge),
	}
	return writeJSON(clientsPath(c.dir, workspace), rec)
}
