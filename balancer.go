package toxy

import (
	"log"
	"math/rand"
)

type Balancer struct {
	services     []*ResolverService
	state        []int
	balancerType string
}

func NewBalancer(services []*ResolverService, balancerType string) *Balancer {
	b := &Balancer{
		services:     services,
		state:        []int{},
		balancerType: balancerType,
	}
	return b
}

func (b *Balancer) random() *ResolverService {
	return b.services[rand.Intn(len(b.services)-0)+0]
}

func (b *Balancer) initSequential() {
	for k := range b.services {
		b.state = append(b.state, k)
	}
}

func (b *Balancer) sequential() *ResolverService {
	for k := range b.state {
		if k+1 < len(b.state) {
			if b.state[k] < b.state[k+1] {
				b.state[k]++
				if b.services[k].State != ServiceUp {
					continue
				}
				return b.services[k]
			}
		} else {
			b.state[k]++
			return b.services[k]
		}
	}
	return nil
}

func (b *Balancer) getServiceAt(index int) *ResolverService {
	return b.services[index]
}

func (b *Balancer) selectService() *ResolverService {
	switch b.balancerType {
	case Random:
		return b.random()
	case Sequential:
		return b.sequential()
	default:
		return b.getServiceAt(0)
	}
}

func (b *Balancer) getCurrentState() {
	log.Printf("%+v \n", b.state)
	log.Printf("%+v \n", b.services)
}
