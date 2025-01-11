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
	
	if err := cli.Start(os.Getenv("TON_SEED")); err != nil {
		log.Fatal(err)
	}

	
}
