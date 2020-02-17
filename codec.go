package grpcdynamic

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const CodecName = "grpcdynamic"

func init() {
	encoding.RegisterCodec(&codec{})
}

type codec struct{}

func (c *codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c *codec) Name() string {
	return CodecName
}
