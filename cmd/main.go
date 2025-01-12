package main

import (
	"log"
	"log/slog"
	"os"
	"umbrellaX/chains/ton"
	"umbrellaX/chains/tron"
	"umbrellaX/server"
)

func main() {
	tonCli := ton.New(false, os.Getenv("TON_SEED"), "UQD6RpM5JZcBwCg7zINlmd2JQToStegSzoJxLb7g7utIcq0d")
	slog.Info("Adding ton client")

	if err := tonCli.Start(); err != nil {
		log.Fatal(err)
	}

	slog.Info("Adding tron client")

	tronCli := tron.New(false, os.Getenv("TRON_SECRET"), "TFz6Tt8k1QYb9aTjwh9NaLtuiScmtVW6rC")

	if err := tronCli.Start(); err != nil {
		log.Fatal(err)
	}

	port := "8080"
	slog.Info("Starting server on ", "port", port)

	log.Fatal(server.New(tonCli, tronCli).Start(port))
}
