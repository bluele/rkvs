package service

import (
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
)

type BaseService struct {
	R      *raft.Raft
	Logger *logrus.Logger
}

func NewBaseService(r *raft.Raft, logger *logrus.Logger) BaseService {
	return BaseService{R: r, Logger: logger}
}

func (s *BaseService) IsLeader() bool {
	return s.R.State() == raft.Leader
}
