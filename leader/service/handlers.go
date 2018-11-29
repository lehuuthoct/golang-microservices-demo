package service

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"github.com/sirupsen/logrus"
	"net"
	"lehuuthoct/lht-microservices/leader/model"
	"io/ioutil"
	"fmt"
)

// init variables
// this help to init connection to other service
var httpClient = &http.Client{}

// this help to update system state
var isSystemHealthy = true

type ServerHealthResponse struct {
	DBStatus string `json:"db_status"`
}

// init service configuration
func init() {

	//	disable connection keep-alive to avoid load-balancing issue
	transport := &http.Transport{DisableKeepAlives:true}
	httpClient.Transport = transport

}

/*
	This method is used to find leader by id
	from leader service
*/
func FindLeaderByID(w http.ResponseWriter, r *http.Request) {
	//	get leader id
	var leaderId = mux.Vars(r)["leaderId"]

	// get leader from db
	leader, err := BoltClient.FindLeaderById(leaderId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// init host ip
	leader.FromHost = findHostIP()

	// init quote
	quote, err := findLeaderQuote()
	if err == nil {
		leader.Quote = quote
	}

	//	init leader json data for resp
	leaderJSON, _ := json.Marshal(leader)
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.Header().Set(CONTENT_LENGTH, strconv.Itoa(len(leaderJSON)))
	w.WriteHeader(http.StatusOK)
	w.Write(leaderJSON)
}

/*
	This method is used to find address
	including IP data of leader service
*/
func findHostIP() string {
	addressList, err := net.InterfaceAddrs()
	if err != nil {
		return "Error finding host"
	}

	for _, address := range addressList {
		//	display address if not a loopback
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	panic("Cannot find IP Address!")
}

/*
	This method is used to check current health state
	of the system running health service
	including database status
*/
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	// check if db is working
	isDBUp := BoltClient.Check()
	status := &ServerHealthResponse{}

	logrus.Infof("isDBUp [%v], isSystemHealthy [%v]", isDBUp, isSystemHealthy)

	if isDBUp && isSystemHealthy {
		status.DBStatus = "UP"
		jsonData, _ := json.Marshal(status)
		initJSONResponse(w, http.StatusOK, jsonData)

		// notify db is not working
	} else {
		status.DBStatus = "Down"
		jsonData, _ := json.Marshal(status)
		initJSONResponse(w, http.StatusServiceUnavailable, jsonData)
	}
}

/*
	This method init json reqponse to user
*/
func initJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.Header().Set(CONTENT_LENGTH, strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}


/*
	This method is used to update system state
	to test feature automatically create new container
	by docker Swarm
*/
func SetDBHealthState(w http.ResponseWriter, r *http.Request) {
	//	get state from path
	state, err := strconv.ParseBool(mux.Vars(r)["state"])

	//	init error if error parsing
	if err != nil {
		logrus.Info("Invalid request to set DB Health! Only true or false are allowed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//	mutate db state
	isSystemHealthy = state
	w.WriteHeader(http.StatusOK)
}

/*
	This method is used to find quote
	from other microservice [Quote]
	developed by spring-boot project
	running in the same Docker Swarm Cluster
*/
func findLeaderQuote() (model.Quote, error)  {

	quote := model.Quote{}

	// init request to quote service
	url := "http://quotes-service:8080/api/quote?strength=4"
	req, err := http.NewRequest("GET", url , nil)

	// get response from quote service
	resp, err := httpClient.Do(req)

	if err != nil {
		return quote, err
	}

	if resp.StatusCode == 200 {

		// read quote data as bytes[] to memory
		quotesBytes, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(quotesBytes, &quote)

		return quote, nil
	} else {
		return quote, fmt.Errorf("Error finding quote %s", url)
	}

}
