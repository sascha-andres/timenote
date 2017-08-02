package factory

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sascha-andres/timenote/persistence"
	"github.com/sascha-andres/timenote/persistence/mysqldriver"
	"github.com/sascha-andres/timenote/persistence/toggldriver"
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
