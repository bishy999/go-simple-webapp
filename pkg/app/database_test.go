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
	env := &Env{DB: db}

	rows := sqlmock.NewRows([]string{"uid", "userid", "lastActivity"}).
		AddRow(uid, "test999@test.com", time.Now())
	mock.ExpectQuery("SELECT").WithArgs(uid).WillReturnRows(rows)

	sess := env.findSession(uid)
	assert.NotEmpty(t, sess)
}
