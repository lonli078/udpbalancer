package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func recipient(port int) {
	saddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	udplisten, _ := net.ListenUDP("udp", saddr)
	var buffer [10]byte
	for {
		n, clientaddr, err := udplisten.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Printf("udpBalance read error: %s", err)
			continue
		}
		log.Println(port)
		log.Println(n)
		log.Println(buffer)
		log.Println(clientaddr)
	}
}

func main() {
	go recipient(8001)
	go recipient(8002)
	for {
		time.Sleep(10000000 * time.Second)
	}

}
