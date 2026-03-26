package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	ln, _ := net.Listen("tcp", ":8070")
	for {
		conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		temp, _ := strconv.Atoi(strings.TrimSpace(msg))
		if temp > 25 { fmt.Printf("[AC] LIGADO (%d°C)\n", temp)
		} else if temp < 18 { fmt.Printf("[AC] DESLIGADO (%d°C)\n", temp) }
		conn.Close()
	}
}