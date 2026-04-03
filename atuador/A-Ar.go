package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var estaLigado bool = false // O "Estado" do Ar-Condicionado

func main() {
	ln, _ := net.Listen("tcp", ":8070")
	fmt.Println("[AC] Atuador iniciado. Estado atual: DESLIGADO")

	for {
		conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		comando := strings.TrimSpace(msg)

		// 1. Processar Comandos Manuais (Texto)
		if comando == "LIGAR" {
			estaLigado = true
		} else if comando == "DESLIGAR" {
			estaLigado = false
		} else {
			// 2. Processar Dados Automáticos (Números)
			temp, err := strconv.Atoi(comando)
			if err == nil {
				if temp > 25 {
					estaLigado = true
				} else if temp < 18 {
					estaLigado = false
				}
			}
		}

		// Exibir o estado real da variável booleana
		statusStr := "DESLIGADO"
		if estaLigado {
			statusStr = "LIGADO"
		}
		fmt.Printf("[STATUS AC] Booleano: %t | Visual: %s\n", estaLigado, statusStr)
		
		conn.Close()
	}
}