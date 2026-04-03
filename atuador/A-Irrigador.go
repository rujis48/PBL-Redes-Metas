package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var valvulaAberta bool = false // Estado da Irrigação

func main() {
	ln, _ := net.Listen("tcp", ":8070")
	fmt.Println("[IRRIGADOR] Atuador iniciado. Válvula: FECHADA")

	for {
		conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		comando := strings.TrimSpace(msg)

		if comando == "LIGAR" {
			valvulaAberta = true
		} else if comando == "DESLIGAR" {
			valvulaAberta = false
		} else {
			umid, err := strconv.Atoi(comando)
			if err == nil {
				if umid < 40 {
					valvulaAberta = true
				} else if umid > 60 {
					valvulaAberta = false
				}
			}
		}

		fmt.Printf("[STATUS IRRIGADOR] Válvula Aberta: %v\n", valvulaAberta)
		conn.Close()
	}
}