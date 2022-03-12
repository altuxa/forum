package main

import (
	"log"

	handlers "forum/package/handlers"
	"forum/package/sqlite3"
)

func main() {
	db, err := sqlite3.OpenDB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	sqlite3.CreateTables(db)
	ch := make(chan int)
	go sqlite3.RemoveExpiredCookie(db, ch)
	sqlite3.CreateHandle(db)
	handle := handlers.CreateDB(db)
	handle.Server()
}
