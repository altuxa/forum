package handlers

import (
	"forum/package/models"
	"forum/package/sqlite3"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (db *Handle) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 303)
		return
	}
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
	str := strings.TrimPrefix(r.URL.Path, "/createcomment/")
	postId := r.FormValue("numberofpost")
	numb, err := strconv.Atoi(postId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	comtext := r.FormValue("comment")
	err = CheckComment(comtext)
	if err != nil {
		CustomError1(err, http.StatusBadRequest, w)
		// http.Redirect(w, r, "/post/"+str, 302)
		return
	}
	if comtext != "" {
		id := sqlite3.IsTokenExists(db.DB, cookie.Value)
		login := sqlite3.GetUserLogin(db.DB, id)
		com := models.Comments{
			UserId: id,
			PostId: numb,
			Text:   comtext,
			Author: login,
		}
		err, comId := sqlite3.AddCommentToDB(db.DB, com)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
		post, err := sqlite3.GetOnePost(db.DB, numb)
		if post.UserId != id {
			err = sqlite3.AddNotificationToDB(db.DB, com.PostId, com.UserId, comId)
			if err != nil {
				CustomError(http.StatusInternalServerError, w)
				return
			}
		}
	}
	http.Redirect(w, r, "/post/"+str, 302)
}

func (db *Handle) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	str := strings.TrimPrefix(r.URL.Path, "/deletecomment/")
	postId, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	delet := r.FormValue("delete")
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 303)
		return
	}
	comId, err := strconv.Atoi(delet)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	userId := sqlite3.GetSessionsFromDB(db.DB, cookie)
	comment := sqlite3.GetOneCommentfromDB(db.DB, comId, userId)
	if comment.Id == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	err = sqlite3.DeleteRatingCommentByCommentId(db.DB, comId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	sqlite3.DeleteNotificationIfCommentDeleted(db.DB, comId)
	err = sqlite3.DeleteCommentFromDB(db.DB, comId, userId)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), 302)
}

func (db *Handle) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	str := strings.TrimPrefix(r.URL.Path, "/updatecomment/")
	postId, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 303)
		return
	}
	update := r.FormValue("update")
	text := r.FormValue("changetext")
	// if text != "" {
	// 	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), 302)
	// 	return
	// }
	err = CheckComment(text)
	if err != nil {
		CustomError1(err, http.StatusBadRequest, w)
		// http.Redirect(w, r, "/post/"+strconv.Itoa(postId), 302)
		return
	}
	comId, err := strconv.Atoi(update)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusBadRequest, w)
		return
	}
	userId := sqlite3.GetSessionsFromDB(db.DB, cookie)
	com := sqlite3.GetOneCommentfromDB(db.DB, comId, userId)
	if com.Id == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	err = sqlite3.UpdateComment(db.DB, text, comId, userId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), 302)
}

func (db *Handle) RatingComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 303)
		return
	}
	user := sqlite3.GetSessionsFromDB(db.DB, cookie)
	// УДЛАЛЕНИЕ КУКОВ У ПОЛЬЗОВАТЕЛЯ ЕСЛИ ОН ЗАШЕЛ С ДВУХ БРАУЗЕРОВ
	tok, _ := sqlite3.IsSessionExist(db.DB, user)
	if tok != cookie.Value {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 303)
		return
	}
	/////////////////////////////
	str := strings.TrimPrefix(r.URL.Path, "/ratingcomment/")
	comId, err := strconv.Atoi(str)
	if err != nil {
		CustomError(http.StatusBadRequest, w)
		return
	}
	var numbLike, numbDislike, likeForComment, dislikeForComment int
	like := r.FormValue("likecomment")
	dislike := r.FormValue("dislikecomment")
	userId := sqlite3.IsTokenExists(db.DB, cookie.Value)
	check, rateComment := sqlite3.CheckRatingCommentTable(db.DB, userId, comId)
	if like == "like" && check == false {
		numbLike = numbLike + 1
		likeForComment = 1
	} else if dislike == "dislike" && check == false {
		numbDislike = numbDislike + 1
		dislikeForComment = 1
	} else if like == "like" && check == true && rateComment.Like == 1 {
		numbLike = 0
		numbDislike = 0
		likeForComment = -1
	} else if dislike == "dislike" && check == true && rateComment.Dislike == 1 {
		numbDislike = 0
		numbLike = 0
		dislikeForComment = -1
	} else if like == "like" && check == true && rateComment.Like == 0 && rateComment.Dislike == 1 {
		numbLike = 1
		numbDislike = 0
		likeForComment = 1
		dislikeForComment = -1
	} else if dislike == "dislike" && check == true && rateComment.Dislike == 0 && rateComment.Like == 1 {
		numbDislike = 1
		numbLike = 0
		dislikeForComment = 1
		likeForComment = -1
	} else if like == "like" && check == true && rateComment.Like == 0 && rateComment.Dislike == 0 {
		numbLike = 1
		likeForComment = 1
	} else if dislike == "dislike" && check == true && rateComment.Like == 0 && rateComment.Dislike == 0 {
		numbDislike = 1
		dislikeForComment = 1
	} else if dislike != "dislike" && len(like) != 0 || like != "like" && len(dislike) != 0 || len(dislike) == 0 && len(like) == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	if check == false {
		err = sqlite3.AddLikeOrDislikeToRatingComment(db.DB, userId, comId, numbLike, numbDislike)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
	} else if check == true {
		err = sqlite3.UpdateLikeOrDislikeToRatingComment(db.DB, userId, comId, numbLike, numbDislike)
		if err != nil {
			log.Println(err)
			CustomError(http.StatusInternalServerError, w)
			return
		}
	}
	err = sqlite3.AddRatingToComment(db.DB, likeForComment, dislikeForComment, comId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	postId, err := sqlite3.GetPostIdByCommentId(db.DB, comId)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), 302)
}
