package service

import (
	"github.com/gorilla/mux"
	"net/http"
)

// init router pointer
func NewRouter() *mux.Router  {

	//  init Gorilla Router
	// Enable StrictSlash allows matching both /path/ and /path
	router := mux.NewRouter().StrictSlash(true)

	// init routes in routes.go
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}

	return router
}

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// init system routes
var routes = Routes{
	Route{
		"GetLeader",
		"GET",
		"/leader/{leaderId}",
		FindLeaderByID,
	},
	Route {
		"HealthCheck",
		"GET",
		"/health",
		HealthCheck,
	},
	Route {
		"SetDBHealth",
		"GET",
		"/init/db-health/{state}",
		SetDBHealthState,
	},
}




