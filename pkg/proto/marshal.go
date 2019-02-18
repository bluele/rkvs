package proto

import (
	"errors"
	fmt "fmt"

	pb "github.com/golang/protobuf/proto"
)

const (
	RaftProto = 0

	OPWrite byte = 0x00
)

func Marshal(msg pb.Message) ([]byte, error) {
	var op byte
	switch msg := msg.(type) {
	case *KVSRequestWrite:
		op = OPWrite
	default:
		return nil, fmt.Errorf("unknown type %T", msg)
	}
	buf := pb.NewBuffer(nil)
	buf.SetBuf([]byte{op})
	if err := buf.Marshal(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(buf []byte) (pb.Message, error) {
	if len(buf) == 0 {
		return nil, errors.New("buf is nil")
	}
	switch buf[0] {
	case OPWrite:
		msg := new(KVSRequestWrite)
		return msg, pb.Unmarshal(buf[1:], msg)
	default:
		return nil, fmt.Errorf("unknown op %v", buf[0])
	}
}
