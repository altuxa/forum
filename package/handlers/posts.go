package handlers

import (
	"forum/package/models"
	"forum/package/sqlite3"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (db *Handle) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/post.html")
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	str := strings.TrimPrefix(r.URL.Path, "/post/")
	numb, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	POST, err := sqlite3.GetOnePost(db.DB, numb)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	if POST.Id == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	delet := r.FormValue("delete")
	change := r.FormValue("change")
	cookie, _ := r.Cookie("session")
	if cookie != nil {
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
	}
	if delet != " " && cookie != nil {
		gR := sqlite3.GetSessionsFromDB(db.DB, cookie)
		for j := 0; j < len(POST.Comment); j++ {
			if POST.Comment[j].UserId == gR {
				POST.Comment[j].AuthorOfCom = true
			} else {
				POST.Comment[j].AuthorOfCom = false
			}
		}
	}
	nu, _ := strconv.Atoi(change)
	for j, i := range POST.Comment {
		if i.Id == nu {
			i.Change = true
			POST.Comment[j].Change = true
		}
	}
	if cookie != nil {
		POST.IsOnline = true
	} else {
		POST.IsOnline = false
	}
	err = ts.Execute(w, POST)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
}

func (db *Handle) CreatePostMethodGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/createPosts.html")
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	id := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if cookie != nil {
		// userId := sqlite3.GetSessionsFromDB(db.DB, cookie)
		// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
		tok, _ := sqlite3.IsSessionExist(db.DB, id)
		if tok != cookie.Value {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
			return
		}
		/////////////////////////////
	}
	err = ts.Execute(w, parse)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	parse.CheckErr = false
}

func (db *Handle) CreatePostMethodPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	id := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if cookie != nil {
		// userId := sqlite3.GetSessionsFromDB(db.DB, cookie)
		// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
		tok, _ := sqlite3.IsSessionExist(db.DB, id)
		if tok != cookie.Value {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
			return
		}
		/////////////////////////////
	}
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	login := sqlite3.GetUserLogin(db.DB, id)
	text := r.FormValue("post")
	title := r.FormValue("title")
	data := r.PostForm["tags"]
	err = CheckPost(title, text, data)
	if err != nil {
		parse.ErrorMessage = err
		parse.CheckErr = true
		http.Redirect(w, r, "/createpost", 302)
		return
	}
	posts := models.Posts{
		UserId: id,
		Title:  title,
		Text:   text,
		Author: login,
		Tags:   data,
	}
	err = sqlite3.AddPostToDB(db.DB, posts)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func (db *Handle) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	id := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if cookie != nil {
		// userId := sqlite3.GetSessionsFromDB(db.DB, cookie)
		// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
		tok, _ := sqlite3.IsSessionExist(db.DB, id)
		if tok != cookie.Value {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
			return
		}
		/////////////////////////////
	}
	idstr := r.FormValue("deletebutton")
	idPost, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	Posts, err := sqlite3.GetPostsByUserId(db.DB, id)
	if err != nil {
		CustomError(http.StatusBadRequest, w)
		return
	}
	check := false
	for _, i := range Posts {
		if idPost == i.Id {
			check = true
		}
	}
	if !check {
		CustomError(http.StatusBadRequest, w)
		return
	}
	comments, err := sqlite3.SelectCommentsFromDB(db.DB, idPost)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	for _, i := range comments {
		err = sqlite3.DeleteRatingCommentByCommentId(db.DB, i.Id)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
	}
	err = sqlite3.DeleteAllCommentFromPost(db.DB, idPost)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = sqlite3.DeleteTags(db.DB, idPost)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = sqlite3.DeleteRatingByPostId(db.DB, idPost)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	sqlite3.DeleteNotificationIfPostDeleted(db.DB, idPost)
	err = sqlite3.DeletePost(db.DB, idPost)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func (db *Handle) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if cookie != nil {
		// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
		tok, _ := sqlite3.IsSessionExist(db.DB, userId)
		if tok != cookie.Value {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
			return
		}
		/////////////////////////////
	}
	str := strings.TrimPrefix(r.URL.Path, "/updatepost/")
	title := r.FormValue("newtitle")
	text := r.FormValue("newtext")
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	data := r.PostForm["tags"]
	err = CheckPost(title, text, data)
	if err != nil {
		CustomError1(err, http.StatusBadRequest, w)
		return
	}
	idPost, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	post, err := sqlite3.GetOnePost(db.DB, idPost)
	if err != nil {
		CustomError(http.StatusBadRequest, w)
		return
	}
	if post.UserId != userId {
		CustomError(http.StatusBadRequest, w)
		return
	}
	err = sqlite3.DeleteTags(db.DB, idPost)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = sqlite3.UpdatePost(db.DB, title, text, idPost, userId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = sqlite3.AddTagsToDB(db.DB, data, int64(idPost))
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func (db *Handle) RatingPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	str := strings.TrimPrefix(r.URL.Path, "/ratingpost/")
	postId, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	var numbLike, numbDislike, likeForPost, dislikeForPost int
	like := r.FormValue("likepost")
	dislike := r.FormValue("dislikepost")
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	if cookie != nil {
		// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
		tok, _ := sqlite3.IsSessionExist(db.DB, userId)
		if tok != cookie.Value {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
			return
		}
		/////////////////////////////
	}
	check, ratePost := sqlite3.CheckRatingPostTable(db.DB, userId, postId)
	if like == "like" && check == false {
		numbLike = numbLike + 1
		likeForPost = 1
	} else if dislike == "dislike" && check == false {
		numbDislike = numbDislike + 1
		dislikeForPost = 1
	} else if like == "like" && check == true && ratePost.Like == 1 {
		numbLike = 0
		numbDislike = 0
		likeForPost = -1
	} else if dislike == "dislike" && check == true && ratePost.Dislike == 1 {
		numbDislike = 0
		numbLike = 0
		dislikeForPost = -1
	} else if like == "like" && check == true && ratePost.Like == 0 && ratePost.Dislike == 1 {
		numbLike = 1
		numbDislike = 0
		likeForPost = 1
		dislikeForPost = -1
	} else if dislike == "dislike" && check == true && ratePost.Dislike == 0 && ratePost.Like == 1 {
		numbDislike = 1
		numbLike = 0
		dislikeForPost = 1
		likeForPost = -1
	} else if like == "like" && check == true && ratePost.Like == 0 && ratePost.Dislike == 0 {
		numbLike = 1
		likeForPost = 1
	} else if dislike == "dislike" && check == true && ratePost.Like == 0 && ratePost.Dislike == 0 {
		numbDislike = 1
		dislikeForPost = 1
	} else {
		CustomError(http.StatusBadRequest, w)
		return
	}
	if check == false {
		err = sqlite3.AddLikeOrDislikeToRatingPost(db.DB, userId, postId, numbLike, numbDislike)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
	} else if check == true {
		err = sqlite3.UpdateLikeOrDislikeToRatingPost(db.DB, userId, postId, numbLike, numbDislike)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
	}
	err = sqlite3.AddRatingToPost(db.DB, likeForPost, dislikeForPost, postId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	post, err := sqlite3.GetOnePost(db.DB, postId)
	if err != nil {
		CustomError(http.StatusBadRequest, w)
		return
	}
	if userId != post.UserId {
		rate := CheckRate(likeForPost, dislikeForPost)
		err = sqlite3.AddRateNotificationToDB(db.DB, postId, userId, rate)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
	}

	http.Redirect(w, r, "/post/"+str, 302)
}
