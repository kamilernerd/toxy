package toxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
)

type Server struct {
	Config   Config
	Services []ResolverService
}

/*
Load x509 keypair
*/
func (s *Server) LoadCertificates() *tls.Config {
	cert, err := tls.LoadX509KeyPair(s.Config.CertPath, s.Config.KeyPath)
	if err != nil {
		log.Fatalf("Failed to load keypair %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: false}
}

/*
Listen for TCP connections
*/
func (s *Server) TcpListener() {
	tlsConfig := s.LoadCertificates()
	ln, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Hostname, s.Config.Port), tlsConfig)
	defer ln.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		defer conn.Close()

		if err != nil {
			log.Printf("Error during accepting a remote connection %v\n", err)
			continue
		}
		go s.connectionHandler(conn)
	}
}

/*
Pipe tcp server to remote connection

Buffers host data and forwards to remote server
Buffers remote data and forwards to host server
*/
func (s *Server) connectionHandler(conn net.Conn) {
	defer conn.Close()

	selectedService := s.Services[0]
	switch s.Config.LoadBalancer {
	case "random":
		index := rand.Intn(len(s.Services)-0) + 0
		selectedService = s.Services[index]
		break
	case "roundrobin":
		break
	default:
		selectedService = s.Services[0]
	}

	proxy := NewProxy(selectedService)
	proxy.connect()
	go proxy.read()

	serverOutBuf := s.read(conn)

	for {
		select {
		case proxyBuf := <-proxy.OutBuf:
			if proxyBuf != nil {
				s.write(conn, proxyBuf)
			}
		case hostBuf := <-serverOutBuf:
			if hostBuf != nil {
				proxy.write(hostBuf)
			}
		}
	}
}

func (s *Server) write(conn net.Conn, buf []byte) {
	n, err := conn.Write(buf)
	if err != nil {
		log.Println(n, err)
		conn.Close()
		return
	}
}

func (s *Server) read(conn net.Conn) chan []byte {
	outBuf := make(chan []byte)
	go func() {
		buf := make([]byte, 4096)
		for {
			if conn == nil {
				log.Printf("Host connection is nil")
				return
			}
			n, err := conn.Read(buf)
			r := make([]byte, n)
			copy(r, buf[:n])
			outBuf <- r
			if err != nil {
				if errors.Is(err, io.EOF) {
					outBuf <- nil
					conn.Close()
					return
				}
				log.Printf("%v", err)
				conn.Close()
				return
			}
		}
	}()
	return outBuf
}
