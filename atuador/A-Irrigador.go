package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var estaLigado bool = false

func main() {
	ln, _ := net.Listen("tcp", ":8070")
	fmt.Println("[IRRIGADOR] Online...")

	for {
		conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		comando := strings.TrimSpace(msg)

		if comando == "LIGAR" || comando == "IRRIG_ON" {
			estaLigado = true
		} else if comando == "DESLIGAR" || comando == "IRRIG_OFF" {
			estaLigado = false
		} else {
			umid, err := strconv.Atoi(comando)
			if err == nil {
				if umid < 35 {
					estaLigado = true
				} else if umid > 75 {
					estaLigado = false
				}
			}
		}

		// Resposta para o Interpretador
		if estaLigado {
			fmt.Fprint(conn, "IRRIG_ON\n")
		} else {
			fmt.Fprint(conn, "IRRIG_OFF\n")
		}
		conn.Close()
	}
}