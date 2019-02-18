package client

import (
	"context"
	"fmt"

	"github.com/bluele/rkvs/pkg/proto"
	"google.golang.org/grpc"
)

var _ proto.KVSClient = &kvsClient{}

type kvsClient struct {
	connector Connector
}

func NewKVSClient(c Connector) *kvsClient {
	return &kvsClient{connector: c}
}

func (c *kvsClient) getClient(ctx context.Context) (proto.KVSClient, error) {
	conn, err := c.connector.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	return proto.NewKVSClient(conn), nil
}

func (c *kvsClient) Read(ctx context.Context, in *proto.KVSRequestRead, opts ...grpc.CallOption) (*proto.KVSResponseRead, error) {
	cl, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	return cl.Read(ctx, in, opts...)
}

func (c *kvsClient) Write(ctx context.Context, in *proto.KVSRequestWrite, opts ...grpc.CallOption) (*proto.KVSResponseWrite, error) {
	cl, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	return cl.Write(ctx, in, opts...)
}

func (c *kvsClient) Ping(ctx context.Context, in *proto.KVSRequestPing, opts ...grpc.CallOption) (*proto.KVSResponsePing, error) {
	cl, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("got client", cl)
	return cl.Ping(ctx, in, opts...)
}
