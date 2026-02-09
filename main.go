package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/elielberra/http-server/parser"
	"github.com/elielberra/http-server/socket"
	"golang.org/x/sys/unix"
)

const pageSize = 4096

func handleConnection(netFD int, sa unix.Sockaddr) {
	defer unix.Close(netFD)
	if _, ok := sa.(*unix.SockaddrInet4); ok {
		buffer := make([]byte, pageSize)
		numBytes, err := unix.Read(netFD, buffer)
		if err != nil {
			fmt.Printf("Error reading the contents of the request: %v", err)
			return
		}
		rawRequest := string(buffer[:numBytes])
		var req parser.Request
		if err := parser.SetRequestData(rawRequest, &req); err != nil {
			fmt.Println(err)
		}
		resBody := fmt.Sprintf("Hello from my custom http server!\r\n"+
			"Request method: %s\r\n"+
			"Request path: %s\r\n",
			req.Method, req.Path)
		if req.Method == parser.POST {
			resBody += fmt.Sprintf("Request body: %s\n", req.Body)
		}
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s",
			len(resBody), resBody)
		if _, err := unix.Write(netFD, []byte(res)); err != nil {
			fmt.Printf("Error sending the response to the client:\n%v", err)
		}
	}
}

func validateIp(ip string) error {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return fmt.Errorf("wrong ip format, must be X.X.X.X")
	}
	for _, octet := range octets {
		intOctet, err := strconv.Atoi(octet)
		if err != nil {
			return fmt.Errorf("couldn't convert ip string octet %v to an integer: %v", octet, err)
		}
		minOctet, maxOctet := 0, 255
		if intOctet < minOctet || intOctet > maxOctet {
			return fmt.Errorf("invalid ip octet: %d (must be between %d and %d)", intOctet, minOctet, maxOctet)
		}
	}
	return nil
}

func validatePort(port int) error {
	minPort, maxPort := 0, 65535
	if port < minPort || port > maxPort {
		return fmt.Errorf("port %d is invalid (must be a value between %d and %d)", port, minPort, maxPort)
	}
	return nil
}

func main() {
	fd, err := socket.CreateTCPSocket()
	if err != nil {
		log.Fatalf("Failed to create the socket: %v", err)
	}
	defer unix.Close(fd)
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		log.Fatalf("Failed to set SO_REUSEADDR: %v", err)
	}
	ip := "0.0.0.0"
	port := 8000
	if err := validateIp(ip); err != nil {
		log.Fatalf("Invalid ip %s. Error: %v", ip, err)
	}
	if err := validatePort(port); err != nil {
		log.Fatal(err)
	}
	if err := socket.BindSocket(ip, port, fd); err != nil { // TODO: Remove socket package
		log.Fatalf("Failed to bind the socket: %v", err)
	}
	if err := socket.ListenSocket(fd); err != nil {
		log.Fatalf("Failed to listen on socket %v", err)
	}
	log.Printf("Created a listening server at %s:%d", ip, port)
	for {
		nfd, sa, err := socket.AcceptConn(fd)
		if err != nil {
			log.Printf("Failed to accept the connection: %s", err)
			continue
		}
		go handleConnection(nfd, sa)
	}
}
