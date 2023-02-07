package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type sqlite struct {
	db *sql.DB
}

func NewSqLite(source string) (*sqlite, error) {
	s := &sqlite{}

	var err error = nil

	s.db, err = sql.Open("sqlite", source)
	if err != nil {
		return s, fmt.Errorf("can't open db: %w", err)
	}

	var new bool

	rows, err := s.db.Query("select iif(count(name) = 0, 'true', 'false') from sqlite_master where type='table' and name='pings'")
	if err != nil {
		return nil, fmt.Errorf("can't execute query for check existing of table 'pings': %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&new); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	if new {
		_, err = s.db.Exec("create table pings (host text, time integer, created_at datetime)")
	}

	return s, nil
}

func (s *sqlite) Write(host string, time int64) error {
	stmt, _ := s.db.Prepare("insert into pings(host, time, created_at) values(?, ?, datetime())")
	defer stmt.Close()

	if _, err := stmt.Exec(host, time); err != nil {
		return fmt.Errorf("can't write data: %w", err)
	}

	return nil
}
