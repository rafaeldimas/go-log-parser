package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rafaeldimas/go-log-parser/internal/database"
	"github.com/rafaeldimas/go-log-parser/internal/parser"
	"github.com/rafaeldimas/go-log-parser/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "[PROCESS - main] ", log.LstdFlags)
	start := time.Now()

	pathFile := flag.String("file", "./tmp/fake_logs.txt", "Path file")
	flag.Parse()

	concurrencyProcess(*pathFile)
	// noConcurrencyProcess(*pathFile)

	logger.Println(time.Since(start))
}

func concurrencyProcess(pathFile string) {
	logger := log.New(os.Stdout, "[PROCESS - concurrency] ", log.LstdFlags)
	logger.Println("Start concurrency process")

	reader := storage.NewLocalStorage(pathFile)
	defer reader.Close()

	if reader.FileIsEmpty() {
		logger.Fatal("File is empty")
		return
	}

	var wg sync.WaitGroup
	done := make(chan bool, 10000)
	defer close(done)

	for {
		line := reader.Read()

		if line == "" {
			break
		}

		wg.Add(1)
		done <- true

		go func() {
			defer wg.Done()
			defer func() { <-done }()

			ps := parser.NewParser()
			logInfo, err := ps.Parse(line)

			if err != nil {
				logger.Printf("Error to parser line: %v", line)
				logger.Printf("Error: %e", err)
			}

			mongoClient, err := database.NewClient()

			if err != nil {
				logger.Fatal(err)
			}

			defer database.Close(mongoClient)

			db := database.NewDatabase(mongoClient)
			collection := database.NewCollection(db)

			collection.InsertOne(context.Background(), logInfo)
		}()
	}

	wg.Wait()
}

func noConcurrencyProcess(pathFile string) {
	logger := log.New(os.Stdout, "[PROCESS - concurrency] ", log.LstdFlags)
	logger.Println("Start no concurrency process")

	reader := storage.NewLocalStorage(pathFile)

	if reader.FileIsEmpty() {
		logger.Fatal("File is empty")
		return
	}

	mongoClient, err := database.NewClient()

	if err != nil {
		logger.Fatal(err)
	}

	defer database.Close(mongoClient)

	db := database.NewDatabase(mongoClient)
	collection := database.NewCollection(db)

	for {
		line := reader.Read()

		if line == "" {
			break
		}

		ps := parser.NewParser()
		logInfo, err := ps.Parse(line)

		if err != nil {
			logger.Printf("Error to parser line: %v", line)
			logger.Printf("Error: %e", err)
			continue
		}

		collection.InsertOne(context.Background(), logInfo)
	}
}
