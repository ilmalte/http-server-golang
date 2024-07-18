package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

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
		// read 4096 bytes at a time
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read data from connection")
			os.Exit(1)
			continue
		}
		s := ""
		if n > 0 {
			l := string(buf[:n])
			i := strings.Index(l, " HTTP")
			if i != -1 {
				s = l[:i]
				i = strings.Index(l, "/")
				s = s[i:]
				fmt.Println("Request: ", s)
			}
		}

		fmt.Println("Len: ", len(s))
		if len(s) > 1 {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			continue
		}
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}
}
