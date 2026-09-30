package cmd

// clientsCreate creates a client in the current workspace.
func clientsCreate() {
	n := requireString("name", *name)
	p := newPersistor()
	if err := p.CreateClient(n); err != nil {
		fatal(err)
	}
}
