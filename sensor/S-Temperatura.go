package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "interpretador:5000")
	conn, _ := net.DialUDP("udp", nil, addr)
	for {
		temp := rand.Intn(40) // Temperaturas de 0-40
		payload := fmt.Sprintf("TEMP:%d", temp)
		conn.Write([]byte(payload))
		time.Sleep(3 * time.Second)
	}
}