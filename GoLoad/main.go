package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct{
	addr string
	proxy *httputil.ReverseProxy
	currentConnection int
	avgResponseTime int
}

func newSimpleServer(addr string) *simpleServer{
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr: addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type Server interface{
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type LoadBalancer struct{
	port string
	roundRobinCount int
	servers []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer{
	return &LoadBalancer{
		port: port,
		roundRobinCount: 0,
		servers: servers,
	}
} 

func handleErr(err error){
	if err != nil {
		fmt.Printf("error: %v\n",err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string {
	return s.addr;
}

func (s *simpleServer) IsAlive() bool{
	return true
} 

func (s *simpleServer) Serve(rw http.ResponseWriter, req  *http.Request){
	s.proxy.ServeHTTP(rw, req)
} 

func (lb *LoadBalancer) getNextServer() Server{
	server:= lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive(){
		lb.roundRobinCount++
		server= lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request){
	targetSever := lb.getNextServer()
	fmt.Printf("forwarding too %s", targetSever.Address())
	targetSever.Serve(rw, req)
}

func main(){
	server := []Server{
		newSimpleServer("https://www.facebook.com/"),
		newSimpleServer("https://www.google.com/"),
		newSimpleServer("https://www.youtube.com/"),
	}
	lb := NewLoadBalancer("8000", server)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request){
		lb.serveProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)
	fmt.Printf("serving is started")
	http.ListenAndServe(":" + lb.port, nil)
}