package client

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
)

type Connector interface {
	GetConn(ctx context.Context) (*grpc.ClientConn, error)
}

type CommonConnector struct {
	opts []grpc.DialOption
	addr string
}

func (c *CommonConnector) GetConn(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, c.addr, grpc.WithInsecure())
}

func NewCommonConnector(addr string, opts ...grpc.DialOption) *CommonConnector {
	return &CommonConnector{addr: addr, opts: opts}
}

type PoolingConnector struct {
	opts  []grpc.DialOption
	addrs []string
	cp    *ConnectionPool
}

func NewPoolingConnector(cp *ConnectionPool, addrs []string, opts ...grpc.DialOption) *PoolingConnector {
	return &PoolingConnector{cp: cp, addrs: addrs, opts: opts}
}

func (c *PoolingConnector) GetConn(ctx context.Context) (*grpc.ClientConn, error) {
	addr := c.addrs[rand.Intn(len(c.addrs))]
	log.Println("GetConn:", addr)
	return c.cp.Get(ctx, addr, c.opts...)
}

type LeaderConnector struct {
	opts []grpc.DialOption
	r    *raft.Raft
	cp   *ConnectionPool
}

func NewLeaderConnector(r *raft.Raft, cp *ConnectionPool, opts ...grpc.DialOption) *LeaderConnector {
	return &LeaderConnector{r: r, cp: cp, opts: opts}
}

func (c *LeaderConnector) lookup() (string, error) {
	leader := c.r.Leader()
	if leader == "" {
		return "", errors.New("no leader")
	}
	return string(leader), nil
}

func (c *LeaderConnector) GetConn(ctx context.Context) (*grpc.ClientConn, error) {
	addr, err := c.lookup()
	if err != nil {
		return nil, err
	}
	return c.cp.Get(ctx, addr, c.opts...)
}

type ConnectionPool struct {
	mu    sync.RWMutex
	conns map[string]*grpc.ClientConn
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{conns: make(map[string]*grpc.ClientConn)}
}

func (cp *ConnectionPool) Get(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	cp.mu.RLock()
	conn, ok := cp.conns[target]
	cp.mu.RUnlock()
	if ok {
		return conn, nil
	}

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	cp.mu.Lock()
	cp.conns[target] = conn
	cp.mu.Unlock()

	return conn, nil
}

func (cp *ConnectionPool) Delete(target string) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	_, ok := cp.conns[target]
	if !ok {
		return
	}
	delete(cp.conns, target)
}
