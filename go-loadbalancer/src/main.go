package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, Servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         Servers,
	}
}
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool { return true }

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	Server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !Server.IsAlive() {
		lb.roundRobinCount++
		Server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return Server

}

func (lb *LoadBalancer) serverProxy(rw http.ResponseWriter, req *http.Request) {
	targetserver := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %q\n", targetserver.Address())
	targetserver.Serve(rw, req)
}

func main() {
	Servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("http://www.youtube.com"),
		newSimpleServer("http://duckduckgo.com"),
	}

	lb := NewLoadBalancer("8000", Servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serverProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving request at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
