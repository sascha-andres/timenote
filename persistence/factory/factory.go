package factory

import (
	"github.com/pkg/errors"
	"github.com/sascha-andres/timenote/persistence"
	"github.com/sascha-andres/timenote/persistence/mysqldriver"
)

// CreatePersistence returns the selected backend
func CreatePersistence(driver, dsn string) (persistence.Persistor, error) {
	if driver != "mysql" {
		return nil, errors.New("Only a MySQL driver exists")
	}
	return mysqldriver.NewMysql(dsn)
}
