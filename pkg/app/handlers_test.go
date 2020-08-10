package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestGetLogin(t *testing.T) {

	userid := "test112@test.com"
	passwd := "password1"
	u := user{UserName: userid, Password: []byte(passwd)}

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	form, _ := url.ParseQuery(req.URL.RawQuery)
	form.Add("email", userid)
	form.Add("password", passwd)
	req.URL.RawQuery = form.Encode()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	bs, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating hash for passwd", err)
	}

	rows := sqlmock.NewRows([]string{"uid", "passwd"}).AddRow(userid, bs)
	mock.ExpectQuery("SELECT").WithArgs(u.UserName).WillReturnRows(rows)
	mock.ExpectPrepare("INSERT INTO session").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	env := &Env{DB: db}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(env.login)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

func TestSignup(t *testing.T) {

	userid := "test113@test.com"
	passwd := "password1"
	u := user{UserName: userid, Password: []byte(passwd)}

	req, err := http.NewRequest("POST", "/signup", nil)
	if err != nil {
		t.Fatal(err)
	}

	form, _ := url.ParseQuery(req.URL.RawQuery)
	form.Add("email", userid)
	form.Add("password", passwd)
	req.URL.RawQuery = form.Encode()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"uid", "passwd"})
	mock.ExpectQuery("SELECT").WithArgs(u.UserName).WillReturnRows(rows)
	mock.ExpectPrepare("INSERT INTO user").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("INSERT INTO session").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	env := &Env{DB: db}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(env.signup)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}
