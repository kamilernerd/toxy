package main

import "github.com/kamilernerd/toxy"

var Config = toxy.LoadConfig()

func main() {
	server := toxy.Server{
		Port:     Config.Port,
		Hostname: Config.Hostname,
		CertPath: Config.CertPath,
		KeyPath:  Config.KeyPath,
	}
	server.TcpListener()
}
