package main

import "github.com/kamilernerd/toxy"

func main() {
	config := toxy.LoadConfig()
	resolver := toxy.ServiceResolver(config)

	go resolver.Resolve()

	serviceLoadBalancer := toxy.NewBalancer(resolver.Services, config.LoadBalancer)

	server := toxy.Server{
		Config:   config,
		Balancer: serviceLoadBalancer,
	}

	server.TcpListener()
}
