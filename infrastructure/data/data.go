package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"os"
)

type Data struct {
	Ctx context.Context
	db  *sql.DB
}

func NewData(ctx context.Context, db *sql.DB) *Data {
	if db == nil {
		panic("db is nil")
	}
	return &Data{Ctx: ctx, db: db}
}

func (r *Data) Master() *sql.DB {
	return r.db
}

func MustConnectPostgres(ctx context.Context, conf *config.Config) *sql.DB {
	connectionStr := fmt.Sprintf("postgres://postgres:mysecretpassword@%s:%s/postgres?sslmode=disable&application_name=quotation&connect_timeout=5",
		conf.Postgres.Host, conf.Postgres.Port)

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(); err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		err = connection.Close()
		if err != nil {
			logger.Errf("can't close database connection, err: %s", err)
			return
		}
	}()

	logger.Infof("connected to db: %+v", connection.Stats())
	return connection
}

func MustMigrate(connection *sql.DB) {
	driver, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrationPath := fmt.Sprintf("file://%s/migration", path)
	fmt.Printf("migrationPath : %s\n", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Infof("no changes in migration, skip")

		} else {
			panic(err)
		}
	}
}
