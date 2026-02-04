package socket

import (
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func ipToBytes(ip string) [4]byte {
	var ipBytes [4]byte
	octets := strings.Split(ip, ".")
	for i, octet := range octets {
		octetInt, _ := strconv.Atoi(octet)
		ipBytes[i] = byte(octetInt)
	}
	return ipBytes
}

func CreateTCPSocket() (int, error) {
	return unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
}

func BindSocket(ip string, port int, fd int) (err error) {
	addr := ipToBytes(ip)
	sockAddr := &unix.SockaddrInet4{Port: port, Addr: addr}
	return unix.Bind(fd, sockAddr)
}

func ListenSocket(fd int) (err error) {
	return unix.Listen(fd, int(unix.SOMAXCONN))
}

func AcceptConn(fd int) (nfd int, sa unix.Sockaddr, err error) {
	return unix.Accept(fd)
}
