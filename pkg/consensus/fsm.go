package consensus

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/bluele/rkvs/pkg/config"
	"github.com/bluele/rkvs/pkg/proto"
	"github.com/bluele/rkvs/pkg/util"
	"github.com/hashicorp/raft"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	errBadMethod = errors.New("bad method")
	errBadAction = errors.New("bad action")
)

type FSM struct {
	cfg config.DB
	db  *leveldb.DB
}

func NewFSM(cfg config.DB) (*FSM, error) {
	db, err := leveldb.OpenFile(cfg.Dir, nil)
	if err != nil {
		return nil, err
	}

	return &FSM{
		cfg: cfg,
		db:  db,
	}, nil
}

func (f *FSM) Get(key []byte) ([]byte, error) {
	return f.db.Get(key, nil)
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	msg, err := proto.Unmarshal(l.Data)
	if err != nil {
		return err
	}
	switch msg := msg.(type) {
	case *proto.KVSRequestWrite:
		return f.db.Put(msg.Key, msg.Value, nil)
	default:
		return errBadAction
	}
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	snapshot, err := f.db.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return &fsmSnapshot{snapshot}, nil
}

func (f *FSM) Restore(r io.ReadCloser) error {
	defer r.Close()

	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	err = f.db.Close()
	if err != nil {
		return err
	}

	oldname := f.cfg.Dir + ".old"
	err = os.Rename(f.cfg.Dir, oldname)
	if err != nil {
		return err
	}
	defer os.RemoveAll(oldname)

	err = util.Untar(f.cfg.Dir, zr)
	if err != nil {
		return err
	}

	db, err := leveldb.OpenFile(f.cfg.Dir, nil)
	if err != nil {
		return err
	}

	f.db = db
	return nil
}

func (f *FSM) Close() error {
	return f.db.Close()
}

// fsmSnapshot implement FSMSnapshot interface
type fsmSnapshot struct {
	snapshot *leveldb.Snapshot
}

// First, walk all kvs, write temp leveldb.
// Second, make tar.gz for temp leveldb dir
func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	// Create a temporary path for the state store
	tmpPath, err := ioutil.TempDir(os.TempDir(), "state")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpPath)

	db, err := leveldb.OpenFile(tmpPath, nil)
	if err != nil {
		return err
	}
	iter := f.snapshot.NewIterator(nil, nil)
	for iter.Next() {
		err = db.Put(iter.Key(), iter.Value(), nil)
		if err != nil {
			db.Close()
			sink.Cancel()
			return err
		}
	}
	iter.Release()
	db.Close()

	// make tar.gz
	w := gzip.NewWriter(sink)
	err = util.Tar(tmpPath, w)
	if err != nil {
		sink.Cancel()
		return err
	}

	err = w.Close()
	if err != nil {
		sink.Cancel()
		return err
	}

	sink.Close()
	return nil
}

func (f *fsmSnapshot) Release() {
	f.snapshot.Release()
}
