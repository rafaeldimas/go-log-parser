package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rafaeldimas/go-log-parser/internal/database"
	"github.com/rafaeldimas/go-log-parser/internal/parser"
	"github.com/rafaeldimas/go-log-parser/internal/storage"
	"github.com/rafaeldimas/go-log-parser/pkg/pipeline"
	"github.com/rafaeldimas/go-log-parser/pkg/pipeline/stage"
)

func main() {
	start := time.Now()
	logger := log.New(os.Stdout, "[PROCESS - MAIN] ", log.LstdFlags)

	pathFile := flag.String("file", "./tmp/fake_logs.txt", "Path file")
	flag.Parse()

	p := pipeline.New(stage.NewGenerator(generateLines(*pathFile), 10))

	p.AddPipe(stage.NewPipe(parserLine(), 10))

	p.Execute(stage.NewExecutor(storageLogInfo()))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-p.Done():
		logger.Printf("pipeline done in %v", time.Since(start))
	case <-ctx.Done():
		logger.Printf("os interrupt in %v", time.Since(start))
	}
}

func generateLines(pathFile string) stage.Generator {
	return func(stage stage.StageGenerator) {
		defer stage.Done()

		logger := log.New(os.Stdout, "[PROCESS - GENERATOR] ", log.LstdFlags)

		reader := storage.NewLocalStorage(pathFile)
		defer reader.Close()

		if reader.FileIsEmpty() {
			logger.Fatal("File is empty")
			return
		}

		for {
			line := reader.Read()

			if line == "" {
				break
			}

			stage.Send(line)
		}
	}
}

func parserLine() stage.Pipe {
	return func(stage stage.StagePipe, in *chan any) {
		defer stage.Done()

		logger := log.New(os.Stdout, "[PROCESS - PIPE - PARSER] ", log.LstdFlags)
		ps := parser.NewParser()

		for line := range *in {
			logInfo, err := ps.Parse(line.(string))

			if err != nil {
				logger.Printf("Error to parser line: %v", line)
				logger.Printf("Error: %e", err)
				continue
			}

			stage.Send(logInfo)
		}
	}
}

func storageLogInfo() stage.Executor {
	return func(stage stage.StageExecutor, in *chan any) {
		defer stage.Done()

		logger := log.New(os.Stdout, "[PROCESS - PIPE - STORAGE] ", log.LstdFlags)
		mongoClient, err := database.NewClient()

		if err != nil {
			logger.Fatal(err)
		}

		defer database.Close(mongoClient)

		db := database.NewDatabase(mongoClient)
		collection := database.NewCollection(db)

		for logInfo := range *in {
			_, err := collection.InsertOne(context.Background(), logInfo)

			if err != nil {
				logger.Printf("Error to insert logInfo: %v", logInfo)
				logger.Printf("Error: %e", err)
				continue
			}
		}
	}
}
