package app

import (
	"net/http"
)


// InitializeRoutes setup route mappings
func InitializeRoutes(env *Env) {

	var dir string

	uiRouter := env.Router.PathPrefix("/").Subrouter()
	apiRouterAuth := env.Router.PathPrefix("/api").Subrouter()
	apiRouterToken := env.Router.PathPrefix("/api").Subrouter()

	//html
	uiRouter.HandleFunc("/", env.index)
	uiRouter.HandleFunc("/favicon.ico", favicon)
	uiRouter.HandleFunc("/internal", env.internal)
	uiRouter.HandleFunc("/login", env.index).Methods("GET")
	uiRouter.HandleFunc("/login", env.login).Methods("POST")
	uiRouter.HandleFunc("/logout", env.logout).Methods("GET")
	uiRouter.HandleFunc("/signup", env.signupIndex).Methods("GET")
	uiRouter.HandleFunc("/signup", env.signup).Methods("POST")
	uiRouter.HandleFunc("/result", env.token).Methods("GET")

	//api
	apiRouterAuth.HandleFunc("/token", env.generateToken).Methods("GET")
	apiRouterToken.HandleFunc("/events", getAllEvents).Methods("GET")
	apiRouterToken.HandleFunc("/events", createEvent).Methods("POST")
	apiRouterToken.HandleFunc("/events", updateEvent).Methods("PATCH")
	apiRouterToken.HandleFunc("/events/{id}", getEvent).Methods("GET")
	apiRouterToken.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	apiRouterToken.Use(AuthMiddleware)

	// Server static files
	uiRouter.PathPrefix("/website/static").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	uiRouter.PathPrefix("/swaggerui").Handler(http.StripPrefix("/", http.FileServer(http.Dir("website"))))
}
