package toxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Proxy struct {
	service ResolverService
	conn    net.Conn
	OutBuf  chan []byte
}

func NewProxy(service ResolverService) *Proxy {
	proxy := &Proxy{
		service: service,
		OutBuf:  make(chan []byte),
	}
	return proxy
}

func (p *Proxy) connect() {
	var err error

	// if p.service.State == "up" {
	p.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", p.service.Hostname, p.service.Port))
	if err != nil {
		p.OutBuf <- nil
		return
	}
	// }
}

func (p *Proxy) write(buf []byte) {
	if p.conn == nil {
		return
	}

	n, err := p.conn.Write(buf)
	if err != nil {
		log.Println(n, err)
		p.conn.Close()
		return
	}
}

func (p *Proxy) read() {
	buf := make([]byte, 4096)
	for {
		if p.conn == nil {
			p.OutBuf <- nil
			log.Printf("Proxy connection is nil")
			return
		}
		n, err := p.conn.Read(buf)
		r := make([]byte, n)
		copy(r, buf[:n])
		p.OutBuf <- r
		if err != nil {
			if errors.Is(err, io.EOF) {
				p.OutBuf <- nil
				p.conn.Close()
				return
			}
			log.Printf("%v", err)
			p.conn.Close()
			return
		}
	}
}

func (p *Proxy) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
