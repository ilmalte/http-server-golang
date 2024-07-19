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
	response := getEmptyResponseStr(405, "Method Not Allowed")
	switch {
	case request.Method == "POST":
		if strings.HasPrefix(request.Path, "/files/") {
			response = postFileContent(request.Path, request.Body)
		}
	case request.Method == "GET":
		switch {
		case request.Path == "/user-agent":
			response = getResponseStr(request.UserAgent, "text/plain")
		case strings.HasPrefix(request.Path, "/echo/"):
			echo := strings.Split(string(request.Path), "/echo/")[1:][0]
			response = getResponseStr(echo, "text/plain")
		case strings.HasPrefix(request.Path, "/files/"):
			response = getFileContent(request.Path)
		case request.Path != "/":
			response = getEmptyResponseStr(404, "Not Found")
		default:
			response = getEmptyResponseStr(200, "OK")
		}
	default:
		break
	}
	return response
}
func postFileContent(path string, content string) string {
	err := error(nil)
	directory := ""
	fileName := strings.Split(string(path), "/files/")[1:][0]

	if len(os.Args) > 1 {
		directory = os.Args[2]
		if directory == "" {
			fmt.Println("No directory specified")
		}
		err = postFile(directory, fileName, content)
	}
	if err != nil {
		return getEmptyResponseStr(500, "Internal Server Error")
	}
	return getEmptyResponseStr(201, "Created")
}
func postFile(directory string, fileName string, content string) error {
	err := os.WriteFile(directory+fileName, []byte(content), 0755)
	if err != nil {
		fmt.Print("unable to write file: %w", err)
		return err
	}
	return nil
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
		return getEmptyResponseStr(500, "Internal Server Error")
	}
	if file == nil {
		return getEmptyResponseStr(404, "Not Found")
	}
	fileContent, err := os.ReadFile(directory + file.Name())
	if err != nil {
		return getEmptyResponseStr(500, "Internal Server Error")
	}
	return getResponseStr(string(fileContent), "application/octet-stream")
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
func getEmptyResponseStr(code int, message string) string {
	return fmt.Sprintf("%s\r\n\r\n", getStatus(code, message))
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
	body := strings.Trim(reqLines[len(reqLines)-1], "\x00")
	return HTTPRequest{
		Method:    method,
		Path:      path,
		Headers:   headers,
		Body:      body,
		UserAgent: headers["User-Agent"],
	}
}
