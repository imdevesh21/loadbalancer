package main

import (
	"loadbalancer/loadbalancer"
	// "loadbalancer/servers"
)
func main(){
	loadbalancer.MakeLoadBalancer(5)
	// servers.RunServers(5);
}