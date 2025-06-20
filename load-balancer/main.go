package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)

	// web server start point
	// tells the server to listen on port 8080 for any incoming connections
	// the second argument, `nil`, tells it to use the default handler I just set up.
	log.Println("Load Balancer skeleton listening on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved request from %s", r.RemoteAddr)

	fmt.Fprintf(w, "Hello!, This is the load balancer skeleton")
}
