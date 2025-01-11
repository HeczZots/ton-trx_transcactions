package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"umbrellaX/network/tron"

	t 
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {
	tronCli := tron.New(true)
	tronCli.Start()
	// tx, err := tronCli.CreateTxTRC20(
	// 	"TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf",
	// 	"TFz6Tt8k1QYb9aTjwh9NaLtuiScmtVW6rC",
	// 	"TBa3KfLYENJX336fZrYdXgx2esaUjHRjAW",
	// 	10000,
	// 	50,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	tx, err := tronCli.CreateTxTRX(
		"TFz6Tt8k1QYb9aTjwh9NaLtuiScmtVW6rC",
		"TBa3KfLYENJX336fZrYdXgx2esaUjHRjAW",
		1000,
		228,
	)
	if err != nil {
		log.Fatal(err)
	}

	res, err := tronCli.SendTx(tx, "1db1510692540cccd2cc6cbfb8e0c53bfd8591770bb4f0c4bf00e283c1234079")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return 
	}

	hash := sha256.Sum256(rawData)
	fmt.Println(hex.EncodeToString(hash[:]))
}

func test() {
	gc := t.NewGrpcClient("")
	if err := gc.Start(grpc.WithInsecure()); err != nil {
		log.Fatal(err)
	}

	_, err := gc.GetAccount("TJGd9GErpVFSuyAhsi5MJ8bGXU6HzaMrWG")
	if err != nil {
		log.Fatal(err)
	}
	tx, err := gc.TriggerConstantContract("", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", "decimals()", "")
	if err != nil {
		log.Fatal(err)
	}
	d, err := gc.TRC20GetDecimals("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	fmt.Println(d.Int64())
	fmt.Println(tx.GetConstantResult()[0])

}
