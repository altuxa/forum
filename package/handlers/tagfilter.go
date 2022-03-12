package handlers

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"forum/package/models"
	"forum/package/sqlite3"
)

func (db *Handle) TagsFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	// if r.URL.Path != "/tagsfilter/" {
	// 	CustomError(http.StatusNotFound, w)
	// 	return
	// }
	ts, err := template.ParseFiles("./ui/html/tagsfilter.html")
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	str := strings.TrimPrefix(r.URL.Path, "/tagsfilter/")
	Tags, err := sqlite3.GetTagsByTagsName(db.DB, str)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
	Posts := []models.Posts{}
	for _, i := range Tags {
		Post, err := sqlite3.GetOnePost(db.DB, i.PostId)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		Posts = append(Posts, Post)
		// post, err := sqlite3.GetPostsByPostId(db.DB, i.PostId)
		// if err != nil {
		// 	log.Println(err)
		// 	CustomError(http.StatusInternalServerError, w)
		// 	return
		// }
		// for _, j := range post {
		// 	Posts = append(Posts, j)
		// }
	}
	//////////////////////////
	if len(Posts) == 0 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	err = ts.Execute(w, Posts)
	if err != nil {
		log.Println(err)
		CustomError(http.StatusInternalServerError, w)
		return
	}
}
