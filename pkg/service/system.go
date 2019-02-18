package service

import (
	"context"
	"time"

	"github.com/bluele/rkvs/pkg/proto"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
)

var _ proto.SystemServer = &systemService{}

type systemService struct {
	BaseService

	fwd proto.SystemClient
}

func NewSystemService(r *raft.Raft, fwd proto.SystemClient, logger *logrus.Logger) *systemService {
	return &systemService{
		BaseService: NewBaseService(r, logger),
		fwd:         fwd,
	}
}

func (s *systemService) Join(ctx context.Context, req *proto.SystemRequestJoin) (*proto.SystemResponseJoin, error) {
	if !s.IsLeader() {
		return s.fwd.Join(ctx, req)
	}

	s.Logger.Info("Join")
	err := s.R.AddVoter(
		raft.ServerID(string(req.Id)),
		raft.ServerAddress(string(req.Addr)),
		0,
		time.Second,
	).Error()
	if err != nil {
		return nil, err
	}
	return &proto.SystemResponseJoin{}, nil
}

func (s *systemService) Servers(ctx context.Context, req *proto.SystemRequestServers) (*proto.SystemResponseServers, error) {
	f := s.R.GetConfiguration()
	if err := f.Error(); err != nil {
		return nil, err
	}
	srvs := f.Configuration().Servers
	return &proto.SystemResponseServers{
		Infos: srvsToMessage(srvs),
	}, nil
}

func srvsToMessage(srvs []raft.Server) []*proto.ServerInfo {
	var infos []*proto.ServerInfo
	for _, srv := range srvs {
		infos = append(infos, &proto.ServerInfo{
			Suffrage: int32(srv.Suffrage),
			Id:       string(srv.ID),
			Address:  string(srv.Address),
		})
	}
	return infos
}
