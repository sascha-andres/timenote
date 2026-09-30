package cmd

import "fmt"

// dispatchToken routes "token" sub-verbs: save, delete.
func dispatchToken(verbs []string) {
	if len(verbs) == 0 {
		fatal(fmt.Errorf("token requires a command: save, delete"))
	}
	switch verbs[0] {
	case "save":
		tokenSave()
	case "delete":
		tokenDelete()
	default:
		fatal(fmt.Errorf("unknown token command %q", verbs[0]))
	}
}
