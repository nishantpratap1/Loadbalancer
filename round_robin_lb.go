package main

import (
	"fmt"
	"sync"
	"time"
)

// LoadBalancer represents a simple load balancer.
type LoadBalancer struct {
	servers  []*Server
	mutex    sync.Mutex
	next     int // Index of the next server to which the request will be routed
}

// Server represents a server in the load balancer.
type Server struct {
	Address string
	Status  bool
}

// NewLoadBalancer creates a new load balancer.
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		servers: make([]*Server, 0),
		mutex:   sync.Mutex{},
		next:    0,
	}
}

// HealthCheck periodically checks the health of each server.
func (lb *LoadBalancer) HealthCheck() {
	for {
		time.Sleep(5 * time.Second)

		lb.mutex.Lock()
		for _, server := range lb.servers {
			// Implement your health check logic here
			// For simplicity, we assume all servers are healthy
			server.Status = true
		}
		lb.mutex.Unlock()
	}
}

// ChooseServer selects the next server in a round-robin manner.
func (lb *LoadBalancer) ChooseServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.servers) == 0 {
		return nil
	}

	chosenServer := lb.servers[lb.next]
	lb.next = (lb.next + 1) % len(lb.servers)
	return chosenServer
}

// AddServer adds a new server to the load balancer.
func (lb *LoadBalancer) AddServer(server *Server) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.servers = append(lb.servers, server)
}

func main() {
	fmt.Println("Load Balancer Implementation with Round-Robin in Golang")

	// Create a new load balancer
	lb := NewLoadBalancer()

	// Add servers to the load balancer
	server1 := &Server{Address: "Server1", Status: true}
	server2 := &Server{Address: "Server2", Status: true}
	server3 := &Server{Address: "Server3", Status: true}
	lb.AddServer(server1)
	lb.AddServer(server2)
	lb.AddServer(server3)

	// Start health checks in a goroutine
	go lb.HealthCheck()

	// Simulate incoming requests
	for i := 0; i < 10; i++ {
		go func(requestNum int) {
			chosenServer := lb.ChooseServer()
			if chosenServer != nil {
				fmt.Printf("Request %d routed to %s\n", requestNum, chosenServer.Address)
				time.Sleep(1 * time.Second) // Simulate processing time
			} else {
				fmt.Printf("Request %d: No servers available\n", requestNum)
			}
		}(i)
	}

	// Allow time for health checks and requests to complete
	time.Sleep(10 * time.Second)
}
