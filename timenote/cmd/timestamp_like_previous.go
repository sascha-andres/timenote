package cmd

import "context"

// timestampLikePrevious selects the last finished entry and starts a new
// one using it as a template.
func timestampLikePrevious() {
	p := newPersistor()
	if err := p.StartPrevious(context.Background()); err != nil {
		fatal(err)
	}
}
