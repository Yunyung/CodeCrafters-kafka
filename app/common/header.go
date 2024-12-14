package common

type MESSAGESIZE = int32
type TaggedFields = byte

type RequestHeader struct {
	MessageSize MESSAGESIZE
	RequestApiKey int16
	RequestApiVersion int16
	CorrelationId int32
	ClientID string
	TaggedFields TaggedFields
}

type RespondHeader struct {
	MessageSize MESSAGESIZE
	CorrelationId int32
}