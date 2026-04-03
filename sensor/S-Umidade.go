package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

// Estado inicial da simulação
var (
	umidadeAtual  = 50
	irrigadorOn   = false
)

func main() {
	// 1. Goroutine para ouvir o status do Irrigador enviado pelo Interpretador (Porta 6000)
	go func() {
		addr, _ := net.ResolveUDPAddr("udp", ":6000") 
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("Erro ao abrir porta 6000:", err)
			return
		}
		defer conn.Close()
		
		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			msg := strings.TrimSpace(string(buf[:n]))
			
			// Atualiza o estado interno baseado no feedback do sistema
			if msg == "IRRIG_ON" {
				irrigadorOn = true
			} else if msg == "IRRIG_OFF" {
				irrigadorOn = false
			}
		}
	}()

	// 2. Conexão para envio de dados ao Interpretador
	addrInterp, _ := net.ResolveUDPAddr("udp", "interpretador:5000")
	connInterp, _ := net.DialUDP("udp", nil, addrInterp)
	defer connInterp.Close()

	fmt.Println("[SENSOR UMID] Iniciado. Monitorando solo...")

	for {
		// 3. Lógica de Simulação Física
		if irrigadorOn {
			// Se o irrigador está ligado, a umidade sobe (entre 1 e 3 unidades)
			umidadeAtual += rand.Intn(3) + 1
		} else {
			// Se está desligado, o solo seca (perda aleatória de 0 a 2 unidades)
			umidadeAtual -= rand.Intn(3)
		}

		// 4. Aplicando os Limites de Segurança (Garantir 0-100%)
		if umidadeAtual < 20 { umidadeAtual = 20 }
		if umidadeAtual > 90 { umidadeAtual = 90 }

		// 5. Envio do Dado
		payload := fmt.Sprintf("UMID:%d", umidadeAtual)
		connInterp.Write([]byte(payload))
		
		// Print no terminal idêntico ao sensor de temperatura
		fmt.Printf("Umid: %d%% | Irrigador: %v\n", umidadeAtual, irrigadorOn)
		
		// Frequência de leitura (mesma do sensor de temperatura para consistência)
		time.Sleep(3 * time.Second)
	}
}