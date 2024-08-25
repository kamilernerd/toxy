package toxy

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Hostname string
	Port     int
	CertPath string
	KeyPath  string
}

func (s *Server) LoadCertificates() *tls.Config {
	cert, err := tls.LoadX509KeyPair(s.CertPath, s.KeyPath)
	if err != nil {
		log.Fatalf("Failed to load keypair %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func (s *Server) TCPListener() {
	tlsConfig := s.LoadCertificates()
	ln, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", s.Hostname, s.Port), tlsConfig)
	defer ln.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		defer conn.Close()

		if err != nil {
			log.Fatal(err)
			conn.Close()
		}

		go s.connectionHandler(conn)
	}
}

func (s *Server) connectionHandler(conn net.Conn) {
	fmt.Printf("%s\n", conn.LocalAddr())
}
