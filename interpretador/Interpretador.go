package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var modoAutomatico = true

func main() {
	go servidorSensores()
	go servidorCliente()

	fmt.Println("Interpretador Online: Sensores(5000/UDP), AC(8070/TCP), Irrigador(8070/TCP), Cliente(8080/TCP)")
	select {}
}

func servidorSensores() {
	addr, _ := net.ResolveUDPAddr("udp", ":5000")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, _, _ := conn.ReadFromUDP(buf)
		msg := strings.TrimSpace(string(buf[:n]))
		partes := strings.Split(msg, ":")
		if len(partes) < 2 { continue }

		tipo, valor := partes[0], partes[1]
		fmt.Printf("[UDP] Recebido %s: %s\n", tipo, valor)

		if modoAutomatico {
			if tipo == "TEMP" {
				enviarTCP("atuador_ac:8070", valor)
			} else if tipo == "UMID" {
				enviarTCP("irrigador:8070", valor)
			}
		}
		enviarUDPCliente("cliente:7000", msg)
	}
}

func servidorCliente() {
	ln, _ := net.Listen("tcp", ":8080")
	for {
		conn, _ := ln.Accept()
		go func(c net.Conn) {
			defer c.Close()
			msg, _ := bufio.NewReader(c).ReadString('\n')
			cmd := strings.TrimSpace(msg)
			if cmd == "AUTO_ON" { modoAutomatico = true } 
			if cmd == "AUTO_OFF" { modoAutomatico = false }
			fmt.Printf("[TCP] Comando Cliente: %s (Auto: %v)\n", cmd, modoAutomatico)
		}(conn)
	}
}

func enviarTCP(target, msg string) {
	conn, err := net.Dial("tcp", target)
	if err == nil {
		fmt.Fprintf(conn, msg+"\n")
		conn.Close()
	}
}

func enviarUDPCliente(target, msg string) {
// Verifica se o target não está vazio
    if target == "" {
        fmt.Println("[ERRO] Target do cliente vazio. Ignorando envio.")
        return
    }

    // Resolver o endereço UDP
    addr, err := net.ResolveUDPAddr("udp", target)
    if err != nil {
        fmt.Printf("[ERRO] Falha ao resolver endereço %s: %v\n", target, err)
        return
    }
    
    // Abre a conexão
    c, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        fmt.Printf("[ERRO] Falha ao conectar via UDP: %v\n", err)
        return
    }
    defer c.Close()

    _, err = c.Write([]byte(msg))
    if err != nil {
        fmt.Printf("[ERRO] Falha ao escrever dados: %v\n", err)
    }
}