package sqlite3

import "database/sql"

//Handle ..
type Handle struct {
	DB *sql.DB
}

//CreateHandle ..
func CreateHandle(db *sql.DB) *Handle {
	return &Handle{
		DB: db,
	}
}
