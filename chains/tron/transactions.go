package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

const trxDecimals = 6

func (c *Client) createTxTRC20(tokenAddr, to string, feeLimit, amount float64) (Tx *api.TransactionExtention, err error) {
	decimals, err := c.cli.TRC20GetDecimals(tokenAddr)
	if err != nil {
		return nil, err
	}

	tx, err := c.cli.TRC20Send(
		c.walletAddress, to,
		tokenAddr,
		big.NewInt(convertAmount(amount, decimals.Int64())), // amount
		convertAmount(amount, trxDecimals))                  // fee
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *Client) createTxTRX( to string, feeLimit, amount float64) (Tx *api.TransactionExtention, err error) {
	tx, err := c.cli.Transfer(
		c.walletAddress, to,
		convertAmount(amount, trxDecimals))
	if err != nil {
		return nil, err
	}
	tx.Transaction.RawData.FeeLimit = convertAmount(feeLimit, trxDecimals)

	return tx, nil
}

func (c *Client) sendTx(tx *api.TransactionExtention) (*api.Return, error) {
	signedTx, err := signTransaction(tx.Transaction, c.walletSecret)
	if err != nil {
		return nil, err
	}

	return c.cli.Broadcast(signedTx)
}

func signTransaction(tx *core.Transaction, privateKeyHex string) (*core.Transaction, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, err
	}

	sk, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	h256h := sha256.New()
	h256h.Write(rawData)

	signature, err := crypto.Sign(h256h.Sum(nil), sk.ToECDSA())
	if err != nil {
		return nil, err
	}

	tx.Signature = append(tx.Signature, signature)

	return tx, nil
}

func convertAmount(value float64, decimals int64) (res int64) {
	return int64(math.Round(value * math.Pow10(int(decimals))))
}
