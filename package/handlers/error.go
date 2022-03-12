package handlers

import (
	"net/http"
	"text/template"
)

func CustomError(code int, w http.ResponseWriter) {
	type Error struct {
		Code    int
		Message string
	}
	message := Error{Code: code, Message: http.StatusText(code)}
	gh, _ := template.ParseFiles("./ui/html/error.html")
	w.WriteHeader(code)
	gh.Execute(w, message)
}

func CustomError1(err error, code int, w http.ResponseWriter) {
	type Error struct {
		Error   error
		Code    int
		Message string
	}
	message := Error{Error: err, Code: code, Message: http.StatusText(code)}
	gh, _ := template.ParseFiles("./ui/html/error1.html")
	w.WriteHeader(code)
	gh.Execute(w, message)
}
