package cmd

// timestampLikePrevious selects the last finished entry and starts a new
// one using it as a template.
func timestampLikePrevious() {
	p := newPersistor()
	if err := p.StartPrevious(); err != nil {
		fatal(err)
	}
}
