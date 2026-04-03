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
	filtroCliente  = "AMBOS"
)

func main() {
	go servidorSensores()
	go servidorCliente()
	fmt.Println("\033[34m[INTERPRETADOR]\033[0m Online | UDP:5000 | TCP:8080")
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

		// Filtro de envio para o Cliente
		if (filtroCliente == "AMBOS") || (filtroCliente == "TEMP" && tipo == "TEMP") || (filtroCliente == "UMID" && tipo == "UMID") {
			enviarUDPCliente("cliente:7000", msg)
		}

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
			case "VER_TEMP":  filtroCliente = "TEMP"
			case "VER_UMID":  filtroCliente = "UMID"
			case "VER_AMBOS": filtroCliente = "AMBOS"
			case "AUTO_ON":   modoAutomatico = true
			case "AUTO_OFF":  modoAutomatico = false
			case "AC_ON":     enviarEConfirmarTCP("atuador_ac:8070", "LIGAR")
			case "AC_OFF":    enviarEConfirmarTCP("atuador_ac:8070", "DESLIGAR")
			case "IRRIG_ON":  enviarEConfirmarTCP("irrigador:8070", "LIGAR")
			case "IRRIG_OFF": enviarEConfirmarTCP("irrigador:8070", "DESLIGAR")
			}
		}(conn)
	}
}

func enviarEConfirmarTCP(target, msg string) {
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil { return }
	defer conn.Close()

	fmt.Fprintf(conn, msg+"\n")
	resposta, err := bufio.NewReader(conn).ReadString('\n')
	if err == nil {
		status := strings.TrimSpace(resposta)
		if status == "" { return }

		enviarUDPCliente("cliente:7000", status)

		// REPASSE DE FEEDBACK PARA OS SENSORES (Porta 6000)
		if strings.Contains(status, "AC") {
			avisarSensor("sensor_temp:6000", status)
		} else if strings.Contains(status, "IRRIG") {
			avisarSensor("sensor_umid:6000", status) // O sensor de umidade deve ter este hostname
		}
	}
}

func enviarUDPCliente(target, msg string) {
	addr, _ := net.ResolveUDPAddr("udp", target)
	c, _ := net.DialUDP("udp", nil, addr)
	if c != nil {
		c.Write([]byte(msg))
		c.Close()
	}
}

func avisarSensor(target, status string) {
	addr, _ := net.ResolveUDPAddr("udp", target)
	c, _ := net.DialUDP("udp", nil, addr)
	if c != nil {
		c.Write([]byte(status))
		c.Close()
	}
}