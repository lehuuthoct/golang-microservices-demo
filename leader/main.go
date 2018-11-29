package main

import (
	"github.com/sirupsen/logrus"
	"lehuuthoct/lht-microservices/leader/service"
	"lehuuthoct/lht-microservices/leader/db"
)

const (
	serviceName = "leaderservice"
	port = "9001"
)

func main() {
	logrus.Infof("Start Service [%s] \n", serviceName)

	// init bolt db
	initBoltDB()

	// init web server
	service.StartWebServer(port)
}

func initBoltDB() {
	// init bolt db instance to interface
	service.BoltClient = &db.BoltClient{}

	// init boltdb database
	service.BoltClient.OpenBoltDB()

	// mock 10 fake leaders
	service.BoltClient.MockLeaderData()
}

