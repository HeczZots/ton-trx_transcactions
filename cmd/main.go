package main

import (
	"log"
	"log/slog"
	"os"
	"umbrellaX/network/ton"
)

func main() {
	cli := ton.New(false)
	slog.Info("Adding ton client")

	if err := cli.Start(os.Getenv("TON_SEED"), 42); err != nil {
		log.Fatal(err)
	}
	slog.Info("Sending tx")
	// usdt "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"
	hash, err := cli.SendTon("", "UQDar3bJMt1cFtVwFxEkb9HYATr-vdDa5XI6BCxExkzuEnne", 0.01, 125)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(hash)
}
