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
		"SHORTURL VARCHAR(255) PRIMARY KEY, " +
		"ORIGINALURL VARCHAR(255), " +
		"EXPIRATIONTIME INT);"
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

// 向数据库添加短链映射
func (d *database) addMapping(m Mapping) error {
	insSQL := "INSERT INTO %s VALUES (?, ?, ?);"
	stm, err := d.db.Prepare(fmt.Sprintf(insSQL, d.name))
	if err != nil {
		return err
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	_, err = stm.Exec(m.ShortURL, m.OriginalURL, m.ExpirationTime)
	if err != nil {
		return err
	}
	return nil
}

// 在数据库里使用短链获取原始链接
func (d *database) getMappingByS(shortURL string) (Mapping, error) {
	var m Mapping
	getSQL := "SELECT * FROM %s WHERE SHORTURL = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(getSQL, d.name))
	if err != nil {
		return m, err
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	err = stm.QueryRow(shortURL).Scan(&m.ShortURL, &m.OriginalURL, &m.ExpirationTime)
	if err != nil {
		return m, err
	}
	return m, err
}

// 在数据库里使用原始链接获取短链
func (d *database) getMappingByO(originalURL string) (Mapping, error) {
	var m Mapping
	getSQL := "SELECT * FROM %s WHERE ORIGINALURL = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(getSQL, d.name))
	if err != nil {
		return m, err
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	err = stm.QueryRow(originalURL).Scan(&m.ShortURL, &m.OriginalURL, &m.ExpirationTime)
	if err != nil {
		return m, err
	}
	return m, err
}

// 检测数据库里是否存在该原始链接
func (d *database) hasOriginalURL(originalURL string) bool {
	hasSQL := "SELECT * FROM %s WHERE ORIGINALURL = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(hasSQL, d.name))
	if err != nil {
		return false
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	rows, err := stm.Query(originalURL)
	if err != nil {
		return false
	}
	if rows.Err() != nil {
		return false
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if !rows.Next() {
		return false
	}
	_ = rows.Close()
	return true
}

// 检测数据库里是否存在该短链
func (d *database) hasShortURL(shortURL string) bool {
	hasSQL := "SELECT * FROM %s WHERE SHORTURL = ?;"
	stm, err := d.db.Prepare(fmt.Sprintf(hasSQL, d.name))
	if err != nil {
		return false
	}
	defer func(stm *sql.Stmt) {
		_ = stm.Close()
	}(stm)
	rows, err := stm.Query(shortURL)
	if err != nil {
		return false
	}
	if rows.Err() != nil {
		return false
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if !rows.Next() {
		return false
	}
	_ = rows.Close()
	return true
}

// 获取数据库里过期的映射
// 参数 time 为某时刻的时间戳
func (d *database) getRemove(time int64) ([]Mapping, error) {
	var ms []Mapping
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
		var m Mapping
		_ = rows.Scan(&m.ShortURL, &m.OriginalURL, &m.ExpirationTime)
		ms = append(ms, m)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	return ms, nil
}

// 删除数据库里过期的映射
// 参数 time 为某时刻的时间戳
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
