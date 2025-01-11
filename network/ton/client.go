package ton

import (
	"context"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type Client struct {
	conn    *liteclient.ConnectionPool
	ctx     context.Context
	cli     *ton.APIClient
	w       *wallet.Wallet
	jettons map[string]string
}

func New(test bool) *Client {
	conn := liteclient.NewConnectionPool()
	ctx := context.Background()

	return &Client{
		jettons: make(map[string]string),
		ctx:     ctx,
		conn:    conn,
	}
}

func (c *Client) Start(seed string) error {
	ctx := context.Background()
	if err := c.conn.AddConnectionsFromConfigUrl(ctx, "https://ton.org/global.config.json"); err != nil {
		return err
	}

	c.cli = ton.NewAPIClient(c.conn)
	// TODO: добавлять множество кошельков
	if err := c.AddWallet(seed); err != nil {
		return err
	}
	// deprecated
	// if err := c.getJettonAddress(usdtContract); err != nil {
	// 	return err
	// }

	return nil
}

var usdtContract = "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"

func (c *Client) AddWallet(recoveryPhrase string) error {
	rp := strings.Split(recoveryPhrase, " ")
	w, err := wallet.FromSeed(c.cli, rp, wallet.V4R2)
	if err != nil {
		return err
	}

	c.w = w
	return nil
}

// deprecated
func (c *Client) getJettonAddress(tokenAddress string) error {
	b, err := c.cli.GetMasterchainInfo(c.ctx)
	if err != nil {
		return err
	}

	ca, err := address.ParseAddr(tokenAddress)
	if err != nil {
		return err
	}

	// get jetton address for transfering other tokens
	payload := cell.BeginCell().MustStoreUInt(0x2c76b973, 32).MustStoreAddr(c.w.Address()).EndCell()

	res, err := c.cli.RunGetMethod(c.ctx, b, ca, "get_wallet_address", payload.BeginParse())
	if err != nil {
		return err
	}

	c.jettons[tokenAddress] = res.MustSlice(0).MustLoadAddr().String()

	return nil
}
