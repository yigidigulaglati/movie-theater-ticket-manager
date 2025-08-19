package dbInstance

import (
	"database/sql"
	"log/slog"
	"os"
	"ticket/api/queries"
)

type StoreDB struct{
	Queries *queries.Queries;
	DB *sql.DB;
}

var Store = StoreDB{}


func CreateDB(){

	dbCore, err := sql.Open(`sqlite`, `./core.sqlite3?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)`);
	if err != nil {
		slog.Error(`Could not open sqlite3 database: `+ err.Error());
		os.Exit(1);
	}

	Store.Queries = queries.New(dbCore);
	Store.DB = dbCore;
}