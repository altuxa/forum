package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/package/models"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./forum.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// SelectUsers.........
func SelectUsers(db *sql.DB, login, password string) (int, error) {
	users := models.Users{}
	row := db.QueryRow("SELECT Id, Password FROM Users WHERE Login = ?", login)
	err := row.Scan(&users.Id, &users.Password)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))
	if err != nil {
		return 0, err
	}
	return users.Id, nil
}

func InsertNewUserIntoDB(db *sql.DB, user models.Users) error {
	if user.Password != user.AcceptPassword {
		return errors.New("Passwords do not match")
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO Users(login, email , password) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Login, user.Email, hashPass)
	if err != nil {
		return fmt.Errorf("Sorry login or email is already taken")
	}
	return nil
}

// test function for add cookie
func InsertCookieIntoDB(db *sql.DB, login string, cookie *http.Cookie) error {
	var id int
	row := db.QueryRow("SELECT Id FROM Users WHERE Login =?;", login)
	row.Scan(&id)

	_, err := db.Exec("DELETE FROM Sessions WHERE UserId =?", id)
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := db.Prepare("INSERT INTO Sessions (UserId, Value,Expires)VALUES(?,?,?)")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(id, cookie.Value, cookie.Expires)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteCookieFromDB(db *sql.DB, cookie *http.Cookie) error {
	_, err := db.Exec("DELETE FROM Sessions WHERE Value = ?", cookie.Value)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetSessionsFromDB(db *sql.DB, cookie *http.Cookie) int {
	var UserId int
	row := db.QueryRow("SELECT UserId FROM Sessions WHERE Value = ?", cookie.Value)
	row.Scan(&UserId)
	return UserId
}

func GetUserLogin(db *sql.DB, id int) string {
	var login string
	row := db.QueryRow("SELECT login FROM Users WHERE Id = ?", id)
	row.Scan(&login)
	return login
}
