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
		umid, _ := strconv.Atoi(strings.TrimSpace(msg))
		if umid < 40 { fmt.Printf("[IRRIGADOR] LIGADO (%d%%)\n", umid)
		} else { fmt.Printf("[IRRIGADOR] DESLIGADO (%d%%)\n", umid) }
		conn.Close()
	}
}