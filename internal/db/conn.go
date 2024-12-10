package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nurtai325/kaspi-service/internal/config"
)

var (
	pool            *sql.DB
	ErrNoConnection = errors.New("error initializing database connection pool")
)

const (
	driver = "postgres"
)

func connect(conf config.Config) error {
	dbUrl := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		driver,
		conf.POSTGRES_USER,
		conf.POSTGRES_PASSWORD,
		conf.POSTGRES_HOST,
		conf.POSTGRES_PORT,
		conf.POSTGRES_DB,
	)
	conn, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return err
	}
	err = conn.Ping()
	if err != nil {
		return err
	}
	pool = conn
	return nil
}

func Conn(conf config.Config) *sql.DB {
	if pool == nil {
		err := connect(conf)
		if err != nil {
			log.Panic(errors.Join(ErrNoConnection, err))
		}
	}
	return pool
}
