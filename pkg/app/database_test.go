package app

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindSession(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uid := uuid.New().String()
	userid := "test111@test.com"
	env := &Env{DB: db}

	rows := sqlmock.NewRows([]string{"uid", "userid", "lastActivity"}).
		AddRow(uid, userid, time.Now())
	mock.ExpectQuery("SELECT").WithArgs(uid).WillReturnRows(rows)

	sess := env.findSession(uid)
	if assert.NotNil(t, sess) {
		assert.Equal(t, uid, sess.id)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userid := "test112@test.com"
	passwd := "password1"
	env := &Env{DB: db}

	rows := sqlmock.NewRows([]string{"uid", "passwd"}).
		AddRow(userid, passwd)
	mock.ExpectQuery("SELECT").WithArgs(userid).WillReturnRows(rows)

	user := env.findUser(userid)
	if assert.NotNil(t, user) {
		assert.Equal(t, userid, user.UserName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	env := &Env{DB: db}
	u := user{UserName: "test113@test.com", Password: []byte("password1")}
	mock.ExpectPrepare("INSERT INTO user").ExpectExec().WithArgs(u.UserName, u.Password).WillReturnResult(sqlmock.NewResult(1, 1))

	user := env.addUser(u)
	assert.Equal(t, true, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestAddSession(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	env := &Env{DB: db}
	s := session{id: uuid.New().String(), un: "test114@test.com", lastActivity: time.Now()}
	mock.ExpectPrepare("INSERT INTO session").ExpectExec().WithArgs(s.id, s.un, s.lastActivity).WillReturnResult(sqlmock.NewResult(1, 1))

	user := env.addSession(s)
	assert.Equal(t, true, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestDeleteSession(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	env := &Env{DB: db}
	s := session{id: uuid.New().String()}
	mock.ExpectPrepare("DELETE FROM session").ExpectExec().WithArgs(s.id).WillReturnResult(sqlmock.NewResult(1, 1))

	user := env.deleteSession(s.id)
	assert.Equal(t, true, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestFindAllSession(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	env := &Env{DB: db}
	s := session{id: uuid.New().String(), un: "test115@test.com", lastActivity: time.Now()}

	rows := sqlmock.NewRows([]string{"uid", "userid", "lastActivity"}).
		AddRow(s.id, s.un, s.lastActivity)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sess := env.findAllSession()
	if assert.NotNil(t, sess) {
		for _, v := range sess {
			assert.Equal(t, v, s)
		}

	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
