package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

var temperaturaAtual = 18
var acLigado = false

func main() {
	go func() {
		addr, _ := net.ResolveUDPAddr("udp", ":6000") // Porta para status
		conn, _ := net.ListenUDP("udp", addr)
		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			msg := strings.TrimSpace(string(buf[:n]))
			if msg == "AC_ON" {
				acLigado = true
			} else if msg == "AC_OFF" {
				acLigado = false
			}
		}
	}()

	addr, _ := net.ResolveUDPAddr("udp", "interpretador:5000")
	conn, _ := net.DialUDP("udp", nil, addr)

	for {
		if acLigado {
			// Se o ar está ligado, a temperatura desce (pode não descer ou descer de maneira aleatoria entre -1 e -2)
			temperaturaAtual -= (rand.Intn(2)+1)
		} else {
			// Se o ar está desligado, a temperatura sobe (pode não subir ou subir enntre +1 a +3)
			temperaturaAtual += (rand.Intn(3)+1)
		}

		// Limites 
		if temperaturaAtual < 15 { temperaturaAtual = 15 }
		if temperaturaAtual > 40 { temperaturaAtual = 40 }

		payload := fmt.Sprintf("TEMP:%d", int(temperaturaAtual))
		conn.Write([]byte(payload))
		
		fmt.Printf("Temp: %d°C | AC Ligado: %v\n", int(temperaturaAtual), acLigado)
		time.Sleep(3 * time.Second)
	}
}