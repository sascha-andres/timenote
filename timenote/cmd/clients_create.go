package cmd

import "context"

// clientsCreate creates a client in the current workspace.
func clientsCreate() {
	n := requireString("name", *name)
	p := newPersistor()
	if err := p.CreateClient(context.Background(), n); err != nil {
		fatal(err)
	}
}
