package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (handle *Handle) Server() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", handle.Home)
	mux.HandleFunc("/registr", handle.RegistrGET)
	mux.HandleFunc("/registrconfirm", handle.RegistrPOST)
	mux.HandleFunc("/login", handle.LoginGET)
	mux.HandleFunc("/loginconfirm", handle.LoginPOST)
	mux.HandleFunc("/logout", handle.LogoutGET)
	mux.HandleFunc("/logoutconfirm", handle.LogoutPOST)

	mux.HandleFunc("/myprofile", handle.Profile)

	mux.HandleFunc("/createpost", handle.CreatePostMethodGet)
	mux.HandleFunc("/addpost", handle.CreatePostMethodPOST)

	mux.HandleFunc("/post/", handle.Post)
	mux.HandleFunc("/deletepost/", handle.DeletePost)
	mux.HandleFunc("/updatepost/", handle.UpdatePost)
	mux.HandleFunc("/ratingpost/", handle.RatingPost)

	mux.HandleFunc("/tagsfilter/", handle.TagsFilter)

	mux.HandleFunc("/createcomment/", handle.CreateComment)
	mux.HandleFunc("/deletecomment/", handle.DeleteComment)
	mux.HandleFunc("/updatecomment/", handle.UpdateComment)
	mux.HandleFunc("/ratingcomment/", handle.RatingComment)

	mux.HandleFunc("/notification", handle.Notification)
	mux.HandleFunc("/readnotification", handle.DeleteNotification)
	mux.HandleFunc("/readratenotification", handle.ReadRateNotif)
	mux.HandleFunc("/readallnotification", handle.ReadAllNotif)
	fmt.Println("starting server at localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
