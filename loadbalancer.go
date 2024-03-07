package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// PriorityServerQueue represents a priority queue of servers based on their active connections.
type PriorityServerQueue []*Server

func (pq PriorityServerQueue) Len() int           { return len(pq) }
func (pq PriorityServerQueue) Less(i, j int) bool { return pq[i].Connections < pq[j].Connections }
func (pq PriorityServerQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityServerQueue) Push(x interface{}) {
	item := x.(*Server)
	*pq = append(*pq, item)
}

func (pq *PriorityServerQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// LoadBalancer represents a simple load balancer.
type LoadBalancer struct {
	servers          []*Server
	mutex            sync.Mutex
	connectionsCount map[*Server]int
	priorityQueue    *PriorityServerQueue
}

// Server represents a server in the load balancer.
type Server struct {
	Address     string
	Status      bool
	Connections int
}

// NewLoadBalancer creates a new load balancer.
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		servers:          make([]*Server, 0),
		mutex:            sync.Mutex{},
		connectionsCount: make(map[*Server]int),
		priorityQueue:    &PriorityServerQueue{},
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

// ChooseServer selects the server with the least active connections.
func (lb *LoadBalancer) ChooseServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	heap.Init(lb.priorityQueue) // Initialize the priority queue

	// Populate the priority queue with healthy servers
	for _, server := range lb.servers {
		if server.Status {
			heap.Push(lb.priorityQueue, server)
		}
	}

	// Get the server with the least active connections
	if lb.priorityQueue.Len() > 0 {
		minConnectionsServer := heap.Pop(lb.priorityQueue).(*Server)
		minConnectionsServer.Connections++
		return minConnectionsServer
	}

	return nil // No healthy servers found
}

// ReleaseConnection releases a connection from a server.
func (lb *LoadBalancer) ReleaseConnection(server *Server) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	// Decrement the active connections count for the server
	if lb.connectionsCount[server] > 0 {
		lb.connectionsCount[server]--
	}
}

// AddServer adds a new server to the load balancer.
func (lb *LoadBalancer) AddServer(server *Server) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.servers = append(lb.servers, server)
	lb.connectionsCount[server] = 0
}

func main() {
	fmt.Println("Load Balancer Implementation in Golang")

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
				lb.ReleaseConnection(chosenServer)
			} else {
				fmt.Printf("Request %d: No healthy servers available\n", requestNum)
			}
		}(i)
	}

	// Allow time for health checks and requests to complete
	time.Sleep(10 * time.Second)
}
