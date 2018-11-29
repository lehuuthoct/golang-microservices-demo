package service

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"lehuuthoct/lht-microservices/leader/db"
)

// init variables
var BoltClient db.IBoltClient
const (
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE = "Content-Type"
	CONTENT_LENGTH = "Content-Length"
)
// start web server at specific port
func StartWebServer(port string) {
	logrus.Infof("Starting Web Server Port[%v] \n", port)

	// init routes
	r := NewRouter()
	http.Handle("/", r)

	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		logrus.Printf("Error starting server %v", err.Error())
	}

}
