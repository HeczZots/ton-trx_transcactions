package network

type Networker interface{
	SendTx(from , to string, amount float64)
}

