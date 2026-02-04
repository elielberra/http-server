package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elielberra/http-server/socket"
	"golang.org/x/sys/unix"
)

func handleConnection(netFD int, sa unix.Sockaddr, logger *log.Logger) {
	defer unix.Close(netFD)
	fmt.Printf("Network File Descriptor %v with socket address %v\n", netFD, sa)
	if saIp4, ok := sa.(*unix.SockaddrInet4); ok {
		fmt.Printf("Client ip: %d.%d.%d.%d.\n", saIp4.Addr[0], saIp4.Addr[1], saIp4.Addr[2], saIp4.Addr[3])
		fmt.Printf("Client port: %d\n", saIp4.Port)
		buffer := make([]byte, 4000)
		numBytes, err := unix.Read(netFD, buffer)
		if err != nil {
			logger.Printf("Error reading the contents of the request: %v", err)
		}
		strBuffer := string(buffer[:numBytes])
		log.Print(strBuffer)
		body := fmt.Sprintf("Server received your message:\n%s", strBuffer)
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s", len(body), body)
		if _, err := unix.Write(netFD, []byte(response)); err != nil {
			logger.Printf("Error sending the response to the client:\n%s", err)
		}
	}
}

func validateIp(ip string) error {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return fmt.Errorf("Wrong format, must be X.X.X.X")
	}
	for _, octet := range octets {
		intOctet, err := strconv.Atoi(octet)
		if err != nil {
			return fmt.Errorf("Couldn't convert string octet %v to an integer: %v", octet, err)
		}
		minOctet, maxOctet := 0, 255
		if intOctet < minOctet || intOctet > maxOctet {
			return fmt.Errorf("Invalid octet %d, must be between %d and %d", intOctet, minOctet, maxOctet)
		}
	}
	return nil
}

func validatePort(port int) error {
	minPort, maxPort := 0, 65535
	if port < minPort || port > maxPort {
		return fmt.Errorf("Port %d is invalid. Must be a value between %d and %d", port, minPort, maxPort)
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
	if err := socket.BindSocket(ip, port, fd); err != nil {
		log.Fatalf("Failed to bind the socket: %v", err)
	}
	if err := socket.ListenSocket(fd); err != nil {
		log.Fatalf("Failed to listen on socket %v", err)
	}
	log.Printf("Created a listening server at %s:%d", ip, port)
	logger := log.New(os.Stderr, "", 0)
	for {
		nfd, sa, err := socket.AcceptConn(fd)
		if err != nil {
			log.Printf("Failed to accept the connection: %s", err)
		}
		go handleConnection(nfd, sa, logger)
	}
}
