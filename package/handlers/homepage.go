package handlers

import (
	"forum/package/models"
	"forum/package/sqlite3"
	"net/http"
	"strconv"
	"text/template"
)

type Index struct {
	Posts []models.Posts
	Check bool
	Notif int
}

func (db *Handle) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	if r.URL.Path != "/" {
		CustomError(http.StatusNotFound, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	var gg Index
	cookie, _ := r.Cookie("session")
	gg.Check = false
	if cookie != nil {
		gg.Check = true
	}
	gg.Posts, err = sqlite3.SelectPostsFromDB(db.DB)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	changeId := r.FormValue("change")
	chislo, _ := strconv.Atoi(changeId)
	for j := range gg.Posts {
		tags, err := sqlite3.GetAllTags(db.DB, gg.Posts[j].Id)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		gg.Posts[j].Tag = tags
	}
	if gg.Check {
		userId := sqlite3.GetSessionsFromDB(db.DB, cookie)
		// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
		tok, _ := sqlite3.IsSessionExist(db.DB, userId)
		if tok != cookie.Value {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
			return
		}
		/////////////////////////////
		for j, i := range gg.Posts {
			if i.UserId == userId {
				gg.Posts[j].IsAuthor = true
			}
			if i.Id == chislo {
				gg.Posts[j].ChangePost = true
			}

		}
		NotificationComments, err := sqlite3.GetNotificationComments(db.DB, userId)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		notifRate, err := sqlite3.GetRateNotification(db.DB, userId)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		gg.Notif = len(NotificationComments) + len(notifRate)
		err = ts.Execute(w, gg)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		return
	}
	err = ts.Execute(w, gg)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
}
