package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Source struct {
	DatabaseDsn string
	db          *sql.DB
}

func NewDB(DatabaseDsn string) *Source {
	db, err := sql.Open("pgx", DatabaseDsn)
	if err != nil {
		log.Fatal("error access into DB")
	}

	return &Source{DatabaseDsn: DatabaseDsn, db: db}
}

func (dbSource *Source) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := dbSource.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (dbSource *Source) InitTables() {
	_, err := dbSource.db.Exec("create table data (hash varchar(20) not null constraint data_pk primary key, url varchar(500))")
	if err != nil {
		log.Println("init tables are NOT created - ", err)
		return
	}
	log.Println("init tables are created")
}

func (dbSource *Source) Save(hash string, url string) {
	log.Println("try to save; hash=", hash, " url=", url)

	row, err := dbSource.db.Exec("insert into data (hash, url) values ($1, $2)", hash, url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("saved", row)
}

func (dbSource *Source) Get(hash string) string {
	var url string
	row := dbSource.db.QueryRow("select url from data where hash = $1", hash)
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ""
		}
		log.Fatal(err)
	}
	return url
}

func (dbSource *Source) GetAll() map[string]string {
	var hash string
	var url string
	var data = make(map[string]string)

	rows, err := dbSource.db.Query("select hash, url from data")
	if err != nil {
		log.Println(err)
		return data
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&hash, &url)
		if err != nil {
			log.Println(err)
			return data
		}
		data[hash] = url
	}
	err = rows.Err()
	if err != nil {
		return data
	}
	return data
}
