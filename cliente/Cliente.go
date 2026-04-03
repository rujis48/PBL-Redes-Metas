package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	valorTemp   = "--"
	valorUmid   = "--"
	statusAC    = "OFF"
	statusIrrig = "OFF"
	modoAuto    = "ATIVADO"
	filtro      = "AMBOS"
)

func main() {
	// Limpa tela, move cursor para o topo e limpa histórico de rolagem
	fmt.Print("\033[2J\033[H\033[3J")
	
	go escutarInterpretador()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		desenharInterface()
		// Posiciona o cursor de entrada sempre na mesma linha (12)
		fmt.Print("\033[12;1H\033[2KEscolha uma opção: ")
		if !scanner.Scan() { break }
		processarComando(strings.TrimSpace(scanner.Text()))
	}
}

func escutarInterpretador() {
	addr, _ := net.ResolveUDPAddr("udp", ":7000")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, _, _ := conn.ReadFromUDP(buf)
		msg := strings.TrimSpace(string(buf[:n]))

		// Processa dados e status
		if strings.HasPrefix(msg, "TEMP:") {
			valorTemp = strings.Split(msg, ":")[1]
		} else if strings.HasPrefix(msg, "UMID:") {
			valorUmid = strings.Split(msg, ":")[1]
		} else if msg == "AC_ON" { statusAC = "ON "
		} else if msg == "AC_OFF" { statusAC = "OFF"
		} else if msg == "IRRIG_ON" { statusIrrig = "ON "
		} else if msg == "IRRIG_OFF" { statusIrrig = "OFF"
		}

		// Atualiza apenas a linha do Monitor (Linha 2)
		fmt.Printf("\033[s\033[2;1H\033[2K\033[36m[ MONITOR ] TEMP: %s°C | AC: [%s] | UMID: %s%% | IRRIG: [%s]\033[0m\033[u", 
			valorTemp, statusAC, valorUmid, statusIrrig)
	}
}

func desenharInterface() {
	fmt.Print("\033[1;1H") // Volta ao topo
	fmt.Println("\033[33m================== PAINEL DE CONTROLE PBL ==================\033[0m")
	fmt.Printf("\033[2;1H") // Linha do monitor (preenchida pela goroutine)
	fmt.Println("\n\033[32m----------------------- MENU PRINCIPAL ---------------------\033[0m")
	fmt.Printf(" MODO: %s | FILTRO: %s\n", modoAuto, filtro)
	fmt.Println(" 1. Ver Temp  | 2. Ver Umid  | 3. Ver Ambos")
	fmt.Println(" A. AUTO ON   | M. AUTO OFF  | Q. Sair")
	fmt.Println("------------------------------------------------------------")
	fmt.Println(" L. Ligar AC  | D. Desligar AC")
	fmt.Println(" I. Ligar IRR | F. Desligar IRR")
	fmt.Println("------------------------------------------------------------")
}

func processarComando(opt string) {
	opt = strings.ToUpper(opt)
	
	conn, err := net.Dial("tcp", "interpretador:8080")
	if err != nil { return }
	defer conn.Close()

	switch opt {
	case "1":
		fmt.Fprint(conn, "VER_TEMP\n")
		filtro = "TEMP"
	case "2":
		fmt.Fprint(conn, "VER_UMID\n")
		filtro = "UMID"
	case "3":
		fmt.Fprint(conn, "VER_AMBOS\n")
		filtro = "AMBOS"
	case "A":
		fmt.Fprint(conn, "AUTO_ON\n")
		modoAuto = "ATIVADO"
	case "M":
		fmt.Fprint(conn, "AUTO_OFF\n")
		modoAuto = "MANUAL "
	case "L": fmt.Fprint(conn, "AC_ON\n")
	case "D": fmt.Fprint(conn, "AC_OFF\n")
	case "I": fmt.Fprint(conn, "IRRIG_ON\n")
	case "F": fmt.Fprint(conn, "IRRIG_OFF\n")
	case "Q": os.Exit(0)
	}
}