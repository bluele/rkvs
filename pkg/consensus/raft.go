package consensus

import (
	"os"
	"path/filepath"
	"time"

	"github.com/bluele/rkvs/pkg/config"
	"github.com/hashicorp/raft"
	raftleveldb "github.com/icexin/raft-leveldb"
)

func NewRaft(cfg config.Raft, fsm raft.FSM, trans raft.Transport, id string) (*raft.Raft, error) {
	raftLogDir := filepath.Join(cfg.DataDir, "log")
	raftMetaDir := filepath.Join(cfg.DataDir, "meta")

	logStore, err := raftleveldb.NewStore(raftLogDir)
	if err != nil {
		return nil, err
	}

	metaStore, err := raftleveldb.NewStore(raftMetaDir)
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(cfg.DataDir, 3, os.Stderr)
	if err != nil {
		return nil, err
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(id)
	raftConfig.SnapshotInterval = time.Duration(cfg.SnapshotInterval)
	raftConfig.SnapshotThreshold = cfg.SnapshotThreshold

	err = raft.ValidateConfig(raftConfig)
	if err != nil {
		return nil, err
	}
	return raft.NewRaft(
		raftConfig,
		fsm,
		logStore,
		metaStore,
		snapshotStore,
		trans,
	)
}
