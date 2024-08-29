package toxy_test

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kamilernerd/toxy"
)

func TestServiceResolver(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}

	r := &toxy.Resolver{
		Services: []toxy.ResolverService{},
	}

	for _, serv := range defaultConfStruct.Server {
		for _, v := range serv {
			r.Services = append(r.Services, toxy.ResolverService{
				Port:     v.Port,
				Hostname: v.Hostname,
				Name:     v.Name,
				State:    "unknown",
			})
		}
	}

	if len(r.Services) == 0 {
		t.Error("no services registered")
	}
}

func TestResolver(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}

	r := &toxy.Resolver{
		Services: []toxy.ResolverService{},
		Quit:     make(chan int, 1),
	}

	for _, serv := range defaultConfStruct.Server {
		for _, v := range serv {
			r.Services = append(r.Services, toxy.ResolverService{
				Port:     v.Port,
				Hostname: v.Hostname,
				Name:     v.Name,
				State:    "unknown",
			})
		}
	}

	if len(r.Services) == 0 {
		t.Error("no services registered")
	}

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

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			for _, v := range r.Services {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", v.Hostname, v.Port))
				if err != nil {
					t.Fatal(err)
				}

				_, err = conn.Write([]byte(testData))
				if err != nil {
					t.Fatal(err)
				}
				conn.(*net.TCPConn).CloseWrite()

				defer conn.Close()
			}
			r.Quit <- 1
		case <-r.Quit:
			ticker.Stop()
			return
		}
	}
}
