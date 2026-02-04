package send

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	MAXMSGSIZE = 8000
)

func send() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("./client [IPv4] [Port]")
	}

	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_IP)
	if err != nil {
		log.Fatal("Socket: ", err)
	}

	port, err := strconv.Atoi(args[1])
	if err != nil || (port < 0 || port > 100000) {
		os.Stderr.WriteString("Invalid port format\n")
		return
	}
	serverAddr := &unix.SockaddrInet4{
		Port: port,
		Addr: inetAddr(args[0]),
	}

	err = unix.Connect(serverFD, serverAddr)
	if err != nil {
		if err == unix.ECONNREFUSED {
			fmt.Println("* Connection failed")
			unix.Close(serverFD)
			return
		}
	}

	var msg string
	var response []byte

	response = make([]byte, MAXMSGSIZE)

	print("> ")
	fmt.Scanln(&msg)
	err = unix.Sendmsg(
		serverFD,
		[]byte(msg),
		nil, serverAddr, unix.MSG_DONTWAIT)
	if err != nil {
		fmt.Println("Sendmsg: ", err)
	}
	_, _, err = unix.Recvfrom(serverFD, response, 0)
	if err != nil {
		fmt.Println("Recvfrom: ", err)
		unix.Close(serverFD)
		return
	}
	fmt.Printf("< %s\n", string(response))
	unix.Close(serverFD)
}

func inetAddr(ipaddr string) [4]byte {
	var (
		ip                 = strings.Split(ipaddr, ".")
		ip1, ip2, ip3, ip4 uint64
	)
	ip1, _ = strconv.ParseUint(ip[0], 10, 8)
	ip2, _ = strconv.ParseUint(ip[1], 10, 8)
	ip3, _ = strconv.ParseUint(ip[2], 10, 8)
	ip4, _ = strconv.ParseUint(ip[3], 10, 8)
	return [4]byte{byte(ip1), byte(ip2), byte(ip3), byte(ip4)}
}