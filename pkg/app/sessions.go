package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

const sessionLength int = 300

func (db *Env) getUser(w http.ResponseWriter, req *http.Request) user {

	var u user
	var s session

	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	// if the user exists already, get user
	s = db.findSession(c.Value)
	if s.un == "" {
		return u
	}
	u = db.findUser(s.un)

	return u
}

// loggedIn return true if the user is already logged in
func (db *Env) loggedIn(w http.ResponseWriter, req *http.Request) bool {

	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	s := db.findSession(c.Value)
	if s.un == "" {
		return false
	}

	u := db.findUser(s.un)
	if u.UserName == "" {
		return false
	}

	// refresh session
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	return true
}

func (db *Env) cleanSessions() {
	log.Println(" ### Session Clean Up ### ")
	dbSessions := db.findAllSession()
	for _, v := range dbSessions {
		//if time.Since().Sub(v.lastActivity) > (time.Second * 30) {
		if time.Since(v.lastActivity) > (time.Second * 30) {
			db.deleteSession(v.id)
			fmt.Println(time.Since(v.lastActivity))
		}
	}
}
