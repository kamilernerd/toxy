package toxy_test

import (
	"crypto/tls"
	"fmt"
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

func TestTcpDial(t *testing.T) {

}
