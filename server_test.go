package toxy_test

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/kamilernerd/toxy"
)

func TestLoadCertificates(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tls.LoadX509KeyPair(defaultConfStruct.CertPath, defaultConfStruct.KeyPath)
	if err != nil {
		t.Fatalf("Failed to load keypair %v", err)
	}
}

func TestTcpListener(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := tls.LoadX509KeyPair(defaultConfStruct.CertPath, defaultConfStruct.KeyPath)
	if err != nil {
		t.Fatalf("Failed to load keypair %v", err)
	}

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", defaultConfStruct.Port), &tls.Config{Certificates: []tls.Certificate{cert}})

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			defer conn.Close()

			if err != nil {
				t.Fatal(err)
				conn.Close()
			}

			go func() {
				conn.Close()
			}()
		}
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", defaultConfStruct.Port))
	defer conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTcpListenerConnectionHanlder(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := tls.LoadX509KeyPair(defaultConfStruct.CertPath, defaultConfStruct.KeyPath)
	if err != nil {
		t.Fatalf("Failed to load keypair %v", err)
	}

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", defaultConfStruct.Port), &tls.Config{Certificates: []tls.Certificate{cert}})

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			defer conn.Close()

			if err != nil {
				t.Fatal(err)
				break
			}
			go HandleTcpListenerHelper(conn)
		}
	}()

	DialTestHelper(t)
}

func HandleTcpListenerHelper(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		_, err := r.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatal(err)
			}
		}

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Fatal(n, err)
		}
	}
}

func DialTestHelper(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", defaultConfStruct.Hostname, defaultConfStruct.Port), conf)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		t.Fatal(n, err)
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatal(n, err)
	}

	if len(string(buf[:n])) == 0 {
		t.Fatal("Response is empty")
	}
}
