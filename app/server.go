package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Program has started!")

	l, err := net.Listen("tcp", "127.0.0.1:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		fmt.Println("Accepted connection from ", conn.RemoteAddr())

		req := make([]byte, 1024)
		conn.Read(req)
		if !strings.HasPrefix(string(req), "GET / HTTP/1.1") {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			conn.Close()
			return
		}
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}
}
