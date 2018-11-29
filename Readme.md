# Introduction
- This codebase demonstrates how to develop microservices using golang

- The examples include Leader and Quote services which are deployed as docker containers to a docker Swarm.

- The services communicate with each other via (HTTP) endpoints.   
 
 
# Packages
## common
- for bolt key-value db 
>
	go get github.com/boltdb/bolt

- for testing
>
	go get github.com/smartystreets/goconvey
	go get github.com/stretchr/testify/mock
	go get gopkg.in/h2non/gock.v1

## module-specific
>
	dep ensure -add github.com/gorilla/mux
	dep ensure -add github.com/boltdb/bolt

## testing
### go-convey-web-ui
- run (default: port 8080)
		~/go/bin/goconvey
- [reference](http://goconvey.co/)

# performance
- Gatling can be used to test performance of different services 
>
	mvn gatling:execute -Dusers=1000 -Dduration=30 -DbaseUrl=http://localhost:8888
