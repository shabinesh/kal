package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	createTableIfNotExists = `create table if not exists migration (version int, dirty bool);`
	setDirty               = `update table migration set dirty = true;`
	updateVersion          = `update table migration set version = $1`
	insertFirstMigration    = `insert into migration (version, dirty) values (1, false);`
)

func migrate(conn *sql.DB, migrations []string) error {
	if len(migrations) == 0 {
		return fmt.Errorf("empty migration set")
	}

	conn.ExecContext(context.Background(), createTableIfNotExists)

	row := conn.QueryRowContext(context.Background(), `select version, dirty from migration`)

	var (
		version int
		dirty   bool
	)

	if err := row.Scan(&version, &dirty); err != nil {
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		logger.Println("insert zero migration")
		conn.ExecContext(context.Background(), insertFirstMigration)
	}

	if dirty {
		panic("last migration was is dirty")
	}

	logger.Println(version)
	for i, m := range migrations[version:] {
		_, err := conn.ExecContext(context.Background(), m)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			conn.ExecContext(context.Background(), setDirty)
			return err
		}
			
		conn.ExecContext(context.Background(), updateVersion, i)
	}

	return nil
}
