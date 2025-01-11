package ton

import (
	"context"
	"log"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

type Client struct {
	ctx context.Context
	cli *ton.APIClient
}

func New(test bool) *Client {
	conn := liteclient.NewConnectionPool()
	ctx := context.Background()
	err := conn.AddConnection(ctx, "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		ctx: ctx,
		cli: ton.NewAPIClient(conn),
	}
}
