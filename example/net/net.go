package net

import "net"


func testNet() {
	ln, err := net.Listen("tcp", ":8080")
if err != nil {
	// handle error
}
for {
	ln.Accept()
	// if err != nil {
	// 	// handle error
	// }
	// go handleConnection(conn)
}
}