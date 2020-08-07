package app

import (
	"testing"
)

func TestCreateTokenValid(t *testing.T) {

	username := "test"
	tkn, err := createToken(username)
	if err != nil {
		t.Fatal(err)
	}

	if tkn == "" {
		t.Errorf("token not correct for user %v got %v", username, tkn)
	}
}

func TestCreateTokenInvalid(t *testing.T) {

	username := ""
	expectedResponse := "provide valid username"
	tkn, err := createToken(username)
	if err != nil {
		if err.Error() != expectedResponse {
			t.Errorf("expected error message %v got %v", expectedResponse, err.Error())
		}
	}
	if tkn != "" {
		t.Errorf("token not correct for user %v got %v", username, tkn)
	}

}
