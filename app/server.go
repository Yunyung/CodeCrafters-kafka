package main

import (
	"fmt"
	"net"
	"os"
	"encoding/binary"
	"bytes"
	"unsafe"
)

type RequestHeader struct {
	messageSize int32
	requestApiKey int16
	requestApiVersion int16
	correlationId int32
	clientID string
	taggedFields string
}

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var BUFFER_SIZE int = 1024
const MESSAGE_SIZE int = 4 // bytes for message size

type TaggedFields struct {
	tagBuffer byte
}
type RespondHeader struct {
	messageSize int32
	correlationId int32
}

type ApiKey struct	{
	apiKey int16
	minVersion int16
	maxVersion int16
}
type Respond struct {
	respondHeader RespondHeader
	errorCode int16
	numApiKeys byte
	apiKeys ApiKey
	taggedFields TaggedFields
	throttleTimeMs int32
	taggedFields2 TaggedFields
}

func handleRequest(conn net.Conn) (RequestHeader, error) {
	request := RequestHeader{}
	err := binary.Read(conn, binary.BigEndian, &request.messageSize)
	if err != nil {
		fmt.Println("Error reading from connection!!!: ", err.Error())
		return RequestHeader{}, err
	}
	binary.Read(conn, binary.BigEndian, &request.requestApiKey)
	binary.Read(conn, binary.BigEndian, &request.requestApiVersion)
	binary.Read(conn, binary.BigEndian, &request.correlationId)
	fmt.Printf("Received request for message size (%d)\n", request.messageSize)
	fmt.Printf("Received request for API key (%d)\n", request.requestApiKey)
	fmt.Printf("Received request for API version (%d)\n", request.requestApiVersion)
	fmt.Printf("Received correlation id (%d)\n", request.correlationId)
	
	conn.Read(make([]byte, 2048))

	return request, nil
}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error closing listener: ", err.Error())
			os.Exit(1)
		}
	}(l)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func() {
		
			for {
				var error_code int16 = 0
				requestHeader, _ := handleRequest(conn)
				if requestHeader.requestApiVersion < 0 || requestHeader.requestApiVersion > 4 {
					error_code = 35
				}

				fmt.Printf("RespondHeader len (%d)\n", unsafe.Sizeof(RespondHeader{}))
				fmt.Printf("ApiKey len (%d)\n", unsafe.Sizeof(ApiKey{}))
				fmt.Printf("TaggedFields len (%d)\n", unsafe.Sizeof(TaggedFields{}))
				fmt.Printf("respond len (%d)\n", unsafe.Sizeof(Respond{}))


				var respond Respond
				respond.respondHeader.correlationId = requestHeader.correlationId
				respond.errorCode = error_code
				respond.numApiKeys = byte(2)
				respond.apiKeys.apiKey = 18
				respond.apiKeys.minVersion = 0
				respond.apiKeys.maxVersion = 4
				
				buf := &bytes.Buffer{}
				
				// binary.Write(buf, binary.BigEndian, respond.respondHeader)
				// binary.Write(buf, binary.BigEndian, respond.errorCode)
				// binary.Write(buf, binary.BigEndian, respond.numApiKeys)
				// binary.Write(buf, binary.BigEndian, respond.apiKeys)
				// binary.Write(buf, binary.BigEndian, respond.taggedFields)
				// binary.Write(buf, binary.BigEndian, respond.throttleTimeMs)
				// binary.Write(buf, binary.BigEndian, respond.taggedFields)
				// // Correctly calculate message size after writing fields.
				// respond.respondHeader.messageSize = int32(buf.Len())
				// binary.BigEndian.PutUint32(buf.Bytes()[:4], uint32(respond.respondHeader.messageSize) - 4)
				binary.Write(buf, binary.BigEndian, respond)
				respond.respondHeader.messageSize = int32(buf.Len())
				binary.BigEndian.PutUint32(buf.Bytes()[:4], uint32(respond.respondHeader.messageSize) - 4)


				fmt.Printf("Afer: Sizeof respond: %d, Sizeof buf: %d, Len of buf: %d\n", unsafe.Sizeof(respond), unsafe.Sizeof(buf), buf.Len() )
				fmt.Println("Buffer as bytes:", buf.Bytes())

				conn.Write(buf.Bytes())
			}
		}()
	}
}