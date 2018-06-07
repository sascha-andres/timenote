package factory

import (
	"fmt"

	"github.com/pkg/errors"
	"livingit.de/code/timenote/persistence"
	"livingit.de/code/timenote/persistence/mysqldriver"
	"livingit.de/code/timenote/persistence/toggldriver"
)

// CreatePersistence returns the selected backend
func CreatePersistence(driver, dsn string) (persistence.Persistor, error) {
	if driver == "mysql" {
		return mysqldriver.NewMysql(dsn)
	}
	if driver == "toggl" {
		return toggldriver.NewToggl(dsn)
	}
	return nil, errors.New(fmt.Sprintf("Driver %s does not exist", driver))
}
