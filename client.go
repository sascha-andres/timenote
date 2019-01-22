package timenote

import "fmt"

type Client struct {
	ID          int
	Name        string
	Description string
}

func (c Client) String() string {
	return fmt.Sprintf("%d: %s\nNote:\n%s\n", c.ID, c.Name, c.Description)
}
