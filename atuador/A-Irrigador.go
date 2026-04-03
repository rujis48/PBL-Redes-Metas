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
    fmt.Println("[IRRIGADOR] Atuador iniciado e aguardando comandos...")

    for {
        conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
        comando := strings.TrimSpace(msg)

        // LOG DE DEBUG
        fmt.Printf("[DEBUG] Recebido: '%s'\n", comando)

        if comando == "LIGAR" || comando == "IRRIG_ON" {
            estaLigado = true
            fmt.Println("COMANDO MANUAL: LIGANDO")
        } else if comando == "DESLIGAR" || comando == "IRRIG_OFF" {
            estaLigado = false
            fmt.Println("COMANDO MANUAL: DESLIGANDO")
        } else {
            // Se não for comando de texto, tenta ler como umidade (número)
            umid, err := strconv.Atoi(comando)
            if err == nil {
                // Lógica Automática solicitada: < 35 liga | > 75 desliga
                if umid < 35 {
                    estaLigado = true
                } else if umid > 75 {
                    estaLigado = false
                }
            } else {
                fmt.Printf("[ERRO] Mensagem inválida: %s\n", comando)
            }
        }

        // FEEDBACK SÍNCRONO: Responde na mesma conexão para o Interpretador repassar
        if estaLigado {
            fmt.Fprint(conn, "IRRIG_ON\n")
        } else {
            fmt.Fprint(conn, "IRRIG_OFF\n")
        }

        status := "OFF"
        if estaLigado { status = "ON" }
        fmt.Printf("[ESTADO ATUAL] Irrigador: %s\n", status)
        
        conn.Close()
    }
}