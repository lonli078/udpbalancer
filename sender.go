package main

import (

	"log"
	"time"
	"net"
	"fmt"

)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8000")
	conn, err := net.DialUDP("udp", nil, addr)
	for i := 1; i <= 100; i++ {
		_, err = conn.Write([]byte(fmt.Sprintf("%d", i)))
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}
