package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"text/template"

	"forum/package/models"
	"forum/package/sqlite3"
)

func (db *Handle) RegistrGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/registration.html")
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	// можно не делать
	if r.URL.Path != "/registr" {
		CustomError(http.StatusNotFound, w)
		return
	}
	cookie, _ := r.Cookie("session")
	if cookie != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Referer() == "http://localhost:8080/registr" {
		err = ts.Execute(w, parse)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		parse.CheckErr = false
		return
	}
	err = ts.Execute(w, parse)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
}

func (db *Handle) RegistrPOST(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if cookie != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	user := models.Users{
		Login:          r.FormValue("login"),
		Email:          r.FormValue("email"),
		Password:       r.FormValue("password"),
		AcceptPassword: r.FormValue("acceptpass"),
	}
	///////////////////
	if user.Login != "" {
		parse.CheckErr = false
		err1 := CheckLogin(user.Login)
		if err1 != nil {
			parse.ErrorMessage = err1
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		if len(user.Login) < 7 {
			parse.ErrorMessage = errors.New("login length is less than 7 characters")
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		} else if len(user.Login) > 20 {
			parse.ErrorMessage = errors.New("login length is more than 20 characters")
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		_, err := mail.ParseAddress(user.Email)
		if err != nil {
			parse.ErrorMessage = err
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		err2 := CheckPass(user.Password)
		if err2 != nil {
			parse.ErrorMessage = err2
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		if len(user.Password) < 12 {
			parse.ErrorMessage = errors.New("password length is less than 12 characters")
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		if len(user.Password) > 25 {
			parse.ErrorMessage = errors.New("password length is more than 25 characters")
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		sqlite3.CreateTables(db.DB)
		err = sqlite3.InsertNewUserIntoDB(db.DB, user)
		if err != nil {
			parse.ErrorMessage = err
			parse.CheckErr = true
			http.Redirect(w, r, "/registr", 302)
			return
		}
		if !parse.CheckErr {
			parse.Cng = "Congratulations you have registered"
			http.Redirect(w, r, "/registr", 302)
		}
	}
}

func (db *Handle) LoginGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/login.html")
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	cookie, _ := r.Cookie("session")
	if cookie != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Referer() == "http://localhost:8080/login" {
		err = ts.Execute(w, parse)
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		parse.CheckErr = false
		return
	}
	err = ts.Execute(w, parse)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	parse.CheckErr = false
}

func (db *Handle) LoginPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, _ := r.Cookie("session")
	if cookie != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	user := models.Users{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}
	err2 := CheckEmptyLogin(user.Login, user.Password)
	if err2 != nil {
		parse.ErrorMessage = err2
		parse.CheckErr = true
		http.Redirect(w, r, "login", 302)
		return
	}
	if user.Login != "" {
		check, err1 := sqlite3.SelectUsers(db.DB, user.Login, user.Password)
		if err1 != nil {
			parse.ErrorMessage = errors.New("Incorrect login or password")
			parse.CheckErr = true
			http.Redirect(w, r, "/login", 302)
			return
		}
		if check != 0 {
			a := sqlite3.MakeCookie()
			fmt.Println(a)
			err := sqlite3.InsertCookieIntoDB(db.DB, user.Login, a)
			if err != nil {
				CustomError(http.StatusInternalServerError, w)
				return
			}
			http.SetCookie(w, a)
			http.Redirect(w, r, "/", 302)
			// return
		}
	}
}

func (db *Handle) LogoutGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/logout.html")
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
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
	err = ts.Execute(w, nil)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
}

func (db *Handle) LogoutPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
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
	choice := r.FormValue("choice")
	if choice == "yes" {
		sqlite3.DeleteCookieFromDB(db.DB, cookie)
		log.Println("Delele cookie from DB happened successfully ")
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 302)
	} else if choice == "no" {
		http.Redirect(w, r, "/", 302)
		return
	}
}
