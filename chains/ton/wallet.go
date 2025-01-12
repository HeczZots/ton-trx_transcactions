package ton

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func (c *Client) AddWallet() error {
	block, err := c.cli.GetMasterchainInfo(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get block info: %w", err)
	}
	addr, err := address.ParseAddr(c.walletAddr)
	if err != nil {
		return err
	}

	account, err := c.cli.GetAccount(context.Background(), block, addr)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	vw := wallet.GetWalletVersion(account)
	w, err := wallet.FromSeed(
		c.cli,
		c.walletSeed,
		vw,
	)
	if err != nil {
		return err
	}
	slog.Info("Wallet initialized", "address: ", w.Address(), " version: ", vw)
	c.w = w
	return nil
}
