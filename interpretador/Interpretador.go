package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	modoAutomatico = true
	filtroCliente  = "AMBOS" // Controle de fluxo externo
)

func main() {
	go servidorSensores()
	go servidorCliente()

	fmt.Println("\033[34m[INTERPRETADOR]\033[0m Online e Roteando...")
	fmt.Println("-> UDP:5000 (Sensores) | TCP:8080 (Cliente)")
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

		// --- LÓGICA DE FILTRO EXTERNO PARA O CLIENTE ---
		enviarParaCliente := false
		switch filtroCliente {
		case "AMBOS":
			enviarParaCliente = true
		case "TEMP":
			if tipo == "TEMP" { enviarParaCliente = true }
		case "UMID":
			if tipo == "UMID" { enviarParaCliente = true }
		}

		if enviarParaCliente {
			enviarUDPCliente("cliente:7000", msg)
		}

		// --- LÓGICA AUTOMÁTICA ---
		if modoAutomatico {
			if tipo == "TEMP" {
				enviarEConfirmarTCP("atuador_ac:8070", valor)
			} else if tipo == "UMID" {
				enviarEConfirmarTCP("irrigador:8070", valor)
			}
		}
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

			switch cmd {
			case "VER_TEMP":
				filtroCliente = "TEMP"
			case "VER_UMID":
				filtroCliente = "UMID"
			case "VER_AMBOS":
				filtroCliente = "AMBOS"
			case "AUTO_ON":
				modoAutomatico = true
			case "AUTO_OFF":
				modoAutomatico = false
			case "AC_ON":
				enviarEConfirmarTCP("atuador_ac:8070", "LIGAR")
			case "AC_OFF":
				enviarEConfirmarTCP("atuador_ac:8070", "DESLIGAR")
			case "IRRIG_ON":
				enviarEConfirmarTCP("irrigador:8070", "LIGAR")
			case "IRRIG_OFF":
				enviarEConfirmarTCP("irrigador:8070", "DESLIGAR")
			}
			fmt.Printf("[LOG] Cliente definiu: %s | Auto: %v\n", cmd, modoAutomatico)
		}(conn)
	}
}

func enviarEConfirmarTCP(target, msg string) {
	conn, err := net.DialTimeout("tcp", target, 2*time.Second)
	if err != nil { return }
	defer conn.Close()

	fmt.Fprintf(conn, msg+"\n")

	// Espera o Atuador responder (Feedback Síncrono)
	resposta, err := bufio.NewReader(conn).ReadString('\n')
	if err == nil {
		status := strings.TrimSpace(resposta)
		if status != "" {
			// Notifica Cliente sobre a mudança de estado (ON/OFF)
			enviarUDPCliente("cliente:7000", status)
			// Se for AC, notifica o sensor para mudar a temperatura física
			if strings.Contains(status, "AC") {
				avisarSensorFisico("sensor_temp:6000", status)
			}
		}
	}
}

func enviarUDPCliente(target, msg string) {
	addr, _ := net.ResolveUDPAddr("udp", target)
	c, err := net.DialUDP("udp", nil, addr)
	if err == nil {
		c.Write([]byte(msg))
		c.Close()
	}
}

func avisarSensorFisico(target, status string) {
	addr, _ := net.ResolveUDPAddr("udp", target)
	c, err := net.DialUDP("udp", nil, addr)
	if err == nil {
		c.Write([]byte(status))
		c.Close()
	}
}