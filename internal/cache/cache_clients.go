package cache

import "github.com/jason0x43/go-toggl"

// Clients returns the cached clients for workspace.
func (c *Cache) Clients(workspace int) ([]toggl.Client, error) {
	var rec clientsRecord
	if err := readJSON(clientsPath(c.dir, workspace), &rec); err != nil {
		return nil, err
	}
	return rec.Clients, nil
}
