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
	fmt.Println("[AC] Atuador iniciado e aguardando comandos...")

	for {
		conn, _ := ln.Accept()
		// Usamos scanner para ler a linha inteira corretamente
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		comando := strings.TrimSpace(msg)

		// LOG DE DEBUG: Importante para você ver o que chegou no container
		fmt.Printf("[DEBUG] Recebido: '%s'\n", comando)

		if comando == "LIGAR" || comando == "AC_ON" {
			estaLigado = true
			fmt.Println("COMANDO MANUAL: LIGANDO")
		} else if comando == "DESLIGAR" || comando == "AC_OFF" {
			estaLigado = false
			fmt.Println("COMANDO MANUAL: DESLIGANDO")
		} else {
			// Se não for comando de texto, tentamos ler como temperatura (número)
			temp, err := strconv.Atoi(comando)
			if err == nil {
				// Lógica Automática baseada em números
				if temp > 25 {
					estaLigado = true
				} else if temp < 18 {
					estaLigado = false
				}
			} else {
				fmt.Printf("[ERRO] Mensagem inválida: %s\n", comando)
			}
		}

		if estaLigado {
            fmt.Fprint(conn, "AC_ON\n")
        } else {
            fmt.Fprint(conn, "AC_OFF\n")
        }

        status := "OFF"
        if estaLigado { status = "ON" }
        fmt.Printf("[ESTADO ATUAL] Ar-Condicionado: %s\n", status)
        
        conn.Close()
    }
}