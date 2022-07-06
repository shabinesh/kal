package main

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type taskrepo struct {
	db *sql.DB
}

var migrations = []string{
	`create table tasks (id integer primary key autoincrement not null, title text, ts timestamp, description blob);`,
}

func newTaskrepo() *taskrepo {
	c, err := sql.Open("sqlite3", ".cal.db")
	if err != nil {
		panic(err)
	}

	err = migrate(c, migrations)
	if err != nil {
		panic(err)
	}

	return &taskrepo{db: c}
}

func (t taskrepo) addTask(inp taskInput) error {
	_, err := t.db.Exec(`insert into tasks (title, ts, description) values ($1, $2, $3)`, inp.Title, inp.Timestamp, inp.Desc)

	return err
}

func (t taskrepo) getTasksOfDay(day time.Time) ([]taskItem, error) {
	items := make([]taskItem, 0)
	rows, err := t.db.Query(`select id, title, description, ts from tasks where date(ts) = date($1);`, day.Format("2006-01-02"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return items, nil
		}

		return items, err
	}

	for rows.Next() {
		var ti taskItem
		var id int
		err := rows.Scan(&id, &ti.title, &ti.desc, &ti.ts)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
		}

		items = append(items, ti)
	}

	return items, rows.Close()
}
