package tron

import (
	t "github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
)

type Client struct {
	rpcAddr string
	cli     *t.GrpcClient
}

func New(test bool) *Client {
	rpcAddr := ""
	if test {
		rpcAddr = "grpc.nile.trongrid.io:50051"
	}
	
	return &Client{
		rpcAddr: rpcAddr,
		cli:     t.NewGrpcClient(rpcAddr),
	}
}

func (c *Client) Start() error {
	return c.cli.Start(grpc.WithInsecure())
}
