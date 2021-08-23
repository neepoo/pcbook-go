package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtobufToJSON(message proto.Message) ([]byte, error) {
	m := protojson.MarshalOptions{
		Indent:          "	",
		UseEnumNumbers:  false,
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	return m.Marshal(message)
}
