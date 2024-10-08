package toxy

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type ResolverService struct {
	Port     int
	Hostname string
	Name     string
	State    string
}

type Resolver struct {
	Services []*ResolverService
	Quit     chan int
	sync     sync.Mutex
	interval int
}

func ServiceResolver(config Config) *Resolver {
	r := &Resolver{
		Services: []*ResolverService{},
		Quit:     make(chan int, 1),
		interval: config.ResolveInterval,
	}

	for _, serv := range config.Server {
		for _, v := range serv {
			r.Services = append(r.Services, &ResolverService{
				Port:     v.Port,
				Hostname: v.Hostname,
				Name:     v.Name,
				State:    ServiceUnknown,
			})
		}
	}
	return r
}

func (r *Resolver) Resolve() {
	ticker := time.NewTicker(time.Duration(r.interval) * time.Second)

	for _, v := range r.Services {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", v.Hostname, v.Port))
		if err != nil {
			log.Printf("Error dialing service %s -> %v", v.Name, err)
			v.State = ServiceDown
			continue
		}
		v.State = ServiceUp
		conn.Close()
	}

	for {
		select {
		case <-ticker.C:
			r.sync.Lock()
			for _, v := range r.Services {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", v.Hostname, v.Port))
				if err != nil {
					log.Printf("Error dialing service %s -> %v", v.Name, err)
					v.State = ServiceDown
					continue
				}
				v.State = ServiceUp
				conn.Close()
			}
			r.sync.Unlock()
		case <-r.Quit:
			ticker.Stop()
			r.sync.Unlock()
			return
		}
	}
}

func (r *Resolver) StopResolver() {
	r.Quit <- 1
}
