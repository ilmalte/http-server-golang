package main

import (
	"fmt"
	"io/fs"
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
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	req := make([]byte, 1024)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Println("Failed to read request")
		return
	}
	httpReq := getHttpRequest(req)
	fmt.Println(httpReq)
	response := getResponse(httpReq)
	conn.Write([]byte(response))
}
func getResponse(request HTTPRequest) string {
	response := ""
	if request.Method != "GET" {
		response = fmt.Sprintf("%s\r\n\r\n", getStatus(405, "Method Not Allowed"))
		return response
	}
	switch {
	case request.Path == "/user-agent":
		response = getResponseStr(request.UserAgent, "text/plain")
	case strings.HasPrefix(request.Path, "/echo/"):
		echo := strings.Split(string(request.Path), "/echo/")[1:][0]
		response = getResponseStr(echo, "text/plain")
	case strings.HasPrefix(request.Path, "/files/"):
		fileContent := getFileContent(request.Path)
		response = getResponseStr(fileContent, "application/octet-stream")
	case request.Path != "/":
		response = fmt.Sprintf("%s\r\n\r\n", getStatus(404, "Not Found"))
	default:
		response = fmt.Sprintf("%s\r\n\r\n", getStatus(200, "OK"))
	}
	return response
}
func getFileContent(path string) string {
	err := error(nil)
	directory := ""
	file := fs.DirEntry(nil)
	fileName := strings.Split(string(path), "/files/")[1:][0]

	if len(os.Args) > 1 {
		directory = os.Args[2]
		if directory == "" {
			fmt.Println("No directory specified")
		}
		file, err = getFile(directory, fileName)
	}
	if err != nil {
		return fmt.Sprintf("%s\r\n\r\n", getStatus(500, "Internal Server Error"))
	}
	if file == nil {
		return fmt.Sprintf("%s\r\n\r\n", getStatus(404, "Not Found"))
	}
	fileContent, err := os.ReadFile(directory + file.Name())
	if err != nil {
		return fmt.Sprintf("%s\r\n\r\n", getStatus(500, "Internal Server Error"))
	}
	return string(fileContent)
}
func getFile(directory string, fileName string) (fs.DirEntry, error) {
	// Check if the specified directory exists
	_, err := os.Stat(directory)
	if err != nil {
		fmt.Println("Directory does not exist")
		return nil, err
	}
	// List files in the specified directory
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Failed to read directory")
		return nil, err
	}
	if len(files) == 0 {
		fmt.Println("No files in the directory")
		return nil, err
	}
	// Find file by name
	for _, file := range files {
		if file.Name() == fileName {
			return file, nil
		}
	}
	return nil, nil
}
func getResponseStr(s string, t string) string {
	return fmt.Sprintf("%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), t, len(s), s)
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
