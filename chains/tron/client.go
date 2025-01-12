package tron

import (
	"crypto/sha256"
	"encoding/hex"

	t "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	rpcAddr       string
	cli           *t.GrpcClient
	walletSecret  string
	walletAddress string
}

func New(test bool, walletSecret, walletAddress string) *Client {
	rpcAddr := ""
	if test {
		rpcAddr = "grpc.nile.trongrid.io:50051"
	}
	// TODO: сделать мапу адресс кошелька: приватный ключ
	return &Client{
		rpcAddr:       rpcAddr,
		cli:           t.NewGrpcClient(rpcAddr),
		walletSecret:  walletSecret,
		walletAddress: walletAddress,
	}
}

func (c *Client) Start() error {
	return c.cli.Start(grpc.WithInsecure())
}

func (c *Client) Name() string {
	return "tron"
}

func (c *Client) SendTx(currencyContract, to string, amount, feeLimit float64) (hash string, err error) {
	var tx *api.TransactionExtention
	switch currencyContract {
	case "", "TRX", "trx":
		tx, err = c.createTxTRX(to, feeLimit, amount)
		if err != nil {
			return
		}
	default:
		tx, err = c.createTxTRC20(currencyContract, to, feeLimit, amount)
		if err != nil {
			return
		}
	}
	_, err = c.sendTx(tx)
	if err != nil {
		return
	}
	
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(rawData)
	return hex.EncodeToString(sum[:]), nil
}
