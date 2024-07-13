package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	baseURL = "http://localhost:808"
)

type LoadBalancer struct {
	RevProxy httputil.ReverseProxy
}

type Endpoints struct {
	List []*url.URL
}

func (e *Endpoints) Shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func MakeLoadBalancer(amount int) {
	// Instantiate Objects
	var lb LoadBalancer
	var ep Endpoints

	// Server + Router
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	// Creating the endpoints
	for i := 0; i < amount; i++ {
		endpoint, err := createEndpoint(baseURL, i)
		if err != nil {
			log.Fatalf("Failed to create endpoint: %v", err)
		}
		ep.List = append(ep.List, endpoint)
	}

	// Handler Functions
	router.HandleFunc("/loadbalancer", makeRequest(&lb, &ep))

	// Listen and serve
	log.Println("Load balancer running on :8090")
	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *LoadBalancer, ep *Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for !testServer(ep.List[0].String()) {
			ep.Shuffle()
		}
		lb.RevProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.Shuffle()
		lb.RevProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, idx int) (*url.URL, error) {
	link := endpoint + strconv.Itoa(idx)
	parsedURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	return parsedURL, nil
}

func testServer(endpoint string) bool {
	resp, err := http.Get(endpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
