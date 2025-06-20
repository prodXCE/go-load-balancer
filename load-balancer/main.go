package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Server now holds a ReverseProxy instance. Each server will have its own
// reverse proxy that knows how to forward requests to it.
type Server struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive sets the alive status of the server.
func (s *Server) SetAlive(alive bool) {
	s.mux.Lock()
	s.Alive = alive
	s.mux.Unlock()
}

// IsAlive checks if the server is alive.
func (s *Server) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.Alive
}

// this basically needed
// ServerPool holds the information about our backends.
type ServerPool struct {
	servers []*Server
	current uint64
	mux     sync.Mutex
}

// AddServer adds a new server to the server pool.
func (sp *ServerPool) AddServer(server *Server) {
	sp.servers = append(sp.servers, server)
}

// GetNextPeer selects the next available and healthy server.
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

// HealthCheck pings the backends and updates their status
// This is a method on the ServerPool
func (sp *ServerPool) HealthCheck() {
	log.Println("Starting health check . . . ")
	for _, s := range sp.servers {
		server := s
		go func() {
			client := http.Client{Timeout: 2 * time.Second}
			resp, err := client.Get(server.URL.String())

			if err != nil || resp.StatusCode != http.StatusOK {
				server.SetAlive(false)
				log.Printf("Server %s is down", server.URL)
				return
			}

			if !server.IsAlive() {
				log.Printf("Server %s is back up.", server.URL)
			}

			server.SetAlive(true)
		}()
	}
	log.Println("Health check completed")
}

// We'll declare our serverPool as a global variable.
var serverPool ServerPool

// main is the entry point of our application.
func main() {
	// List of backend server addresses.
	backendAddresses := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, addr := range backendAddresses {
		serverUrl, err := url.Parse(addr)
		if err != nil {
			log.Fatalf("Could not parse server URL: %v", err)
		}

		// For each server, create its dedicated reverse proxy.
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)

		serverPool.AddServer(&Server{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s", serverUrl)
	}

	// This is our main handler for the load balancer.
	http.HandleFunc("/", lb)

	log.Println("Load Balancer listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

// lb is our load balancer handler.
func lb(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s", r.RemoteAddr)

	peer := serverPool.GetNextPeer()
	if peer == nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Forwarding request to %s", peer.URL)
	peer.ReverseProxy.ServeHTTP(w, r)
}
