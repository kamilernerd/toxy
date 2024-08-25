package toxy

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
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

func (s *Server) TcpListener() {
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
			continue
		}

		tlscon, ok := conn.(*tls.Conn)
		if ok {
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		go s.connectionHandler(conn)
	}
}

func (s *Server) connectionHandler(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		_, err := r.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatal(err)
			}
		}

		// TODO - Loadbalance servers
		// TODO - Connect, write and read response then return to the proxy

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Fatal(n, err)
		}
	}
}
