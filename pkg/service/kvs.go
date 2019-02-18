package service

import (
	"context"
	"time"

	"github.com/bluele/rkvs/pkg/consensus"
	"github.com/bluele/rkvs/pkg/proto"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
)

var _ proto.KVSServer = &kvsService{}

type kvsService struct {
	BaseService

	fsm *consensus.FSM
	fwd proto.KVSClient
}

func NewKVSService(r *raft.Raft, fsm *consensus.FSM, fwd proto.KVSClient, logger *logrus.Logger) *kvsService {
	return &kvsService{
		BaseService: NewBaseService(r, logger),
		fsm:         fsm,
		fwd:         fwd,
	}
}

func (s *kvsService) Read(ctx context.Context, req *proto.KVSRequestRead) (*proto.KVSResponseRead, error) {
	s.Logger.Info("Read")
	if !s.IsLeader() {
		return s.fwd.Read(ctx, req)
	}
	return s.read(ctx, req)
}

func (s *kvsService) read(ctx context.Context, req *proto.KVSRequestRead) (*proto.KVSResponseRead, error) {
	v, err := s.fsm.Get(req.Key)
	if err != nil {
		return nil, err
	}
	return &proto.KVSResponseRead{Value: v}, nil
}

func (s *kvsService) Write(ctx context.Context, req *proto.KVSRequestWrite) (*proto.KVSResponseWrite, error) {
	s.Logger.Info("Write")
	if !s.IsLeader() {
		return s.fwd.Write(ctx, req)
	}
	return s.write(ctx, req)
}

func (s *kvsService) write(ctx context.Context, req *proto.KVSRequestWrite) (*proto.KVSResponseWrite, error) {
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret := s.R.Apply(buf, time.Second)

	if err := ret.Error(); err != nil {
		return nil, err
	}
	if ret.Response() != nil {
		return nil, ret.Response().(error)
	}

	return &proto.KVSResponseWrite{}, nil
}

func (s *kvsService) Ping(ctx context.Context, req *proto.KVSRequestPing) (*proto.KVSResponsePing, error) {
	s.Logger.Info("Ping")
	if !s.IsLeader() {
		return s.fwd.Ping(ctx, req)
	}
	s.Logger.Info("I am Leader")
	return &proto.KVSResponsePing{}, nil
}
