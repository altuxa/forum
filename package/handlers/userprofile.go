package handlers

import (
	"forum/package/models"
	"forum/package/sqlite3"
	"log"
	"net/http"
	"text/template"
)

func (db *Handle) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/profile.html")
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 303)
		return
	}
	com := []models.Comments{}
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
	FavoritePost := []models.Posts{}
	FavoritePostId, err := sqlite3.GetMyFavoritePosts(db.DB, userId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	for _, i := range FavoritePostId {
		pp, err := sqlite3.GetOnePost(db.DB, i.PostId)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
		FavoritePost = append(FavoritePost, pp)
	}
	Post, err := sqlite3.GetPostsByUserId(db.DB, userId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	for j, i := range Post {
		comments, err := sqlite3.SelectCommentsFromDB(db.DB, i.Id)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		Post[j].Comment = comments
		for _, t := range comments {
			com = append(com, t)
		}
	}
	type Home struct {
		Posts              []models.Posts
		Comment            []models.Comments
		LikePost           []models.Posts
		CheckMyPosts       bool
		CheckFavoritePosts bool
		CheckComments      bool
	}
	var gg Home
	seePost := r.FormValue("seemypost")
	if seePost == "see" {
		gg.CheckMyPosts = true
	}
	closePost := r.FormValue("closemypost")
	if closePost == "close" {
		gg.CheckMyPosts = false
	}
	seeFavorite := r.FormValue("seemyfavoriteposts")
	if seeFavorite == "see" {
		gg.CheckFavoritePosts = true
	}
	closeFavorite := r.FormValue("closemyfavoriteposts")
	if closeFavorite == "close" {
		gg.CheckFavoritePosts = false
	}
	seeComment := r.FormValue("seecomments")
	if seeComment == "see" {
		gg.CheckComments = true
	}
	closeComment := r.FormValue("closecomments")
	if closeComment == "close" {
		gg.CheckComments = false
	}
	gg.Posts = Post
	gg.Comment = com
	gg.LikePost = FavoritePost
	err = ts.Execute(w, gg)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
}
