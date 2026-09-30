package cache

import "github.com/jason0x43/go-toggl"

// ClientByID returns the cached client with the given id.
func (c *Cache) ClientByID(clientID, workspace int) (*toggl.Client, error) {
	clients, err := c.Clients(workspace)
	if err != nil {
		return nil, err
	}
	for _, cl := range clients {
		if cl.ID == clientID {
			return &cl, nil
		}
	}
	return &toggl.Client{}, nil
}
