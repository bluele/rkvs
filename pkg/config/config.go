package config

import (
	"path"
	"time"
)

type Config struct {
	Raft   Raft
	Server Server
	DB     DB
}

type Server struct {
	Listen    string
	APIListen string
}

type Raft struct {
	Advertise string
	DataDir   string

	SnapshotInterval  time.Duration
	SnapshotThreshold uint64
}

type DB struct {
	Dir string
}

func MakeDefaultConfig(addr, rootDir string) Config {
	return Config{
		Raft: Raft{
			Advertise:         addr,
			DataDir:           path.Join(rootDir, "raft"),
			SnapshotInterval:  3 * time.Second,
			SnapshotThreshold: 10000,
		},
		Server: Server{
			Listen:    addr,
			APIListen: addr,
		},
		DB: DB{
			Dir: path.Join(rootDir, "db"),
		},
	}
}
