package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	go func() {
		addr, _ := net.ResolveUDPAddr("udp", ":7000")
		conn, _ := net.ListenUDP("udp", addr)
		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			fmt.Printf("\r[DADOS] %s          \n> ", string(buf[:n]))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Comandos: 'on' (Auto ON), 'off' (Auto OFF)")
	for scanner.Scan() {
		cmd := strings.ToLower(scanner.Text())
		conn, _ := net.Dial("tcp", "interpretador:8080")
		if cmd == "on" { fmt.Fprint(conn, "AUTO_ON\n")
		} else if cmd == "off" { fmt.Fprint(conn, "AUTO_OFF\n") }
		conn.Close()
		fmt.Print("> ")
	}
}