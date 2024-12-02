package main

import (
	"fmt"
	"net"
	"os"
	"encoding/binary"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var BUFFER_SIZE int = 1024
const MESSAGE_SIZE int = 4 // bytes for message size

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	
	buffer := make([]byte, 1024)

	conn.Read(buffer)
	request_api_version := int16(binary.BigEndian.Uint16(buffer[6:8]))
	var error_code uint16 = 0
	if request_api_version < 0 || request_api_version > 4 {
		error_code = 35
	}
	fmt.Printf("Received request for API version (%d)\n", request_api_version)
	fmt.Printf("Received message %v (%d)", buffer[8:12], int32(binary.BigEndian.Uint32(buffer[8:12])))
	respond := make([]byte, 24)
	copy(respond, make([]byte, MESSAGE_SIZE)) // copy the first 4 bytes from the request[MESSAGE_SIZE]byte)
	copy(respond[4:8], buffer[8:12]) // correlation id
	//copy(respond[8:12], error_code)
	binary.BigEndian.PutUint16(respond[8:10], error_code)
	conn.Write(respond)

	conn.Close()
}
