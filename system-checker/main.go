package main

import (
	"flag"
	"net/http"
	"os"
)

func main()  {

	//	init to receive arguments
	port := flag.String("port", "80", "port to check")
	flag.Parse()

	url := "http://localhost:" + *port + "/health"
	resp, err := http.Get(url)


	//	error case - exit & return 1
	if err != nil || resp.StatusCode != 200 {
		os.Exit(1)

		//	successful case - exit & return 0
	} else {
		os.Exit(0)
	}
}
