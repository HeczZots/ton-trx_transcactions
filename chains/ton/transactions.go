package ton

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const tonDecimals = 9
const usdtDecimals = 6

func (c *Client) sendTon(from, to string, amount float64, feeLimit float64) (hash string, err error) {
	// TODO: выбирать кошелек с какого отправлять
	if c.w == nil {
		return "", fmt.Errorf("wallet not initialized")
	}

	toAddr, err := address.ParseAddr(to)
	if err != nil {
		return "", err
	}

	amountTokens, err := tlb.FromDecimal(strconv.FormatFloat(amount, 'f', -1, 64), tonDecimals)
	if err != nil {
		return "", err
	}

	tx, _, err := c.w.TransferWaitTransaction(c.ctx, toAddr, amountTokens, "")
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx.Hash[:]), nil
}

func (c *Client) sendJetton(from, to, tokenAddress string, amount, feeLimit float64) (hash string, err error) {
	// TODO: выбирать кошелек с какого отправлять
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
	// TODO: test
	decimals := c.getDecimals(token)

	amountTokens, err := tlb.FromDecimal(strconv.FormatFloat(amount, 'f', -1, 64), decimals)
	if err != nil {
		return "", err
	}

	// comment, _ := wallet.CreateCommentCell("")
	toAddr, err := address.ParseAddr(to)
	if err != nil {
		return "", err
	}

	transferPayload, err := tokenWallet.BuildTransferPayloadV2(toAddr, toAddr, amountTokens, tlb.ZeroCoins, nil, nil)
	if err != nil {
		return "", err
	}

	// TODO: determine commission or select fee limit
	msg := wallet.SimpleMessage(toAddr, tlb.MustFromTON("0.05"), transferPayload)
	tx, _, err := c.w.SendWaitTransaction(c.ctx, msg)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx.Hash[:]), nil
}

func (c *Client) getDecimals(token *jetton.Client) int {
	decimals := usdtDecimals
	jd, err := token.GetJettonData(c.ctx)
	if err != nil {
		return decimals
	}

	content, ok := jd.Content.(*nft.ContentOnchain)
	if ok {
		v := content.GetAttribute("decimals")
		if f, err := strconv.ParseInt(v, 10, 64); err == nil && f > 0 {
			decimals = int(f)
		}
	}

	return decimals
}
