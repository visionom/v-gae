package mysql

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/visionom/v-gae/adapter/mysql/config"
	"github.com/visionom/v-gae/adapter/mysql/domain"
	"github.com/visionom/v-gae/adapter/mysql/errors"
)

type DB struct {
	dbcp *sql.DB
}

var instance *DB
var once sync.Once

func GetDB() (*DB, error) {
	if instance == nil {
		return nil, errors.New("database connection pool is uninitialized")
	}
	return instance, nil
}

func Init(conf *config.DBConfig) error {
	var err error
	once.Do(func() {
		instance = &DB{}
		instance.dbcp, err = newDB(conf.User, conf.Passwd, conf.Host, conf.Name, conf.Port)
	})
	if err != nil {
		return err
	}

	if err = instance.dbcp.Ping(); err != nil {
		return errors.Wrap(err, "database connect fails")
	}
	return nil
}

func newDB(user, passwd, host, dbName string, port int) (*sql.DB, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd,
		host, port, dbName)
	cp, err := sql.Open("mysql", source)
	if err != nil {
		return nil, errors.Wrapf(err,
			"database connection open fails, source: %s", source)
	}
	return cp, nil
}

func Query(sql string, args ...interface{}) (domain.ActiveRecordList,
	error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.dbcp.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrapf(
			err, "sql query fails, sql: %s; args: %v", sql, args)
	}
	defer rows.Close()
	rds, err := domain.GetActiveRecordList(rows)
	if err != nil {
		return nil, err
	}
	return rds, nil
}

func QueryFirst(sql string, args ...interface{}) (*domain.ActiveRecord,
	error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.dbcp.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err,
			"sql query fails, sql: %s; args: %v", sql, args)
	}
	defer rows.Close()
	rd, err := domain.GetActiveRecord(rows)
	if err != nil {
		return nil, err
	}
	return rd, nil
}

func Update(sql string, args ...interface{}) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}
	stmt, err := db.dbcp.Prepare(sql)
	if err != nil {
		return -1, errors.Wrapf(err,
			"prepare sql fails, sql: %s; args: %v", sql, args)
	}
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return -1, errors.Wrapf(err,
			"exec sql fails, sql: %s; args: %v", sql, args)
	}
	affect, err := res.RowsAffected()

	if err != nil {
		return -1, errors.Wrapf(err,
			"read affect fails")
	}
	return affect, nil
}

func Stat() (sql.DBStats, error) {
	db, err := GetDB()
	if err != nil {
		return sql.DBStats{}, err
	}
	return db.dbcp.Stats(), nil
}

func Close() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	if err := db.dbcp.Close(); err != nil {
		return errors.Wrapf(err,
			"close fails")
	}
	return nil
}

func BeginTx() (*Tx, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return BeginTxByDB(db)
}
