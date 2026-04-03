package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var visualizacao = "AMBOS" // Opções: TEMP, UMID, AMBOS

func main() {
	// Goroutine para receber dados em tempo real
	go func() {
		addr, _ := net.ResolveUDPAddr("udp", ":7000")
		conn, _ := net.ListenUDP("udp", addr)
		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			msg := string(buf[:n])
			
			// Lógica de filtro de exibição
			exibir := false
			if visualizacao == "AMBOS" || 
			   (visualizacao == "TEMP" && strings.Contains(msg, "TEMP")) ||
			   (visualizacao == "UMID" && strings.Contains(msg, "UMID")) {
				exibir = true
			}

			if exibir {
				// \033[s salva a posição do cursor, \033[u restaura
				fmt.Printf("\033[s\033[H\033[2K\r[DADOS ATUAIS] %s\033[u", msg)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		desenharMenu()
		fmt.Print("Escolha uma opção: ")
		if !scanner.Scan() { break }
		
		opcao := strings.TrimSpace(scanner.Text())
		processarComando(opcao)
	}
}

func desenharMenu() {
	fmt.Println("\n\033[32m======= SISTEMA DE CONTROLE PBL =======\033[0m")
	fmt.Printf("Visualização Atual: [%s]\n", visualizacao)
	fmt.Println("1. Ver apenas Temperatura (AC)")
	fmt.Println("2. Ver apenas Umidade (Irrigador)")
	fmt.Println("3. Ver Ambos")
	fmt.Println("---------------------------------------")
	fmt.Println("A. Ativar MODO AUTOMÁTICO")
	fmt.Println("M. Desativar MODO AUTOMÁTICO (Manual)")
	fmt.Println("---------------------------------------")
	fmt.Println("L. Ligar AC (Manual)    | D. Desligar AC (Manual)")
	fmt.Println("I. Ligar Irrigador (M)  | F. Desligar Irrigador (M)")
	fmt.Println("Q. Sair")
}

func processarComando(opt string) {
	conn, err := net.Dial("tcp", "interpretador:8080")
	if err != nil {
		fmt.Println("Erro ao conectar com o Interpretador")
		return
	}
	defer conn.Close()

	switch strings.ToUpper(opt) {
	case "1": visualizacao = "TEMP"
	case "2": visualizacao = "UMID"
	case "3": visualizacao = "AMBOS"
	case "A": fmt.Fprint(conn, "AUTO_ON\n")
	case "M": fmt.Fprint(conn, "AUTO_OFF\n")
	case "L": fmt.Fprint(conn, "AC_ON\n")
	case "D": fmt.Fprint(conn, "AC_OFF\n")
	case "I": fmt.Fprint(conn, "IRRIG_ON\n")
	case "F": fmt.Fprint(conn, "IRRIG_OFF\n")
	case "Q": os.Exit(0)
	default: fmt.Println("Opção inválida!")
	}
}