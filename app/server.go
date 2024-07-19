package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type HTTPRequest struct {
	Method    string
	Path      string
	Headers   map[string]string
	Body      string
	UserAgent string
}

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
		response := ""
		req := make([]byte, 1024)
		conn.Read(req)
		request := getHttpRequest([]byte(req))
		if request.Method != "GET" {
			response = getStatus(405, "Method Not Allowed")
			conn.Write([]byte(response))
			conn.Close()
		}

		switch {
		case request.Path == "/user-agent":
			response = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(request.UserAgent), request.UserAgent)
		case strings.HasPrefix(request.Path, "/echo/"):
			echo := strings.Split(string(request.Path), "/echo/")[1:][0]
			response = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(echo), echo)
		case request.Path != "/":
			response = getStatus(404, "Not Found")
		default:
			response = getStatus(200, "OK")
		}
		fmt.Println("Response ", string([]byte(response)))
		conn.Write([]byte(response))
		conn.Close()
	}
}

func getStatus(code int, message string) string {
	return fmt.Sprintf("HTTP/1.1 %d %s", code, message)
}
func getHttpRequest(req []byte) HTTPRequest {
	reqStr := string(req)
	reqLines := strings.Split(reqStr, "\r\n")
	reqLine := strings.Split(reqLines[0], " ")
	method := reqLine[0]
	path := reqLine[1]
	headers := make(map[string]string)
	for _, line := range reqLines[1:] {
		if line == "" {
			break
		}
		header := strings.Split(line, ": ")
		headers[header[0]] = header[1]
	}
	body := reqLines[len(reqLines)-1]
	return HTTPRequest{
		Method:    method,
		Path:      path,
		Headers:   headers,
		Body:      body,
		UserAgent: headers["User-Agent"],
	}
}
