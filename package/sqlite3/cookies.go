package sqlite3

import (
	"database/sql"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func genToken() string {
	u1 := uuid.NewV4()
	return u1.String()
}

func RemoveExpiredCookie(db *sql.DB, ch chan int) {
	for {
		_, err := db.Exec("DELETE FROM Sessions WHERE Expires <?", time.Now())
		if err != nil {
			ch <- 0
		}
		time.Sleep(10 * time.Second)
	}
}

func MakeCookie() *http.Cookie {
	expiration := time.Now().Add(10 * time.Minute)
	cookie := http.Cookie{
		Name:     "session",
		Value:    genToken(),
		Expires:  expiration,
		HttpOnly: false,
		Path:     "/",
	}
	return &cookie
}

func DeleteCookie(token string) *http.Cookie {
	cookie := http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   -1,
	}
	return &cookie
}

func IsUserInSession(db *sql.DB, cookie *http.Cookie) bool {
	if cookie == nil {
		return false
	}
	var val string
	err := db.QueryRow("SELECT Value FROM Sessions WHERE Value = ?", cookie.Value).Scan(&val)
	if err == nil && err == sql.ErrNoRows {
		return false
	}
	if val == "" {
		return false
	}
	return true
}

func IsTokenExists(db *sql.DB, token string) int {
	var id int
	err := db.QueryRow("SELECT UserId FROM Sessions WHERE Value = ?", token).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func IsSessionExist(db *sql.DB, id int) (string, bool) {
	var token string
	row := db.QueryRow("SELECT Value FROM Sessions WHERE UserId = ?", id)
	row.Scan(&token)
	if token == "" {
		return "", false
	}
	return token, true
}
