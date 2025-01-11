package main

import (
	"log"
	"log/slog"
	"os"
	"umbrellaX/network/ton"
	"umbrellaX/network/tron"
	"umbrellaX/server"
)

func main() {
	tonCli := ton.New(false)
	slog.Info("Adding ton client")

	if err := tonCli.Start(os.Getenv("TON_SEED"), 42); err != nil {
		log.Fatal(err)
	}

	tronCli := tron.New(false)

	if err := tronCli.Start(); err != nil {
		log.Fatal(err)
	}
	
	log.Fatal(server.New(tonCli, tronCli).Start("8080"))
}
