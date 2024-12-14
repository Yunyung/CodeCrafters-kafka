package requests

import (
	"net"
	"encoding/binary"
	"fmt"
	"github.com/codecrafters-io/kafka-starter-go/app/common"
)

/* No struct required for ApiVersions request < 3 */

func HandleRequest(conn net.Conn) (common.RequestHeader, error) {
	request := common.RequestHeader{}
	err := binary.Read(conn, binary.BigEndian, &request.MessageSize)
	if err != nil {
		fmt.Println("Error reading from connection!!!: ", err.Error())
		return common.RequestHeader{}, err
	}
	binary.Read(conn, binary.BigEndian, &request.RequestApiKey)
	binary.Read(conn, binary.BigEndian, &request.RequestApiVersion)
	binary.Read(conn, binary.BigEndian, &request.CorrelationId)
	fmt.Printf("Received request for message size (%d)\n", request.MessageSize)
	fmt.Printf("Received request for API key (%d)\n", request.RequestApiKey)
	fmt.Printf("Received request for API version (%d)\n", request.RequestApiVersion)
	fmt.Printf("Received correlation id (%d)\n", request.CorrelationId)
	
	conn.Read(make([]byte, 2048))

	return request, nil
}