package mysqldriver

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sascha-andres/timenote"
	"github.com/sascha-andres/timenote/persistence"
	log "github.com/sirupsen/logrus"
)

type (
	MySQLPersistor struct {
		databaseConnection *sql.DB
		dsn                string
	}
)

func NewMysql(dsn string) (persistence.Persistor, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MySQL")
	}
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(500)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS timenotes (
	id INT AUTO_INCREMENT NOT NULL,
	tag NVARCHAR(100) NOT NULL DEFAULT '',
	text TEXT NOT NULL DEFAULT '',
	start TIMESTAMP NOT NULL DEFAULT NOW(),
	stop TIMESTAMP,
	PRIMARY KEY ( id )
);`)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating table")
	}
	return &MySQLPersistor{databaseConnection: db, dsn: dsn}, nil
}

func (mysql *MySQLPersistor) New() error {
	if err := mysql.prepareDb(); err != nil {
		return errors.Wrap(err, "Connection to DB not valid")
	}
	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not start transaction")
	}
	_, err = tx.Exec("update timenotes set stop = NOW() where stop = '0000-00-00 00:00:00'")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not update old entry")
	}
	_, err = tx.Exec("insert into timenotes set start = NOW()")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not create new entry")
	}
	return tx.Commit()
}

func (mysql *MySQLPersistor) Append(line string) error {
	if err := mysql.prepareDb(); err != nil {
		return errors.Wrap(err, "Connection to DB not valid")
	}
	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not start transaction")
	}
	stmt, err := tx.Prepare("update timenotes set `text`=CONCAT(`text`, ?) where stop = '0000-00-00 00:00:00'")
	if err != nil {
		return errors.Wrap(err, "Could not prepare statement")
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %#v\n", err)
		}
	}()
	_, err = stmt.Exec(fmt.Sprintf("\n%s", line))
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not append line")
	}
	return tx.Commit()
}

func (mysql *MySQLPersistor) Tag(name string) error {
	if err := mysql.prepareDb(); err != nil {
		return errors.Wrap(err, "Connection to DB not valid")
	}
	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not start transaction")
	}
	if _, err = tx.Exec("update timenotes set `tag`=? where `stop` = '0000-00-00 00:00:00'", name); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not append line")
	}
	return tx.Commit()
}

func (mysql *MySQLPersistor) Done() error {
	if err := mysql.prepareDb(); err != nil {
		return errors.Wrap(err, "Connection to DB not valid")
	}
	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not start transaction")
	}
	_, err = tx.Exec("update timenotes set stop = NOW() where stop = '0000-00-00 00:00:00'")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not stop entry")
	}
	return tx.Commit()
}

func (mysql *MySQLPersistor) Close() error {
	return mysql.databaseConnection.Close()
}

// duration in sec
// SELECT *, UNIX_TIMESTAMP(`stop`)-UNIX_TIMESTAMP(`start`) FROM `timenotes` WHERE `stop` <> '0000-00-00 00:00:00'
// unix_timestamp(maketime(_,_,_)

func (mysql *MySQLPersistor) ListForDay(delta int) ([]timenote.TimeEntry, error) {
	if err := mysql.prepareDb(); err != nil {
		return nil, errors.Wrap(err, "Connection to DB not valid")
	}
	return nil, errors.New("Not yet implemented")
}

func (mysql *MySQLPersistor) Current() (*timenote.TimeEntry, error) {
	if err := mysql.prepareDb(); err != nil {
		return nil, errors.Wrap(err, "Connection to DB not valid")
	}
	row := mysql.databaseConnection.QueryRow("select tag, `text` from timenotes where stop = '0000-00-00 00:00:00'")
	var te timenote.TimeEntry
	if err := row.Scan(&te.Tag, &te.Note); err != nil {
		return nil, errors.Wrap(err, "Could not load record")
	}
	return &te, nil
}

func (mysql *MySQLPersistor) prepareDb() error {
	if err := mysql.databaseConnection.Ping(); err != nil {
		db, err := sql.Open("mysql", mysql.dsn)
		if err != nil {
			return errors.Wrap(err, "Could not connect to MySQL")
		}
		mysql.databaseConnection = db
	}
	return nil
}
