// Package classification simple crud app API.
//
// Swagger API for simple crud app
//
// Terms Of Service: http://swagger.io/terms/
//
//     Schemes: http, https
//	   BasePath: /api
//     Version: 0.0.1
//	   License: MIT http://opensource.org/licenses/MIT
//     Contact: John Bishop<bishy999@hotmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - Bearer: []
//
//     SecurityDefinitions:
//     Bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//          description: |
//           For accessing the API a valid JWT token must be passed in all the queries in
//           the 'Authorization' header.
//
//
// swagger:meta
package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/bishy999/go-simple-webapp/pkg/app"
)

func main() {

	db := app.InitDB()
	tpl := template.Must(template.ParseGlob("/website/templates/*"))
	router := mux.NewRouter().StrictSlash(true)

	env := &app.Env{DB: db, Tpl: tpl, Router: router, DbSessionsCleaned: time.Now()}
	app.InitializeRoutes(env)

	log.Println("#########################")
	log.Println("Server is up on 8080 port")
	log.Println("#########################")

	log.Fatalln(http.ListenAndServe(":8080", router))

}
