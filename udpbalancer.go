package main

import (
	"container/ring"
	"fmt"
	"log"
	"net"
	"sync"
)

type Backend struct {
	Host  string
	Port  int
	Sconn *net.UDPConn
	Dconn *net.UDPConn
}

func (bk *Backend) RunBackend() {
	daddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bk.Host, bk.Port))
	dconn, err := net.DialUDP("udp", nil, daddr)
	if err != nil {
		log.Printf("udpBalance RunBackend connect error: %s", err)
	}
	bk.Dconn = dconn
	var buffer [1500]byte
	for {
		n, err := bk.Dconn.Read(buffer[0:])
		if err != nil {
			log.Printf("udpBalance RunBackend error: %s", err)
		}
		_, err = bk.Sconn.Write(buffer[0:n])
		if err != nil {
			log.Printf("udpBalance RunBackend error 2: %s", err)
		}
	}
}

type BackendRouter interface {
	Choose() *Backend
	Add(*Backend)
	SetSconn(*net.UDPConn)
}

type UdpbackendSimple struct {
	rlist *ring.Ring
	mu    sync.RWMutex
}

func (br *UdpbackendSimple) Add(b *Backend) {
	br.mu.Lock()
	defer br.mu.Unlock()
	value := &ring.Ring{Value: b}
	if br.rlist == nil {
		br.rlist = value
	} else {
		br.rlist = br.rlist.Link(value).Next()
	}
	go b.RunBackend()
}

func (br *UdpbackendSimple) Choose() *Backend {
	br.mu.Lock()
	defer br.mu.Unlock()
	if br.rlist == nil {
		return nil
	}
	ret := br.rlist.Value.(*Backend)
	br.rlist = br.rlist.Next()
	return ret
}

func (br *UdpbackendSimple) SetSconn(sconn *net.UDPConn) {
	br.mu.Lock()
	defer br.mu.Unlock()
	r := br.rlist

	for i := br.rlist.Len(); i > 0; i-- {
		r = r.Next()
		bk := r.Value.(*Backend)
		bk.Sconn = sconn
	}

}

func udpBalance(port int, bkrouter BackendRouter) {
	log.Println("udpBalance start")
	saddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	udplisten, err := net.ListenUDP("udp", saddr)
	if err != nil {
		log.Printf("udpBalance Listen error: %s", err)
		return
	}
	bkrouter.SetSconn(udplisten)
	var buffer [1500]byte
	for {
		n, clientaddr, err := udplisten.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Printf("udpBalance read error: %s", err)
			continue
		}
		bk := bkrouter.Choose()
		_, err = bk.Dconn.Write(buffer[0:n])
		if err != nil {
			log.Printf("udpBalance write error: %s", err)
		}
	}
}

func main() {
	router := &UdpbackendSimple{}
	router.Add(&Backend{Host: "127.0.0.1", Port: 8001})
	router.Add(&Backend{Host: "127.0.0.1", Port: 8002})
	udpBalance(8000, router)
}
