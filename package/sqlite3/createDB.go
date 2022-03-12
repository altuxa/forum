package sqlite3

import "database/sql"

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE Users (
		"id"	INTEGER NOT NULL UNIQUE,
		"login"	TEXT UNIQUE,
		"email"	TEXT UNIQUE,
		"password"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE "Sessions" (
		"UserId"	INTEGER NOT NULL UNIQUE,
		"Value"	TEXT,
		"Expires"	DATETIME NOT NULL,
		FOREIGN KEY("UserId") REFERENCES "Users"("id")
	)`)

	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE "Posts" (
		"Id"	INTEGER NOT NULL,
		"UserId"	INTEGER NOT NULL,
		"Title"	TEXT NOT NULL,
		"Text"	TEXT NOT NULL,
		"Author"	TEXT,
		"Likes"	INTEGER,
		"Dislikes"	INTEGER,
		PRIMARY KEY("Id" AUTOINCREMENT),
		FOREIGN KEY("UserId") REFERENCES "Users"("id")
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE "Comments" (
		"Id"	INTEGER NOT NULL UNIQUE,
		"UserId"	INTEGER,
		"PostId"	INTEGER,
		"Text"	TEXT NOT NULL,
		"Author"	TEXT,
		"Likes"	INTEGER,
		"Dislikes"	INTEGER,
		FOREIGN KEY("UserId") REFERENCES "Users"("id"),
		FOREIGN KEY("PostId") REFERENCES "Posts"("Id"),
		PRIMARY KEY("Id" AUTOINCREMENT)
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE "Tags" (
		"Id"	INTEGER NOT NULL UNIQUE,
		"PostId"	INTEGER,
		"Value"	TEXT NOT NULL,
		FOREIGN KEY("PostId") REFERENCES "Posts"("Id"),
		PRIMARY KEY("Id" AUTOINCREMENT)
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE "RatingPost" (
		"PostId"	INTEGER,
		"UserId"	INTEGER,
		"Like"	INTEGER,
		"Dislike"	INTEGER,
		FOREIGN KEY("UserId") REFERENCES "Users"("id"),
		FOREIGN KEY("PostId") REFERENCES "Posts"("Id")
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE "RatingComment" (
		"CommentId"	INTEGER,
		"UserId"	INTEGER,
		"Like"	INTEGER,
		"Dislike"	INTEGER,
		FOREIGN KEY("CommentId") REFERENCES "Comments"("Id"),
		FOREIGN KEY("UserId") REFERENCES "Users"("id")
	)`)
	if err != nil {
		return err
	}
	return nil
}
