package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/NerdBow/Grinders-API/internal/server"
)

func main() {
	// If no fatal is thrown then all env vars are not empty.
	checkSQLiteEnv()
	checkArgonEnv()
	checkHTTPEnv()
	initilizeLogging()

	server.Run()
}

func checkSQLiteEnv() {
	dbfile := os.Getenv("DBFILE")
	if dbfile == "" {
		log.Fatalf("Unable to get \"DBFILE\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		log.Fatalf("File %s could not be found.\nPlease specify a .db file to use for the database of the API.", dbfile)
	}
}

func checkArgonEnv() {
	time := os.Getenv("ARGON_TIME")
	if time == "" {
		log.Fatalf("Unable to get \"ARGON_TIME\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
	n, err := strconv.Atoi(time)
	if err != nil || n <= 0 {
		log.Fatalf("Variable \"ARGON_TIME\" must be a positive integer greater than 0.")
	}

	memory := os.Getenv("ARGON_MEMORY")
	if memory == "" {
		log.Fatalf("Unable to get \"ARGON_MEMORY\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
	n, err = strconv.Atoi(memory)
	if err != nil || n <= 0 {
		log.Fatalf("Variable \"ARGON_MEMORY\" must be a positive integer greater than 0.")
	}

	threads := os.Getenv("ARGON_THREADS")
	if threads == "" {
		log.Fatalf("Unable to get \"ARGON_THREADS\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
	n, err = strconv.Atoi(threads)
	if err != nil || n <= 0 {
		log.Fatalf("Variable \"ARGON_THREADS\" must be a positive integer greater than 0.")
	}

	hashLength := os.Getenv("ARGON_HASH_LENGTH")
	if hashLength == "" {
		log.Fatalf("Unable to get \"ARGON_HASH_LENGTH\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
	n, err = strconv.Atoi(hashLength)
	if err != nil || n <= 0 {
		log.Fatalf("Variable \"ARGON_HASH_LENGTH\" must be a positive integer greater than 0.")
	}
}

func checkHTTPEnv() {
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		log.Fatalf("Unable to get \"ADDRESS\" environmnet variable.\nPlease make sure to set it before starting the API.")
	}
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
