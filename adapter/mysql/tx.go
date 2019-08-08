package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/visionom/v-gae/adapter/mysql/domain"
	"github.com/visionom/v-gae/adapter/mysql/errors"
)

type Tx struct {
	tx *sql.Tx
}

func BeginTxByDB(db *DB) (*Tx, error) {
	tx, err := db.dbcp.Begin()
	if err != nil {
		wrapErr := errors.Wrap(err,
			"The transaction start fails")
		if txErr := tx.Rollback(); txErr != nil {
			return nil, errors.Wrapf(txErr, "rollback fails, when catch fails: %s", wrapErr.Error())
		}
		return nil, wrapErr
	}
	return &Tx{
		tx,
	}, nil
}

func (tx *Tx) Query(sql string, args ...interface{}) (domain.ActiveRecordList,
	error) {
	rows, err := tx.tx.Query(sql, args...)
	if err != nil {
		wrapErr := errors.Wrapf(err,
			"sql query fails, sql: %s; args: %v", sql, args)
		if txErr := tx.tx.Rollback(); txErr != nil {
			return nil, errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", wrapErr.Error())
		}
		return nil, wrapErr
	}
	defer rows.Close()

	rds, err := domain.GetActiveRecordList(rows)
	if err != nil {
		wrapErr := err
		if txErr := tx.tx.Rollback(); txErr != nil {
			return nil, errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", wrapErr.Error())
		}
		return nil, wrapErr
	}

	return rds, nil
}

func (tx *Tx) QueryFirst(sql string, args ...interface{}) (*domain.ActiveRecord,
	error) {
	rows, err := tx.tx.Query(sql, args...)
	if err != nil {
		if txErr := tx.tx.Rollback(); txErr != nil {
			return nil, errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", err.Error())
		}
		return nil, err
	}
	defer rows.Close()

	rd, err := domain.GetActiveRecord(rows)
	if err != nil {
		if txErr := tx.tx.Rollback(); txErr != nil {
			return nil, errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", err.Error())
		}
		return nil, err
	}
	return rd, nil
}

func (tx *Tx) Update(sql string, args ...interface{}) (int64, error) {
	stmt, err := tx.tx.Prepare(sql)
	if err != nil {
		wrapErr := errors.Wrapf(err,
			"prepare sql fails, sql: %s; args: %v", sql, args)
		if txErr := tx.tx.Rollback(); txErr != nil {
			return -1, errors.Wrapf(txErr, "rollback fails! when catch fails: %s", wrapErr.Error())
		}
		return -1, wrapErr
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		wrapErr := errors.Wrapf(err,
			"exec sql fails, sql: %s; args: %v", sql, args)
		if txErr := tx.tx.Rollback(); txErr != nil {
			return -1, errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", wrapErr.Error())
		}
		return -1, wrapErr
	}

	affect, err := res.RowsAffected()
	if err != nil {
		wrapErr := errors.Wrapf(err,
			"read affect fails")
		if txErr := tx.tx.Rollback(); txErr != nil {
			return -1, errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", wrapErr.Error())
		}
		return -1, wrapErr
	}

	return affect, nil
}

func (tx *Tx) RollBack() error {
	if err := tx.tx.Rollback(); err != nil {
		return errors.Wrap(err, "rollback fails")
	}
	return nil
}

func (tx *Tx) Commit() error {
	if err := tx.tx.Commit(); err != nil {
		wrapErr := errors.Wrap(err, "commit fails")
		if txErr := tx.tx.Rollback(); txErr != nil {
			return errors.Wrapf(txErr,
				"rollback fails, when catch fails: %s", wrapErr.Error())
		}
		return wrapErr
	}
	return nil
}
