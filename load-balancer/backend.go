package main

import (
	"fmt"
	"log"
	"net/http"
)

func startBackendServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Backend :%d] Recieved request", port)
		fmt.Fprintf(w, "Hello from backend server on port: %d\n", port)
	}

	log.Printf("Starting backend server on %s", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(handler)); err != nil {
		log.Fatalf("Failed to start backend server on %s: %v", addr, err)
	}
}

func main() {
	go startBackendServer(8081)
	go startBackendServer(8082)
	go startBackendServer(8083)

	select {}
}
