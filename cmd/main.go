package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	initilizeLogging()
}

func initilizeLogging() {
	debugFlag := os.Getenv("DEBUG")
	if debugFlag == "" {
		log.Fatalf("Unable to get \"DEBUG\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
	if debugFlag != "0" && debugFlag != "1" {
		log.Fatalf("\"DEBUG\" environmnet variable must be values 0 (debug mode off) or 1 (debug mode on).")
	}

	var slogOpts *slog.HandlerOptions
	switch debugFlag {
	case "0":
		slogOpts = nil
	case "1":
		slogOpts = &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
	slog.SetDefault(logger)
}
