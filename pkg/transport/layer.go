package transport

import (
	"net"
	"time"

	"github.com/bluele/rkvs/pkg/proto"
	"github.com/hashicorp/raft"
)

type RaftLayer struct {
	advertise net.Addr
	listener  net.Listener
}

func NewRaftLayer(advertise net.Addr, l net.Listener) *RaftLayer {
	return &RaftLayer{
		advertise: advertise,
		listener:  l,
	}
}

func (t *RaftLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", string(address), timeout)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write([]byte{proto.RaftProto})
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// Accept implements the net.Listener interface.
func (t *RaftLayer) Accept() (c net.Conn, err error) {
	return t.listener.Accept()
}

// Close implements the net.Listener interface.
func (t *RaftLayer) Close() (err error) {
	return t.listener.Close()
}

// Addr implements the net.Listener interface.
func (t *RaftLayer) Addr() net.Addr {
	// Use an advertise addr if provided
	if t.advertise != nil {
		return t.advertise
	}
	return t.listener.Addr()
}
