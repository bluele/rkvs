package node

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/bluele/rkvs/pkg/client"
	"github.com/bluele/rkvs/pkg/config"
	"github.com/bluele/rkvs/pkg/consensus"
	"github.com/bluele/rkvs/pkg/proto"
	"github.com/bluele/rkvs/pkg/service"
	"github.com/bluele/rkvs/pkg/transport"
	"github.com/soheilhy/cmux"

	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Node struct {
	logger *logrus.Logger

	r     *raft.Raft
	trans *raft.NetworkTransport
	layer *transport.RaftLayer

	mux *proto.Mux
	fsm *consensus.FSM

	Closers []io.Closer
}

func MakeNode(cfg config.Config, id string, logger *logrus.Logger) (*Node, error) {
	// setup listener
	logger.Info("listen:", cfg.Server.Listen)

	l, err := net.Listen("tcp", cfg.Server.Listen)
	if err != nil {
		return nil, err
	}

	mux := proto.NewMux(l, nil)
	raftl := mux.Handle(proto.RaftProto)
	gl := mux.HandleThird(cmux.HTTP2())

	// setup raft transporter
	advertise, err := net.ResolveTCPAddr("tcp", cfg.Raft.Advertise)
	if err != nil {
		return nil, err
	}
	layer := transport.NewRaftLayer(advertise, raftl)
	trans := raft.NewNetworkTransport(
		layer,
		5,
		time.Second,
		os.Stderr,
	)

	fsm, err := consensus.NewFSM(cfg.DB)
	if err != nil {
		return nil, err
	}

	r, err := consensus.NewRaft(cfg.Raft, fsm, trans, id)
	if err != nil {
		return nil, err
	}

	srv := grpc.NewServer()
	cp := client.NewLeaderConnector(r, client.NewConnectionPool(), grpc.WithInsecure())

	proto.RegisterKVSServer(srv, service.NewKVSService(r, fsm, client.NewKVSClient(cp), logger))
	proto.RegisterSystemServer(srv, service.NewSystemService(r, client.NewSystemClient(cp), logger))

	go srv.Serve(gl)

	return &Node{
		logger: logger,
		r:      r,
		mux:    mux,
		layer:  layer,
		trans:  trans,
		fsm:    fsm,
	}, nil
}

func (n *Node) Serve() error {
	return n.mux.Serve()
}

func (n *Node) Raft() *raft.Raft {
	return n.r
}

func (n *Node) Close() error {
	n.layer.Close()
	n.trans.Close()
	ret := n.r.Shutdown()
	// wait raft shutdown
	ret.Error()
	n.fsm.Close()
	return n.mux.Close()
}

func (n *Node) Bootstrap(id, addr string) {
	var configuration raft.Configuration
	configuration.Servers = append(configuration.Servers, raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(id),
		Address:  raft.ServerAddress(addr),
	})
	if err := n.r.BootstrapCluster(configuration).Error(); err != nil {
		panic(err)
	}
}

func (n *Node) AddNode(id, addr string) error {
	return n.r.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr), 0, 0).Error()
}
