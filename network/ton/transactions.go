package ton

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func (c *Client) SendTx(from, to string, amount float64, feeLimit float64) error {
	return nil
}

func (c *Client) SendJetton(from, to, tokenAddress string, amount, feeLimit float64) (hash string, err error) {
	if c.w == nil {
		return "", fmt.Errorf("wallet not initialized")
	}
	ca, err := address.ParseAddr(tokenAddress)
	if err != nil {
		return "", err
	}

	token := jetton.NewJettonMasterClient(c.cli, ca)

	tokenWallet, err := token.GetJettonWallet(c.ctx, c.w.WalletAddress())
	if err != nil {
		return "", err
	}
	// TODO: get decimals
	// jd, err := token.GetJettonData(c.ctx)
	// if err != nil {
	// 	return "", err
	// }

	toAddr, err := address.ParseAddr(to)
	if err != nil {
		return "", err
	}
	// decimals hardcode for now
	amountTokens := tlb.MustFromDecimal(strconv.FormatFloat(amount, 'f', -1, 64), 9)
	comment, _ := wallet.CreateCommentCell("")

	transferPayload, err := tokenWallet.BuildTransferPayloadV2(toAddr, toAddr, amountTokens, tlb.ZeroCoins, comment, nil)
	if err != nil {
		return "", err
	}

	msg := wallet.SimpleMessageAutoBounce(toAddr, tlb.ZeroCoins, transferPayload)
	tx, _, err := c.w.SendWaitTransaction(c.ctx, msg)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx.Hash[:]), nil
}
