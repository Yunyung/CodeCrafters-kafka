package main

import (
	"fmt"
	"net"
	"os"
	"unsafe"
)

import . "github.com/codecrafters-io/kafka-starter-go/app/requests"

func main() {
	// Print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Error closing listener: ", err.Error())
			os.Exit(1)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func() {
			for {
				requestHeader, _ := HandleRequest(conn)
				buf := CreateRespond(requestHeader)
				fmt.Printf("Sizeof buf: %d, Len of buf: %d\n", unsafe.Sizeof(buf), buf.Len() )
				fmt.Println("Buffer as bytes:", buf.Bytes())
				conn.Write(buf.Bytes())
			}
		}()
	}
}