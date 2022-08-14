package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

type Source struct {
	DatabaseDsn string
}

func NewDB(DatabaseDsn string) *Source {
	return &Source{DatabaseDsn: DatabaseDsn}
}

func (dbSource *Source) Ping() error {
	if len(dbSource.DatabaseDsn) < 2 {
		return errors.New("DatabaseDsn too short")
	}
	db, err := sql.Open("pgx", dbSource.DatabaseDsn)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (dbSource *Source) InitTables() {
	db, err := sql.Open("pgx", dbSource.DatabaseDsn)
	if err != nil {
		log.Println("db connection error - init tables are NOT created")
		return
	}
	defer db.Close()

	_, err = db.Exec("create table data (hash varchar(20) not null constraint data_pk primary key, url varchar(500))")
	if err != nil {
		log.Println("init tables are NOT created - ", err)
		return
	}
	log.Println("init tables are created")
}
