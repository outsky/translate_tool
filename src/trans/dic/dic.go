package dic

import (
	"bytes"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sql_create string = `create table if not exists dic (cn bold not null primary key, trans bold not null);`
	sql_insert string = `insert into dic(trans, cn) values(?, ?)`
	sql_update string = `update dic set trans = ? where cn = ?`
	sql_query  string = `select trans from dic where cn = ?`
)

type dic struct {
	db *sql.DB
}

func New(name string) *dic {
	ins := &dic{}
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(sql_create)
	if err != nil {
		panic(err)
	}
	ins.db = db
	return ins
}

func (d *dic) Close() {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			panic(err)
		}
		d.db = nil
	}
}

func (d *dic) Insert(cn, trans []byte) error {
	sql := sql_insert
	ret, err := d.Query(cn)
	if err == nil {
		if bytes.Compare(ret, trans) == 0 {
			return nil
		}
		sql = sql_update
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(trans, cn)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *dic) Query(text []byte) ([]byte, error) {
	var trans []byte
	stmt, err := d.db.Prepare(sql_query)
	if err != nil {
		return trans, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(text).Scan(&trans)
	if err != nil {
		return trans, err
	}
	return trans, nil
}
