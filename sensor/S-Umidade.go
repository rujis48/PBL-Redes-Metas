package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

var umidadeAtual = 50
var irrigadorOn  = false


func main() {
	go func() {
		addr, _ := net.ResolveUDPAddr("udp", ":6000")
		conn, _ := net.ListenUDP("udp", addr)
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			msg := strings.TrimSpace(string(buf[:n]))
			if msg == "IRRIG_ON" {
				irrigadorOn = true
			} else if msg == "IRRIG_OFF" {
				irrigadorOn = false
			}
		}
	}()

	addrInterp, _ := net.ResolveUDPAddr("udp", "interpretador:5000")
	connInterp, _ := net.DialUDP("udp", nil, addrInterp)
	defer connInterp.Close()

	for {
		if irrigadorOn {
			umidadeAtual += rand.Intn(3) + 1
		} else {
			umidadeAtual -= rand.Intn(2)
		}

		if umidadeAtual < 20 { umidadeAtual = 20 }
		if umidadeAtual > 90 { umidadeAtual = 90 }

		payload := fmt.Sprintf("UMID:%d", umidadeAtual)
		connInterp.Write([]byte(payload))
		
		fmt.Printf("Umid: %d%% | Irrigador: %v\n", umidadeAtual, irrigadorOn)
		time.Sleep(3 * time.Second)
	}
}