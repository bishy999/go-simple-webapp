package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

// swagger:operation POST /events events createEvent
// ---
// summary: Create a new event
// description: If event creation is a success, event will be returned with Created (201).
// parameters:
// - name: event
//   description: event to add to the list of events
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/event"
// responses:
//   "201":
//     "$ref": "#/responses/success"
//   "400":
//     "$ref": "#/responses/badRequest"
//   "404":
//     "$ref": "#/responses/notFound"
func createEvent(w http.ResponseWriter, req *http.Request) {

	data := event{}
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, " Error reading data. Please ensure your request is correct", http.StatusForbidden)
		return
	}

	if len(reqBody) == 0 {
		http.Error(w, "Enter data with the event title and description only in order to update", http.StatusForbidden)
		return
	}

	json.Unmarshal(reqBody, &data)
	events = append(events, data)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

}

// swagger:operation GET /events/{id} events getEvent
// ---
// summary: Return an event provided by the id.
// description: If the event is found, it will be returned else Error Not Found (404) will be returned.
// parameters:
// - name: id
//   in: path
//   description: id of the event
//   type: string
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/success"
//   "400":
//     "$ref": "#/responses/badRequest"
//   "404":
//     "$ref": "#/responses/notFound"
func getEvent(w http.ResponseWriter, req *http.Request) {

	var match bool

	w.Header().Add("Content-Type", "application/json")
	id := mux.Vars(req)["id"]
	for _, singleEvent := range events {
		if singleEvent.ID == id {
			json.NewEncoder(w).Encode(singleEvent)
			match = true
		}
	}
	if !match {
		log.Print("Requested event was not found for id: ", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}

}

// swagger:operation GET /events events getAllEvents
// ---
// summary: Return all events
// description: If the events are found, they will be returned else Error Not Found (404) will be returned.
// responses:
//   "200":
//     "$ref": "#/responses/success"
//   "400":
//     "$ref": "#/responses/badRequest"
//   "404":
//     "$ref": "#/responses/notFound"
func getAllEvents(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	if len(events) == 0 {
		log.Print("No events were found ")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}
	json.NewEncoder(w).Encode(events)
}

// swagger:operation DELETE /events/{id} events deleteEvent
// ---
// summary: Delete an event
// description: If the event is found, it will be deleted else Error Not Found (404) will be returned.
// parameters:
// - name: id
//   in: path
//   description: id of the event
//   type: string
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/success"
//   "400":
//     "$ref": "#/responses/badRequest"
//   "404":
//     "$ref": "#/responses/notFound"
func deleteEvent(w http.ResponseWriter, req *http.Request) {

	var match bool

	w.Header().Add("Content-Type", "application/json")
	id := mux.Vars(req)["id"]
	for i, singleEvent := range events {
		if singleEvent.ID == id {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", id)
			match = true
		}
	}
	if !match {
		log.Print("DELETE: Requested event was not found for id: ", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}
}

// swagger:operation PATCH /events events updateEvent
// ---
// summary: Update an event
// description: If event update is a success, event will be returned with Created (200).
// parameters:
// - name: event
//   description: update event details
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/event"
// responses:
//   "200":
//     "$ref": "#/responses/success"
//   "400":
//     "$ref": "#/responses/badRequest"
//   "403":
//     "$ref": "#/responses/forbidden"
//   "404":
//     "$ref": "#/responses/notFound"
func updateEvent(w http.ResponseWriter, req *http.Request) {

	var updatedEvent event
	var match bool

	w.Header().Add("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, " Error reading data. Please ensure your request is correct", http.StatusForbidden)
		return
	}
	if len(reqBody) == 0 {
		http.Error(w, "Enter data with the event title and description only in order to update", http.StatusForbidden)
		return
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == updatedEvent.ID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
			match = true
		}
	}
	if !match {
		log.Print("PATCH: Requested event was not found for id: ", updatedEvent.ID)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}
}

// swagger:operation GET /token token generateToken
// ---
// summary: Generate a token
// description: If the email and password are valid return token for further use.
// parameters:
// - name: email
//   description: email
//   in: query
//   required: true
// - name: password
//   description: password
//   in: query
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/success"
//   "400":
//     "$ref": "#/responses/badRequest"
//   "403":
//     "$ref": "#/responses/forbidden"
//   "404":
//     "$ref": "#/responses/notFound"
func (db *Env) generateToken(w http.ResponseWriter, req *http.Request) {

	log.Println(" ### Generate Token ### ")
	w.Header().Add("Content-Type", "application/json")

	username := req.FormValue("email")
	passwd := req.FormValue("password")

	// check if username exists
	u, err := db.checkCredentials(username, passwd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	log.Printf(" Generating token for user: %s ", username)

	tkn, err := createToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}

	u.Token = tkn
	io.WriteString(w, `{"token":"`+tkn+`"}`)

}

func createToken(username string) (string, error) {

	// We are happy with the credentials, so build a token. We've given it
	// an expiry of 1 hour.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(AppKey))
	if err != nil {
		return tokenString, errors.New("Pleases authenticate with valid credentials")
	}
	return tokenString, nil

}
