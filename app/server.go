package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
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
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}

	// if err != nil {
	// 	fmt.Println("Error accepting connection: ", err.Error())
	// 	os.Exit(1)
	// }
}
