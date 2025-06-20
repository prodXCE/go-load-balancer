package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Server struct {
	URL *url.URL
	Alive bool
	mux   sync.RWMutex
}

func (s *Server) SetAlive(alive bool) {
	s.mux.Lock()
	s.Alive = alive
	s.mux.Unlock()
}

func (s *Server) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.Alive
}

type ServerPool struct {
	servers []*Server
	current uint64
	mux     sync.Mutex
}

func (sp *ServerPool) AddServer(server *Server) {
	sp.servers = append(sp.servers, server)
}

func (sp *ServerPool) GetNextPeer() *Server {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	for i := 0; i < len(sp.servers); i++ {
		sp.current = (sp.current + 1) % uint64(len(sp.servers))
		if sp.servers[sp.current].IsAlive() {
			return sp.servers[sp.current]
		}
	}

	return nil
}

var serverPool ServerPool


func main() {
	// --- Server Pool Setup ---
	// List of backend server address
	backendAddress := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, addr := range backendAddress {
		// The `url.Parse` function takes a string and returns a URL object.
		// It's safer than just using strings for addresses
		serverUrl, err := url.Parse(addr)
		if err != nil {
			log.Fatalf("Could not parse server URL: %v", err)
		}

		// Creating a new server struct for each backend address
		// im using the pointer (*Server) so that when a change is need for server status
		// the change is reflected everywhere
		serverPool.AddServer(&Server{
			URL:   serverUrl,
			Alive: true, // Assume servers are alive initially
		})
		log.Printf("Configured server: %s", serverUrl)
	}

	// --- End of server Pool Setup ---



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

	peer := serverPool.GetNextPeer()
	if peer == nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Forwarding request to %s", peer.URL)
	fmt.Fprintf(w, "Hello!, This is the load balancer skeleton")
}
