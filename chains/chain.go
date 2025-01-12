package chain

type Chain interface {
	SendTx(currencyContract, to string, amount, feeLimit float64) (hash string, err error)
	Name() string
}
