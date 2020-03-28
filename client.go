// A cache client used to test the LRU cache library.
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	// connect to server
	conn, _ := net.Dial("tcp", "127.0.0.1:9000")
	text := "set name adnan"
	// send to server
	fmt.Fprintf(conn, text+"\n")
	// wait for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)

	text = "get name"
	// send to server
	fmt.Fprintf(conn, text+"\n")
	// wait for reply
	message, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}
