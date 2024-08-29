package main

import "github.com/kamilernerd/toxy"

func main() {
	config := toxy.LoadConfig()
	resolver := toxy.ServiceResolver(config)

	go resolver.Resolve()

	server := toxy.Server{
		Config:   config,
		Services: resolver.Services,
	}

	server.TcpListener()
}
