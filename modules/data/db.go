package data

import (
	"database/sql"
	"errors"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"strconv"
)

type database struct {
	path string
	name string
	db   *sql.DB
}

func getDatabase() database {
	wd, _ := os.Getwd()
	p := path.Join(wd, "ssug.db")
	return database{
		name: "mappings",
		path: p,
	}
}

func (d *database) open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", d.path)
	if err != nil {
		return db, errors.New("open sqlite error " + err.Error())
	}
	createSql := "CREATE TABLE IF NOT EXISTS %s (" +
		"KEY VARCHAR(255) PRIMARY KEY," +
		"VALUE VARCHAR(255)," +
		"EXPIRATIONTIME INT" +
		");"
	stm, err := db.Prepare(fmt.Sprintf(createSql, d.name))
	if err != nil {
		return nil, errors.New("create sqlite table error " + err.Error())
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	_, err = stm.Exec()
	if err != nil {
		return nil, errors.New("create sqlite table error " + err.Error())
	}
	return db, nil
}

func (d *database) close() {
	_ = d.db.Close()
}

func (d *database) addMapping(time int64, key string, value string) error {
	insSQL := "INSERT INTO %s VALUES (?, ?, ?);"
	stm, err := d.db.Prepare(fmt.Sprintf(insSQL, d.name))
	if err != nil {
		return err
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	_, err = stm.Exec(key, value, time)
	if err != nil {
		return err
	}
	return nil
}

func (d *database) getMappingByV(value string) (error, mapping) {
	var m mapping
	getSQL := "SELECT * FROM %s WHERE VALUE = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(getSQL, d.name))
	if err != nil {
		return err, m
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	err = stm.QueryRow(value).Scan(&m.Key, &m.Value, &m.ExpirationTime)
	if err != nil {
		return err, m
	}
	return nil, m
}

func (d *database) getMappingByK(key string) (error, mapping) {
	var m mapping
	getSQL := "SELECT * FROM %s WHERE KEY = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(getSQL, d.name))
	if err != nil {
		return err, m
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	err = stm.QueryRow(key).Scan(&m.Key, &m.Value, &m.ExpirationTime)
	if err != nil {
		return err, m
	}
	return nil, m
}

func (d *database) hasKey(key string) bool {
	hasSQL := "SELECT * FROM %s WHERE KEY = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(hasSQL, d.name))
	if err != nil {
		return false
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	rows, err := stm.Query(key)
	if err != nil {
		return false
	}
	if rows.Err() != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	return true
}

func (d *database) hasValue(value string) bool {
	hasSQL := "SELECT * FROM %s WHERE VALUE = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(hasSQL, d.name))
	if err != nil {
		return false
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	rows, err := stm.Query(value)
	if err != nil {
		return false
	}
	if rows.Err() != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	return true
}

func (d *database) getRemove(time int64) ([]mapping, error) {
	var ms []mapping
	qSQL := "SELECT * FROM %s WHERE EXPIRATIONTIME <= ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(qSQL, d.name))
	if err != nil {
		return ms, err
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	rows, err := stm.Query(strconv.FormatInt(time, 10))
	if err != nil {
		return ms, err
	}
	for rows.Next() {
		var m mapping
		_ = rows.Scan(&m.Key, &m.Value, &m.ExpirationTime)
		ms = append(ms, m)
	}
	return ms, nil
}

func (d *database) autoRemove(time int64) error {
	delSQL := "DELETE FROM %s WHERE EXPIRATIONTIME <= ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(delSQL, d.name))
	if err != nil {
		return err
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	_, err = stm.Exec(strconv.FormatInt(time, 10))
	if err != nil {
		return err
	}
	return nil
}