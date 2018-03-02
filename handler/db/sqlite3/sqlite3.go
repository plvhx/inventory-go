package sqlite3

import(
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitHandler(f string) (db *sql.DB, e error) {
	db, err := sql.Open("sqlite3", f)

	return db, err
}

func DestroyHandler(db *sql.DB) {
	db.Close()
}