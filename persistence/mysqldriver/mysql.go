package mysqldriver

import (
	"context"
	"database/sql"

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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS timenote (
	id INT AUTO_INCREMENT NOT NULL,
	tag NVARCHAR(100) NOT NULL DEFAULT '',
	start TIMESTAMP NOT NULL DEFAULT NOW(),
	stop TIMESTAMP,
	PRIMARY KEY ( id )
);`)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating table - timenote")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS project (
	id INT AUTO_INCREMENT NOT NULL,
	name NVARCHAR(100) NOT NULL DEFAULT '',
	PRIMARY KEY ( id )
);`)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating table - project")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS timenote_project (
	id INT AUTO_INCREMENT NOT NULL,
	id_timenote INT NOT NULL,
	id_project INT NOT NULL,
	PRIMARY KEY ( id )
);`)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating table - timenote_project")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS line (
	id INT AUTO_INCREMENT NOT NULL,
	id_timenote INT NOT NULL,
	text TEXT NOT NULL DEFAULT '',
	entered TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY ( id )
);`)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating table - line")
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
	_, err = tx.Exec("update timenote set stop = NOW() where stop = '0000-00-00 00:00:00'")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not update old entry")
	}
	_, err = tx.Exec("insert into timenote set start = NOW()")
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
	var id int
	if err := mysql.databaseConnection.QueryRow("select id from timenote where stop = '0000-00-00 00:00:00'").Scan(&id); err != nil {
		return errors.Wrap(err, "Could not start get active id")
	}
	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not start transaction")
	}
	stmt, err := tx.Prepare("insert into line (id_timenote, text) values (?, ?)")
	if err != nil {
		return errors.Wrap(err, "Could not prepare statement")
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %#v\n", err)
		}
	}()
	_, err = stmt.Exec(id, line)
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
	if _, err = tx.Exec("update timenote set `tag`=? where `stop` = '0000-00-00 00:00:00'", name); err != nil {
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
	_, err = tx.Exec("update timenote set stop = NOW() where stop = '0000-00-00 00:00:00'")
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

/* func (mysql *MySQLPersistor) ListForDay(delta int) ([]timenote.TimeEntry, error) {
	if err := mysql.prepareDb(); err != nil {
		return nil, errors.Wrap(err, "Connection to DB not valid")
	}
	return nil, errors.New("Not yet implemented")
}*/

func (mysql *MySQLPersistor) Current() (*timenote.TimeEntry, error) {
	if err := mysql.prepareDb(); err != nil {
		return nil, errors.Wrap(err, "Connection to DB not valid")
	}
	row := mysql.databaseConnection.QueryRow("select id, tag, `text` from timenote where stop = '0000-00-00 00:00:00'")
	var te timenote.TimeEntry
	if err := row.Scan(&te.Id, &te.Tag, &te.Note); err != nil {
		if err == sql.ErrNoRows {
			return nil, timenote.ErrNoCurrentTimeEntry
		}
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

func (mysql *MySQLPersistor) Project(name string) error {
	if err := mysql.prepareDb(); err != nil {
		return errors.Wrap(err, "Connection to DB not valid")
	}

	var (
		projectID int
		err       error
		te        *timenote.TimeEntry
	)
	projectID, err = mysql.getProjectID(name)
	if err != nil {
		return errors.Wrap(err, "Unable to select project")
	}
	te, err = mysql.Current()
	if err != nil {
		return err
	}
	if projectID == 0 {
		projectID, err = mysql.getProjectID(name)
		if err != nil {
			return errors.Wrap(err, "Unable to select project")
		}
	}
	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not start transaction")
	}
	_, err = tx.Exec("insert into timenote_project (timenote_id, project_id) values (?, ?)", te.Id, projectID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Could not stop entry")
	}
	return tx.Commit()
}

func (mysql *MySQLPersistor) createProject(name string) (int, error) {
	if err := mysql.prepareDb(); err != nil {
		return 0, errors.Wrap(err, "Connection to DB not valid")
	}

	tx, err := mysql.databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, errors.Wrap(err, "Could not start transaction")
	}
	_, err = tx.Exec("insert into project (name) values (?)", name)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "Could not insert project")
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return mysql.getProjectID(name)
}

func (mysql *MySQLPersistor) getProjectID(name string) (int, error) {
	if err := mysql.prepareDb(); err != nil {
		return 0, errors.Wrap(err, "Connection to DB not valid")
	}
	row := mysql.databaseConnection.QueryRow("select id, `name` from project where name = ?", name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			return 0, errors.Wrap(err, "Could not load record")
		}
	}
	return id, nil
}
