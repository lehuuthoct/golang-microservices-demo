package model

type Leader struct {
	Id string `json:"id"`
	Name string `json:"name"`
	FromHost string `json:"fromHost"`
	Quote Quote `json:"quote"`
}

// Quote json tag is based on Quote service [springboot]
type Quote struct {
	Text string `json:"quote"`
	FromHost string `json:"ipAddress"`
	Language string `json:"language"`
}







