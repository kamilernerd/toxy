package toxy_test

import (
	// "bufio"
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
		}
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", defaultConfStruct.Port))
	defer conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTcpListenerConnectionHandler(t *testing.T) {
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

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", defaultConfStruct.Port+1), &tls.Config{Certificates: []tls.Certificate{cert}})

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			defer conn.Close()

			if err != nil {
				log.Println(err)
				break
			}
			go HandleTcpListenerHelper(conn)
		}
	}()

	DialTestHelper(t)
}

func HandleTcpListenerHelper(conn net.Conn) {
	defer conn.Close()

	var buf [4096]byte

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("EOF reading tcp data %v", err)
				break
			}
			log.Printf("%v", err)
			break
		}

		if n == 0 {
			log.Fatal("Could not receive any bytes")
			return
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

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", defaultConfStruct.Hostname, defaultConfStruct.Port+1), conf)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	testData := `
GET / HTTP/1.1
Host: localhost:8081
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Sec-GPC: 1
Connection: keep-alive
Upgrade-Insecure-Requests: 1
Cache-Control: no-cache
	`

	n, err := conn.Write([]byte(testData))
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
