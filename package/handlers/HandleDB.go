package handlers

import "database/sql"

//Handle ..
type Handle struct {
	DB *sql.DB
}

//CreateHandle ..
func CreateDB(db *sql.DB) *Handle {
	return &Handle{
		DB: db,
	}
}
