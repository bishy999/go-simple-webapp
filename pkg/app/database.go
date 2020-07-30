package app

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//InitDB setup database for use
func InitDB() *sql.DB {

	var c conf
	c.getConf()

	db, err := sql.Open("mysql", c.Usename+":"+c.Password+"@tcp("+c.Host+":"+c.Port+")/"+c.Name+"?charset=utf8&parseTime=true")
	check(err)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	check(err)

	return db

}

/**
func (env *Env) createDB(w http.ResponseWriter, req *http.Request) {

	stmt, err := env.DB.Prepare(`CREATE TABLE user (userid VARCHAR(30), passwd CHAR(60));`)
	//stmt, err := db.Prepare(`CREATE TABLE session (uid CHAR(60), userid VARCHAR(30), lastActivity DATETIME);`)
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "CREATED TABLE user", n)
}
**/

// findUserSession searches for the user associated with this session
// returns session if found
func (env *Env) findSession(id string) session {

	rows, err := env.DB.Query("SELECT * FROM session WHERE uid=?", id)
	check(err)
	defer rows.Close()

	s := session{}
	for rows.Next() {
		var uid string
		var userid string
		var lastActivity time.Time
		err = rows.Scan(&uid, &userid, &lastActivity)
		check(err)
		s.id = uid
		s.un = userid
		s.lastActivity = lastActivity
	}

	return s
}

// findUser searches for the user associated with the username provided and
// returns user if found
func (env *Env) findUser(username string) user {

	rows, err := env.DB.Query("SELECT * FROM user WHERE userid=?", username)
	check(err)
	defer rows.Close()

	u := user{}
	for rows.Next() {
		var id string
		var passwd []byte
		err = rows.Scan(&id, &passwd)
		check(err)
		u.UserName = id
		u.Password = passwd
	}

	return u
}

// addUser add user to persistence layer
// returns true if added successfully
func (env *Env) addUser(u user) bool {

	ustmt, err := env.DB.Prepare("INSERT INTO user(userid, passwd) VALUES(?,?)")
	check(err)
	defer ustmt.Close()

	r, err := ustmt.Exec(u.UserName, u.Password)
	check(err)

	n, err := r.RowsAffected()
	check(err)
	log.Printf("Added %d record in user table for user %s", n, u.UserName)

	return true

}

// addSession add session to persistence layer
// returns true if added successfully
func (env *Env) addSession(s session) bool {

	stmt, err := env.DB.Prepare("INSERT INTO session(uid, userid, lastActivity) VALUES(?,?,?)")
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(s.id, s.un, s.lastActivity)
	check(err)

	n, err := r.RowsAffected()
	check(err)
	log.Printf("Added %d record in session table for user %s with a session of %s", n, s.un, s.id)

	return true

}

// deleteSession delete session from persistence layer
// returns true if added successfully
func (env *Env) deleteSession(s string) bool {

	stmt, err := env.DB.Prepare("DELETE FROM session WHERE uid=?")
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(s)
	check(err)

	n, err := r.RowsAffected()
	check(err)
	log.Printf("Deleted %v records in session table for session %s", n, s)

	return true
}

// findUserSession searches for the user associated with this session
// returns session if found
func (env *Env) findAllSession() []session {

	rows, err := env.DB.Query("SELECT * FROM session")
	check(err)
	defer rows.Close()

	var s []session
	for rows.Next() {
		var uid string
		var userid string
		var lastActivity time.Time
		err = rows.Scan(&uid, &userid, &lastActivity)
		check(err)
		s = append(s, session{id: uid, un: userid, lastActivity: lastActivity})
	}

	return s

}
