package timenote

import "fmt"

type Client struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"notes"`
}

func (c Client) String() string {
	return fmt.Sprintf("%d: %s\nNote:\n%s\n", c.ID, c.Name, c.Description)
}
