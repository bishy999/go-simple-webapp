package app

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// AppKey ghg
	AppKey = "golangcode.com"
)

// Env contains data required to start
type Env struct {
	DB                *sql.DB
	Tpl               *template.Template
	Router            *mux.Router
	DbSessionsCleaned time.Time
}

type user struct {
	UserName string
	Password []byte
	Token    string
}

type session struct {
	id           string
	un           string
	lastActivity time.Time
}

func favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "website/static/favicon.ico")
}

func (env *Env) index(w http.ResponseWriter, req *http.Request) {
	env.Tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func (env *Env) signupIndex(w http.ResponseWriter, req *http.Request) {
	env.Tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func (env *Env) internal(w http.ResponseWriter, req *http.Request) {
	u := env.getUser(w, req)
	env.Tpl.ExecuteTemplate(w, "internal.gohtml", u)
}

func (env *Env) login(w http.ResponseWriter, req *http.Request) {

	if env.loggedIn(w, req) {
		http.Redirect(w, req, "/internal", http.StatusSeeOther)
		return
	}

	log.Println(" ### Login ### ")
	username := req.FormValue("email")
	passwd := req.FormValue("password")

	u, err := env.checkCredentials(username, passwd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}
	log.Printf("User %v successfully logged in", u.UserName)

	// create session
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	s := session{c.Value, username, time.Now()}
	ok := env.addSession(s)
	if !ok {
		http.Error(w, "Unable to signup successfully!", http.StatusForbidden)
		return
	}

	// redirect
	http.Redirect(w, req, "/internal", http.StatusSeeOther)
	return

}

func (env *Env) signup(w http.ResponseWriter, req *http.Request) {

	if env.loggedIn(w, req) {
		http.Redirect(w, req, "/internal", http.StatusSeeOther)
		return
	}

	log.Println(" ### Signup ### ")
	username := req.FormValue("email")
	passwd := req.FormValue("password")

	// check if username exists
	u := env.findUser(username)
	if u.UserName != "" {
		http.Error(w, "Username already taken", http.StatusForbidden)
		return
	}

	// create session
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	s := session{c.Value, username, time.Now()}
	bs, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	u = user{username, bs, ""}
	userok := env.addUser(u)
	if !userok {
		http.Error(w, "Unable to signup successfully!", http.StatusForbidden)
		return
	}
	sessok := env.addSession(s)
	if !sessok {
		http.Error(w, "Unable to signup successfully!", http.StatusForbidden)
		return
	}

	//redirect
	http.Redirect(w, req, "/internal", http.StatusSeeOther)
	return

}

func (env *Env) logout(w http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	log.Println(" ### Logout ### ")

	// delete the session
	ok := env.deleteSession(c.Value)
	if !ok {
		http.Error(w, "Unable to logout successfully!", http.StatusForbidden)
		return
	}

	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up persisted sessions
	fmt.Println(time.Now().Sub(env.DbSessionsCleaned))
	if time.Now().Sub(env.DbSessionsCleaned) > (time.Second * 60) {
		go env.cleanSessions()
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (env *Env) token(w http.ResponseWriter, req *http.Request) {

	if !env.loggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	log.Println(" ### Generate Token ### ")
	u := env.getUser(w, req)
	log.Printf(" Generating token for user: %s ", u.UserName)
	tkn, err := createToken(u.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}

	u.Token = tkn
	env.Tpl.ExecuteTemplate(w, "internal.gohtml", u)
}

// AuthMiddleware  is middleware to check token is valid. Returning
// a 401 status to the client if it is not valid.
func AuthMiddleware(next http.Handler) http.Handler {
	if len(AppKey) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(AppKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}

func (env *Env) checkCredentials(username, passwd string) (user, error) {

	// check if username exists
	u := env.findUser(username)
	if u.UserName == "" {
		return u, errors.New("Pleases authenticate with valid credentials")
	}

	// does the entered password match the stored password?
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwd))
	if err != nil {
		return u, errors.New("Username and/or password do not match")
	}
	return u, nil
}
