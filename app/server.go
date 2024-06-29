// ----------------------- Building a HTTP Server by stages in Go --------------------

// #1. bind to a port + #2. Respond with 200

// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// )

// func main() {
// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Println("Logs from your program will appear here!")
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}

// 	conn, err := l.Accept()

// 	if err != nil {
// 		fmt.Println("Error accepting connection: ", err.Error())
// 		os.Exit(1)
// 	}

// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

// }

// #3. Extract URL path

// package main
// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// )
// func main() {
// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Println("Logs from your program will appear here!")
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	conn, err := l.Accept()
// 	if err != nil {
// 		fmt.Println("Error accepting connection: ", err.Error())
// 		os.Exit(1)
// 	}
// 	req := make([]byte, 1024)
// 	conn.Read(req)
// 	if !strings.HasPrefix(string(req), "GET / HTTP/1.1") {
// 		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 		return
// 	}
// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
// }

// #4. Respond with a body

// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// )

// func main() {
// 	fmt.Println("Logs from your program will appear here!")
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	defer l.Close()
// 	for {
// 		connection, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			continue
// 		}
// 		go handleConnection(connection)
// 	}
// }
// func handleConnection(connection net.Conn) {
// 	defer connection.Close()
// 	requestBuffer := make([]byte, 1024)
// 	n, err := connection.Read(requestBuffer)
// 	if err != nil {
// 		fmt.Println("Failed to read the request:", err)
// 		return
// 	}
// 	fmt.Printf("Request: %s\n", requestBuffer[:n])
// 	request := string(requestBuffer[:n])
// 	path := strings.Split(request, " ")[1]
// 	if path == "/" {
// 		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
// 	} else if strings.Split(path, "/")[1] == "echo" {
// 		message := strings.Split(path, "/")[2]
// 		connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
// 	} else {
// 		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 	}
// }

// #5. Read Headers

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// )

// type HTTPRequest struct {
// 	Method    string
// 	Path      string
// 	Headers   map[string]string
// 	Body      string
// 	UserAgent string
// }

// func main() {
// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Println("Logs from your program will appear here!")
// 	// Uncomment this block to pass the first stage
// 	//
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	conn, err := l.Accept()
// 	if err != nil {
// 		fmt.Println("Error accepting connection: ", err.Error())
// 		os.Exit(1)
// 	}
// 	// status, err := bufio.NewReader(conn).ReadString('\n')
// 	scanner := bufio.NewScanner(conn)
// 	req, err := parseStatus(scanner)
// 	// fmt.Println(req.Headers)
// 	if err := scanner.Err(); err != nil {
// 		fmt.Fprintln(conn, "reading standard input:", err)
// 	}
// 	var response string
// 	switch path := req.Path; {
// 	case strings.HasPrefix(path, "/echo/"):
// 		content := strings.TrimLeft(path, "/echo/")
// 		response = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(content), content)
// 	case path == "/user-agent":
// 		response = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(req.UserAgent), req.UserAgent)
// 	case path == "/":
// 		response = getStatus(200, "OK") + "\r\n\r\n"
// 	default:
// 		response = getStatus(404, "Not Found") + "\r\n\r\n"
// 	}
// 	conn.Write([]byte(response))
// 	// fmt.Fprintln(conn)
// 	conn.Close()
// }
// func parseStatus(scanner *bufio.Scanner) (*HTTPRequest, error) {
// 	var req HTTPRequest = HTTPRequest{}
// 	req.Headers = make(map[string]string)
// 	for i := 0; scanner.Scan(); i++ {
// 		if i == 0 {
// 			parts := strings.Split(scanner.Text(), " ")
// 			req.Method = parts[0]
// 			req.Path = parts[1]
// 			continue
// 		}
// 		headers := strings.Split(scanner.Text(), ": ")
// 		if len(headers) < 2 {
// 			req.Body = headers[0]
// 			break
// 		}
// 		if headers[0] == "User-Agent" {
// 			req.UserAgent = headers[1]
// 		}
// 		req.Headers[headers[0]] = headers[1]
// 	}
// 	return &req, nil
// }
// func getStatus(statusCode int, statusText string) string {
// 	return fmt.Sprintf("HTTP/1.1 %d %s", statusCode, statusText)
// }

// #6. Support Concurrent Connections

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"net"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// var statusCodeToString = map[int]string{
// 	200: "OK",
// 	404: "Not Found",
// }

// type Request struct {
// 	Method  string
// 	Path    string
// 	Headers map[string]string
// 	Body    string
// }
// type Response struct {
// 	StatusCode int
// 	Headers    map[string]string
// 	Body       string
// }

// func (r Response) String() string {
// 	statusText, ok := statusCodeToString[r.StatusCode]
// 	if !ok {
// 		statusText = "Unknown"
// 	}
// 	// No headers so assume plain text result.
// 	if r.Headers == nil {
// 		r.Headers = map[string]string{
// 			"Content-Type": "text/plain",
// 		}
// 	}
// 	// Figure out content length if not set.
// 	if _, ok = r.Headers["Content-Length"]; !ok {
// 		r.Headers["Content-Length"] = strconv.Itoa(len(r.Body))
// 	}
// 	var headerString strings.Builder
// 	for k, v := range r.Headers {
// 		headerString.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
// 	}
// 	return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s", r.StatusCode, statusText, headerString.String(), r.Body)
// }
// func main() {
// 	l, err := net.Listen("tcp", ":4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		// Spawn a go thread
// 		go handleConnection(conn)
// 	}
// }
// func handleConnection(conn net.Conn) {
// 	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
// 	defer func(conn net.Conn) {
// 		err := conn.Close()
// 		if err != nil {
// 			fmt.Println("Failed to close connection: ", err.Error())
// 		}
// 	}(conn)
// 	request, err := parseRequest(stream.Reader)
// 	if err != nil {
// 		fmt.Println("Failed to parse request: ", err.Error())
// 		os.Exit(1)
// 	}
// 	response := Response{StatusCode: 200}
// 	if strings.HasPrefix(request.Path, "/echo") {
// 		pathParts := strings.SplitN(request.Path, "/echo/", 2)
// 		response.Body = pathParts[1]
// 	} else if request.Path == "/user-agent" {
// 		userAgent := request.Headers["User-Agent"]
// 		response.Body = userAgent
// 	} else if request.Path != "/" {
// 		response.StatusCode = 404
// 	}
// 	_, err = stream.WriteString(response.String())
// 	if err != nil {
// 		fmt.Println("Failed to write to socket: ", err.Error())
// 		os.Exit(1)
// 	}
// 	err = stream.Flush()
// 	if err != nil {
// 		fmt.Println("Failed to flush to socket")
// 		os.Exit(1)
// 	}
// }
// func parseRequest(reader *bufio.Reader) (Request, error) {
// 	request := Request{
// 		Headers: make(map[string]string),
// 	}
// 	firstLine, err := reader.ReadString('\n')
// 	if err != nil {
// 		return Request{}, fmt.Errorf("malformed HTTP request")
// 	}
// 	parts := strings.Split(firstLine, " ")
// 	request.Method = parts[0]
// 	request.Path = parts[1]
// 	for {
// 		curLine, err := reader.ReadString('\n')
// 		if curLine == "\r\n" {
// 			break
// 		}
// 		if err == io.EOF {
// 			return request, nil
// 		} else if err != nil {
// 			return Request{}, err
// 		}
// 		headerParts := strings.SplitN(curLine, ":", 2)
// 		request.Headers[headerParts[0]] = strings.TrimSpace(headerParts[1])
// 	}
// 	// If the content length is set read the body.
// 	contentLenStr, ok := request.Headers["Content-Length"]
// 	if !ok {
// 		return request, nil
// 	}
// 	contentLen, _ := strconv.Atoi(contentLenStr)
// 	buf := make([]byte, contentLen)
// 	// This probably should read in chunks.
// 	_, err = io.ReadFull(reader, buf)
// 	if err != nil {
// 		return request, err
// 	}
// 	request.Body = string(buf)
// 	return request, nil
// }

// #7. Return a file

// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// )

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 1024)
// 	conn.Read(buf)
// 	fmt.Print(string(buf))
// 	req := string(buf)
// 	path := strings.Split(req, "\r\n")[0]
// 	path = strings.TrimSpace(path)
// 	path = strings.Split(path, " ")[1]
// 	var pathUA string
// 	headers := strings.Split(req, "\r\n")
// 	for _, header := range headers {
// 		if strings.HasPrefix(header, "User-Agent:") {
// 			uaParts := strings.SplitN(header, ":", 2)
// 			pathUA = strings.TrimSpace(uaParts[1])
// 			break
// 		}
// 	}
// 	response := ""
// 	if strings.HasPrefix(req, "GET / HTTP") {
// 		response = "HTTP/1.1 200 OK\r\n\r\n"
// 	} else if strings.Contains(req, "/echo/") {
// 		echo := strings.TrimPrefix(path, "/echo/")
// 		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echo), echo)
// 	} else if strings.Contains(req, "/user-agent") && pathUA != "" {
// 		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(pathUA), pathUA)
// 	} else if strings.Contains(req, "/files/") {
// 		dir := os.Args[2]
// 		fileName := strings.TrimPrefix(path, "/files/")
// 		fmt.Print(fileName)
// 		data, err := os.ReadFile(dir + fileName)
// 		if err != nil {
// 			response = "HTTP/1.1 404 Not Found\r\n\r\n"
// 		} else {
// 			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
// 		}
// 	} else {
// 		response = "HTTP/1.1 404 Not Found\r\n\r\n"
// 	}
// 	conn.Write([]byte(response))
// }
// func main() {
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		go handleConnection(conn)
// 	}
// }

// #8. Read request body - Last Base Stage which should support all the before stages

// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// func main() {
// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Println("Logs from your program will appear here!")
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		handleConnection(conn)
// 	}
// }
// func handleConnection(conn net.Conn) {
// 	req := make([]byte, 1024)
// 	conn.Read(req)
// 	if strings.Contains(string(req), "GET / HTTP/1.1") {
// 		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
// 	} else if strings.Contains(string(req), "GET /echo/") {
// 		reqParam := GetStringInBetween(string(req), "GET /echo/", "HTTP/1.1")
// 		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(reqParam)) + "\r\n\r\n" + reqParam))
// 	} else if strings.Contains(string(req), "GET /user-agent HTTP/1.1") {
// 		reqSplitByN := strings.Split(string(req), "\r\n")
// 		for _, eachLine := range reqSplitByN {
// 			if strings.Contains(strings.ToLower(eachLine), "user-agent: ") {
// 				userAgentHeader := eachLine[12:]
// 				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(userAgentHeader)) + "\r\n\r\n" + userAgentHeader))
// 			}
// 		}
// 	} else if strings.Contains(string(req), "GET /files/") {
// 		fileName := GetStringInBetween(string(req), "GET /files/", "HTTP/1.1")
// 		if len(os.Args) < 3 {
// 			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 			conn.Close()
// 		} else {
// 			data, err := os.ReadFile(os.Args[2] + fileName)
// 			if err != nil {
// 				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 				conn.Close()
// 			}
// 			conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(string(data))) + "\r\n\r\n" + string(data)))
// 		}
// 	} else if strings.Contains(string(req), "POST /files/") {
// 		fileName := GetStringInBetween(string(req), "POST /files/", "HTTP/1.1")
// 		if len(os.Args) < 3 {
// 			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 			conn.Close()
// 		} else {
// 			reqSplitByN := strings.Split(string(req), "\r\n")
// 			contentLength := 0
// 			for _, eachLine := range reqSplitByN {
// 				if strings.Contains(eachLine, "Content-Length: ") {
// 					contentLength, _ = strconv.Atoi(eachLine[16:])
// 				}
// 			}
// 			fileDataString := ""
// 			for i := range contentLength {
// 				fileDataString += string(reqSplitByN[len(reqSplitByN)-1][i])
// 			}
// 			fileDataBytes := []byte(fileDataString)
// 			os.WriteFile(os.Args[2]+fileName, fileDataBytes, 0644)
// 			conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
// 			conn.Close()
// 		}
// 	} else {
// 		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 	}
// }
// func GetStringInBetween(str string, start string, end string) (result string) {
// 	s := strings.Index(str, start)
// 	if s == -1 {
// 		return
// 	}
// 	s += len(start)
// 	e := strings.Index(str[s:], end)
// 	if e == -1 {
// 		return
// 	}
// 	e = s + e - 1
// 	return str[s:e]
// }

// --------------------- HTTP Compression ------------------------------

// #1. Compression Headers

// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// func main() {
// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Println("Logs from your program will appear here!")
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		handleConnection(conn)
// 	}
// }
// func handleConnection(conn net.Conn) {
// 	req := make([]byte, 1024)
// 	conn.Read(req)
// 	if strings.Contains(string(req), "GET / HTTP/1.1") {
// 		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
// 	} else if strings.Contains(string(req), "GET /echo/") {
// 		reqParam := GetStringInBetween(string(req), "GET /echo/", "HTTP/1.1")
// 		reqSplitByN := strings.Split(string(req), "\r\n")
// 		for _, eachLine := range reqSplitByN {
// 			if strings.Contains(eachLine, "Accept-Encoding: ") {
// 				acceptEncoding := eachLine[17:]
// 				if acceptEncoding != "gzip" {
// 					acceptEncoding = ""
// 				}
// 				if len(acceptEncoding) == 0 {
// 					conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(reqParam)) + "\r\n\r\n" + reqParam))
// 					conn.Close()
// 				}
// 				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: " + acceptEncoding + "\r\nContent-Length: " + strconv.Itoa(len(reqParam)) + "\r\n\r\n" + reqParam))
// 				conn.Close()
// 			}
// 		}
// 		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(reqParam)) + "\r\n\r\n" + reqParam))
// 	} else if strings.Contains(string(req), "GET /user-agent HTTP/1.1") {
// 		reqSplitByN := strings.Split(string(req), "\r\n")
// 		for _, eachLine := range reqSplitByN {
// 			if strings.Contains(strings.ToLower(eachLine), "user-agent: ") {
// 				userAgentHeader := eachLine[12:]
// 				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(userAgentHeader)) + "\r\n\r\n" + userAgentHeader))
// 			}
// 		}
// 	} else if strings.Contains(string(req), "GET /files/") {
// 		fileName := GetStringInBetween(string(req), "GET /files/", "HTTP/1.1")
// 		if len(os.Args) < 3 {
// 			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 			conn.Close()
// 		} else {
// 			data, err := os.ReadFile(os.Args[2] + fileName)
// 			if err != nil {
// 				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 				conn.Close()
// 			}
// 			conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(string(data))) + "\r\n\r\n" + string(data)))
// 		}
// 	} else if strings.Contains(string(req), "POST /files/") {
// 		fileName := GetStringInBetween(string(req), "POST /files/", "HTTP/1.1")
// 		if len(os.Args) < 3 {
// 			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 			conn.Close()
// 		} else {
// 			reqSplitByN := strings.Split(string(req), "\r\n")
// 			contentLength := 0
// 			for _, eachLine := range reqSplitByN {
// 				if strings.Contains(eachLine, "Content-Length: ") {
// 					contentLength, _ = strconv.Atoi(eachLine[16:])
// 				}
// 			}
// 			fileDataString := ""
// 			for i := range contentLength {
// 				fileDataString += string(reqSplitByN[len(reqSplitByN)-1][i])
// 			}
// 			fileDataBytes := []byte(fileDataString)
// 			os.WriteFile(os.Args[2]+fileName, fileDataBytes, 0644)
// 			conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
// 			conn.Close()
// 		}
// 	} else {
// 		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
// 	}
// }
// func GetStringInBetween(str string, start string, end string) (result string) {
// 	s := strings.Index(str, start)
// 	if s == -1 {
// 		return
// 	}
// 	s += len(start)
// 	e := strings.Index(str[s:], end)
// 	if e == -1 {
// 		return
// 	}
// 	e = s + e - 1
// 	return str[s:e]
// }

// #2. Multiple Compression Schemes

// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"net"
// 	"os"
// 	"regexp"
// 	"strings"
// )

// const CRLF = "\r\n"

// type StatusLine struct {
// 	method  string
// 	path    string
// 	version string
// }
// type Headers = map[string]string
// type Body = string
// type Request struct {
// 	status  StatusLine
// 	headers Headers
// 	body    Body
// }

// func parseRequest(buffer []byte) Request {
// 	buffer = bytes.Trim(buffer, "\x00")
// 	stringBuffer := string(buffer)
// 	return Request{
// 		status:  makeStatusLine(stringBuffer),
// 		headers: makeHeaders(stringBuffer),
// 		body:    makeBody(stringBuffer),
// 	}
// }
// func makeStatusLine(buffer string) StatusLine {
// 	statusLineIndex := strings.Index(buffer, CRLF)
// 	statusLine := buffer[:statusLineIndex]
// 	stringStatusLine := strings.Fields(statusLine)
// 	return StatusLine{
// 		method:  stringStatusLine[0],
// 		path:    stringStatusLine[1],
// 		version: stringStatusLine[2],
// 	}
// }
// func makeHeaders(buffer string) Headers {
// 	headers := make(map[string]string)
// 	lines := strings.Split(buffer, CRLF)
// 	for _, line := range lines {
// 		if line == "" {
// 			break
// 		}
// 		parts := strings.SplitN(line, ":", 2)
// 		if len(parts) == 2 {
// 			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
// 		}
// 	}
// 	return headers
// }
// func makeBody(buffer string) Body {
// 	bodyIndex := strings.LastIndex(buffer, CRLF)
// 	body := buffer[bodyIndex:]
// 	body = body[2:]
// 	return body
// }
// func handleConnection(connection net.Conn) {
// 	buffer := make([]byte, 1024)
// 	connection.Read(buffer)
// 	request := parseRequest(buffer)
// 	homePattern := "/"
// 	echoPattern := regexp.MustCompile(`^/echo/[a-zA-Z0-9]+$`)
// 	userAgentPattern := "/user-agent"
// 	filenamePattern := regexp.MustCompile(`^/files/[a-zA-Z0-9_]+$`)
// 	var response string
// 	if echoPattern.MatchString(request.status.path) {
// 		pathComponents := strings.Split(request.status.path, "/")
// 		path := pathComponents[2]
// 		value, ok := request.headers["Accept-Encoding"]
// 		if ok {
// 			if value == "gzip" {
// 				response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(path), path)
// 			} else if strings.Contains(value, "gzip") {
// 				response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(path), path)
// 			} else {
// 				response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(path), path)
// 			}
// 		} else {
// 			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(path), path)
// 		}
// 	} else if userAgentPattern == request.status.path {
// 		userAgentValue := request.headers["User-Agent"]
// 		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgentValue), userAgentValue)
// 	} else if homePattern == request.status.path {
// 		response = "HTTP/1.1 200 OK\r\n\r\n"
// 	} else if filenamePattern.MatchString(request.status.path) {
// 		pathComponents := strings.Split(request.status.path, "/")
// 		filename := pathComponents[2]
// 		directory := os.Args[2]
// 		if request.status.method == "GET" {
// 			data, err := os.ReadFile(directory + filename)
// 			if err != nil {
// 				fmt.Println("cannot find")
// 				response = "HTTP/1.1 404 Not Found\r\n\r\n"
// 			} else {
// 				fmt.Println("fined")
// 				response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
// 			}
// 		} else {
// 			os.WriteFile(directory+filename, []byte(request.body), os.ModeTemporary)
// 			response = "HTTP/1.1 201 Created\r\n\r\n"
// 		}
// 	} else {
// 		response = "HTTP/1.1 404 Not Found\r\n\r\n"
// 	}
// 	byteResponse := []byte(response)
// 	connection.Write(byteResponse)
// 	connection.Close()
// }
// func main() {
// 	fmt.Println("Server was started")
// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221")
// 		os.Exit(1)
// 	}
// 	for {
// 		connection, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		go handleConnection(connection)
// 	}
// }

// #3. Gzip Compression

package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

var filesDir string

func main() {
	flag.StringVar(&filesDir, "directory", "./files", "specify a directory where files uploaded to server are saved")
	flag.Parse()
	fmt.Println("Saving files to: ", filesDir)
	fmt.Println("Got args: ", os.Args)
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Listening")
	for {
		conn, err := l.Accept()
		fmt.Println("Accepted")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
		fmt.Println("Handling started")
	}
}
func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	request := string(buffer[:n])
	resp := dispatch(parseRequest(request))
	conn.Write([]byte(resp))
	fmt.Println("Handled")
}
func parseRequest(request string) (method string, path string, headers map[string]string, body string) {
	fmt.Println("~~~ request ~~~\n", request)
	requestParts := strings.Split(request, "\r\n")
	requestLine := requestParts[0]
	lineParts := strings.Split(requestLine, " ")
	method = lineParts[0]
	path = lineParts[1]
	body = requestParts[len(requestParts)-1]
	headers = make(map[string]string)
	for i := 1; i < len(requestParts)-1; i++ {
		key, value, _ := strings.Cut(requestParts[i], ":")
		headers[strings.ToLower(strings.Trim(key, " "))] = strings.Trim(value, " ")
	}
	return // named result parameters
}
func dispatch(method string, urlPath string, headers map[string]string, body string) string {
	if method != "GET" {
		if method == "POST" && strings.HasPrefix(urlPath, "/files/") {
			fname := urlPath[len("/files/"):]
			f, err := os.Create(path.Join(filesDir, fname))
			if err != nil {
				return "HTTP/1.1 422 Unprocessable Entity\r\n\r\n"
			}
			defer f.Close()
			f.WriteString(body)
			return "HTTP/1.1 201 Created\r\n\r\n"
		} else {
			return "HTTP/1.1 405 Method Not Allowed\r\n\r\n"
		}
	}
	if urlPath == "/" || urlPath == "/index.html" {
		return "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(urlPath, "/echo/") {
		text := urlPath[len("/echo/"):]
		return contentResponse(text, "text/plain", headers)
	} else if urlPath == "/user-agent" {
		agent := headers["user-agent"]
		return contentResponse(agent, "text/plain", headers)
	} else if strings.HasPrefix(urlPath, "/files/") {
		fname := urlPath[len("/files/"):]
		dat, err := os.ReadFile(path.Join(filesDir, fname))
		if err != nil {
			return "HTTP/1.1 404 Not Found\r\n\r\n"
		}
		return contentResponse(string(dat), "application/octet-stream", headers)
	} else {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
}
func contentResponse(content string, contentType string, requestHeaders map[string]string) string {
	values, ok := requestHeaders["accept-encoding"]
	compression := ""
	if ok {
		for _, elem := range strings.Split(values, ",") {
			if strings.TrimSpace(elem) == "gzip" {
				var buffer bytes.Buffer
				w := gzip.NewWriter(&buffer)
				w.Write([]byte(content))
				w.Close()
				content = buffer.String()
				compression = "Content-Encoding: gzip\r\n"
				break
			}
		}
	}
	return "HTTP/1.1 200 OK\r\n" +
		compression +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + fmt.Sprint(len(content)) +
		"\r\n" +
		"\r\n" + content
}
