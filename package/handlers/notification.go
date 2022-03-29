package handlers

import (
	"forum/package/models"
	"forum/package/sqlite3"
	"net/http"
	"strconv"
	"text/template"
)

func (db *Handle) Notification(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/notification.html")
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if userId == 0 {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	Comments := []models.Comments{}
	NotifComments, err := sqlite3.GetNotificationComments(db.DB, userId)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	for _, i := range NotifComments {
		comment := sqlite3.GetOneCommentfromDB(db.DB, i.CommentId, i.UserId)
		Comments = append(Comments, comment)
	}
	RateNotif, err := sqlite3.GetRateNotification(db.DB, userId)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	for j, i := range RateNotif {
		login := sqlite3.GetUserLogin(db.DB, i.UserId)
		RateNotif[j].RateAuthor = login
	}
	type Parse struct {
		Comments []models.Comments
		Rate     []models.RateNotification
	}
	var MyNotification Parse
	MyNotification.Comments = Comments
	MyNotification.Rate = RateNotif
	err = ts.Execute(w, MyNotification)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
}

func (db *Handle) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if userId == 0 {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	but := r.FormValue("ReadButton")
	// fmt.Println(but)
	// all := r.FormValue("readAll")
	// fmt.Println(all)
	if len(but) != 0 {
		comId, err := strconv.Atoi(but)
		if err != nil {
			CustomError(http.StatusBadRequest, w)
			return
		}
		err = sqlite3.DeleteOneNotificationComment(db.DB, comId, userId)
		if err != nil {
			CustomError(http.StatusBadRequest, w)
			return
		}
		http.Redirect(w, r, "/notification", 302)
	} else if len(but) == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
}

func (db *Handle) ReadRateNotif(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if userId == 0 {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	rate := r.FormValue("ReadRate")
	if len(rate) == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	notifId, err := strconv.Atoi(rate)
	if err != nil {
		CustomError(http.StatusBadRequest, w)
		return
	}
	err = sqlite3.DeleteRateNotification(db.DB, userId, notifId)
	if err != nil {
		CustomError(http.StatusBadRequest, w)
		return
	}
	http.Redirect(w, r, "/notification", 302)
}

func (db *Handle) ReadAllNotif(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if userId == 0 {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	r.FormValue("readAll")
	// fmt.Println(all, 1)
	sqlite3.DeleteAllNotification(db.DB, userId)
	http.Redirect(w, r, "/notification", 302)
}
