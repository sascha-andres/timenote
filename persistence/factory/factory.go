package factory

import (
	"livingit.de/code/timenote/persistence/toggldriver"

	"livingit.de/code/timenote/persistence"
)

// CreatePersistence returns the selected backend
func CreatePersistence(dsn string) (persistence.Persistor, error) {
	return toggldriver.NewToggl(dsn, 0)
}
