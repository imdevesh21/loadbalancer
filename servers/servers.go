package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	Ports []int
	mu    sync.Mutex
}

func (s *ServerList) Populate(amount int) {
	if amount > 10 {
		log.Fatal("Amount of ports can't exceed 10")
	}

	for x := 0; x < amount; x++ {
		s.Ports = append(s.Ports, x)
	}
}

func (s *ServerList) Pop() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Ports) == 0 {
		log.Fatal("No ports available")
	}
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServers(amount int) {
	// ServerList Object
	var myServerList ServerList
	myServerList.Populate(amount)

	// Waitgroup
	var wg sync.WaitGroup
	wg.Add(amount)

	for x := 0; x < amount; x++ {
		go makeServers(&myServerList, &wg)
	}

	wg.Wait()
}

func makeServers(sl *ServerList, wg *sync.WaitGroup) {
	defer wg.Done()

	// Router
	r := http.NewServeMux()

	// Server
	port := sl.Pop()

	server := &http.Server{
		Addr:    fmt.Sprintf(":808%d", port),
		Handler: r,
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server %d", port)
	})

	r.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Server Shut Down!"))
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	})

	log.Printf("Starting server on port :808%d", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("ListenAndServe error: %v", err)
	}
}
