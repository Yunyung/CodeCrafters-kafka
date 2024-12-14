package requests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/codecrafters-io/kafka-starter-go/app/common"
	"unsafe"
)

type ApiVersionsResponse struct {
	errorCode int16
	numApiKeys byte
	apiKeys []ApiKey
	throttleTimeMs int32
	taggedFields common.TaggedFields
}

type ApiKey struct	{
	apiKey int16
	minVersion int16
	maxVersion int16
	taggedFields common.TaggedFields
}

func CreateRespond(requestHeader common.RequestHeader) *bytes.Buffer {
	var error_code int16 = 0
	if requestHeader.RequestApiVersion < 0 || requestHeader.RequestApiVersion > 4 {
		error_code = 35
	}
	fmt.Printf("RespondHeader len (%d)\n", unsafe.Sizeof(common.RespondHeader{}))
	fmt.Printf("ApiKey len (%d)\n", unsafe.Sizeof(ApiKey{}))
	fmt.Printf("ApiVersionsResponse len (%d)\n", unsafe.Sizeof(ApiVersionsResponse{}))

	numApiKeys := 1
	var apiVersionsRespond ApiVersionsResponse
	var respondHeader common.RespondHeader
	respondHeader.CorrelationId = requestHeader.CorrelationId
	apiVersionsRespond.errorCode = error_code
    apiVersionsRespond.numApiKeys = byte(numApiKeys)
    apiVersionsRespond.apiKeys = make([]ApiKey, numApiKeys)
	// APIVersions
	apiVersionsRespond.apiKeys[0] = ApiKey{apiKey: 18, minVersion: 0, maxVersion: 4}

	// Wrtie to bytes buffer
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, respondHeader)
	binary.Write(buf, binary.BigEndian, apiVersionsRespond.errorCode)
	binary.Write(buf, binary.BigEndian, apiVersionsRespond.numApiKeys)
	// Write each API key individually
	for _, apiKey := range apiVersionsRespond.apiKeys {
		binary.Write(buf, binary.BigEndian, apiKey)
	}
	//binary.Write(buf, binary.BigEndian, &apiVersionsRespond)
	messageBytes := unsafe.Sizeof(common.MESSAGESIZE(0))
	binary.BigEndian.PutUint32(buf.Bytes()[:4], uint32(buf.Len()) - uint32(messageBytes))
	fmt.Printf("Sizeof buf: %d, Len of buf: %d\n", unsafe.Sizeof(buf), buf.Len() )
	fmt.Println("Buffer as bytes:", buf.Bytes())
	return buf
}