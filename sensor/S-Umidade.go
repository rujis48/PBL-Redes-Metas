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
		umid := rand.Intn(41) + 30 // 30-70
		payload := fmt.Sprintf("UMID:%d", umid)
		conn.Write([]byte(payload))
		time.Sleep(5 * time.Second)
	}
}