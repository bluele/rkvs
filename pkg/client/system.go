package client

import (
	"context"

	"github.com/bluele/rkvs/pkg/proto"
	"google.golang.org/grpc"
)

var _ proto.SystemClient

type systemClient struct {
	connector Connector
}

func NewSystemClient(c Connector) *systemClient {
	return &systemClient{connector: c}
}

func (c *systemClient) getClient(ctx context.Context) (proto.SystemClient, error) {
	conn, err := c.connector.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	return proto.NewSystemClient(conn), nil
}

func (c *systemClient) Join(ctx context.Context, in *proto.SystemRequestJoin, opts ...grpc.CallOption) (*proto.SystemResponseJoin, error) {
	cl, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	return cl.Join(ctx, in, opts...)
}

func (c *systemClient) Servers(ctx context.Context, in *proto.SystemRequestServers, opts ...grpc.CallOption) (*proto.SystemResponseServers, error) {
	cl, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	return cl.Servers(ctx, in, opts...)
}
