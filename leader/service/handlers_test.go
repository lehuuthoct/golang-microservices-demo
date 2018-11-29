package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"net/http/httptest"
	"lehuuthoct/lht-microservices/leader/db"
	"lehuuthoct/lht-microservices/leader/model"
	"fmt"
	"encoding/json"
	"gopkg.in/h2non/gock.v1"
)

var mockDB *db.MockBoltClient

func initInitialData()  {
	mockDB = &db.MockBoltClient{}
	gock.InterceptClient(httpClient)
}

func TestGetLeaderByIdIncorrectPathReturnsNotFound(t *testing.T) {

	Convey("Given a url /notvalid/123  to get leader by id ", t, func() {
		req := httptest.NewRequest("GET", "/notvalid/1", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the Router", func() {
			NewRouter().ServeHTTP(resp, req)

			Convey("Then the response should be a 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
		})
	})
}

func TestGetLeaderByIDSuccessReturnsLeader(t *testing.T)  {
	/*
		Mock Behaviors Before Testing
	*/
	initInitialData()

	// init gock to intercepts requests for Quote service
	quote := `{"quote": "never stop learning!", "ipAddress":"10.1.1.1:8080", "language":"en"}`
	gock.New("http://quotes-service:8080").Get("/api/quote").MatchParam("strength","4").Reply(200).BodyString(quote)
	defer gock.Off()

	//	init mock instance
	leader1 := model.Leader{Id: "1", Name: "Leader_1"}
	method := "FindLeaderById"

	//	mock behavior for successful response (Leader: leader1, err: nil)
	mockDB.On(method, "1").Return(leader1, nil)

	//	 mock behavior for error response
	mockDB.On(method, "99").Return(model.Leader{}, fmt.Errorf("No leader is found"))

	// init boltdb client
	BoltClient = mockDB


	Convey("Given a HTTP request for /leader/1", t, func() {
		req := httptest.NewRequest("GET", "/leader/1", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the Router", func() {
			NewRouter().ServeHTTP(resp, req)

			Convey("Then the response status should be a 200", func() {
				// assert response code
				So(resp.Code, ShouldEqual, 200)

				//	assert received leader data
				leader := &model.Leader{}
				json.Unmarshal(resp.Body.Bytes(), &leader)
				So(leader.Id, ShouldEqual, "1")
				So(leader.Name, ShouldEqual, "Leader_1")
				So(leader.Quote.Text, ShouldEqual, "never stop learning!")
			})
		})
	})

	Convey("Given a HTTP request for /leader/99", t, func() {
		req := httptest.NewRequest("GET", "/leader/99", nil)
		resp := httptest.NewRecorder()
		Convey("When the request is handled by the Router", func() {
			NewRouter().ServeHTTP(resp, req)
			Convey("Then the response status should be a 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
		})
	})

}

func TestHealthCheckReturnsOK(t *testing.T)  {
	initInitialData()

	mockDB.On("Check").Return(true)
	BoltClient = mockDB

	Convey("Given a HTTP Request for /health", t, func() {
		req := httptest.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is served", func() {
			NewRouter().ServeHTTP(resp, req)
			Convey("Then expect 200 OK", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

func TestHealthCheckReturnsFalseWhenDBNotWorking(t *testing.T) {
	initInitialData()

	mockDB.On("Check").Return(false)
	BoltClient = mockDB
	Convey("Given a HTTP request for /health", t, func() {
		req := httptest.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is served", func() {
			NewRouter().ServeHTTP(resp, req)
			Convey("Then expect 503 Service Unavailable", func() {
				So(resp.Code, ShouldEqual, 503)
			})
		})
	})
}
