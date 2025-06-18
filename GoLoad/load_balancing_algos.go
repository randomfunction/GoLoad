package main

import (
	// "fmt"
	"net/http"
	// "net/http/httputil"
	// "net/url"
	// "os"
)

type Algorithms interface {
	// weightedRoundRobin() Server
	ipHashing(r *http.Request) Server
	leastConnections() Server
	leastResponseTime() Server
}

func (lb *LoadBalancer) ipHashing(r *http.Request) Server {
	clientIP := r.RemoteAddr
	hash := 0
	for i := 0; i < len(clientIP); i++ {
		hash += int(clientIP[i])
	}
	return lb.servers[hash%len(lb.servers)]
}

func (lb *LoadBalancer) leastConnections() Server {
	var selected Server
	minConn := int(^uint(0) >> 1) 

	for _, s := range lb.servers {
		ss := s.(*simpleServer)
		if ss.currentConnection < minConn {
			selected = ss
			minConn = ss.currentConnection
		}
	}
	return selected
}

func (lb *LoadBalancer) leastResponseTime() Server {
	var selected Server
	minResp := int(^uint(0) >> 1)

	for _, s := range lb.servers {
		ss := s.(*simpleServer)
		if ss.avgResponseTime < minResp {
			selected = ss
			minResp = ss.avgResponseTime
		}
	}
	return selected
}

