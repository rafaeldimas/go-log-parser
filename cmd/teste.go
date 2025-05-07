package main

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/rafaeldimas/go-log-parser/pkg/pipeline"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p := pipeline.NewPipeline[string](ctx)
	p.AddStage(addDataStage())
	p.AddStage(toUpperCaseStage())
	p.AddStage(logStage())

	start := time.Now()
	p.Execute()

	select {
	case <-p.Done():
		log.Printf("pipeline done in %v", time.Since(start))
	case <-ctx.Done():
		log.Printf("context timeout error: %v, in %v", ctx.Err(), time.Since(start))
	default:
		log.Printf("context default")
	}
}

func addDataStage() pipeline.Stage[string] {
	return pipeline.NewPipelineStage(func(ctx context.Context, _ *chan string, output *chan string) {
		wg := ctx.Value(pipeline.WgContextKey{}).(*sync.WaitGroup)

		tests := generateStrings(10_000_000)
		for _, line := range tests {
			*output <- line
		}

		wg.Done()
		close(*output)
	}).DefineConcorrency(10000)
}

func toUpperCaseStage() pipeline.Stage[string] {
	return pipeline.NewPipelineStage(func(ctx context.Context, input *chan string, output *chan string) {
		wg := ctx.Value(pipeline.WgContextKey{}).(*sync.WaitGroup)

		for line := range *input {
			lineUpper := strings.ToUpper(line)
			*output <- lineUpper
		}

		wg.Done()
		close(*output)
	}).DefineConcorrency(10000)
}

func logStage() pipeline.Stage[string] {
	return pipeline.NewPipelineStage(func(ctx context.Context, input *chan string, output *chan string) {
		wg := ctx.Value(pipeline.WgContextKey{}).(*sync.WaitGroup)

		for line := range *input {
			log.Printf("log stage: %s", line)
		}

		wg.Done()
		close(*output)
	})
}

func generateStrings(tamanho int) []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	stringsAleatorias := make([]string, tamanho)
	caracteres := "abcdefghijklmnopqrstuvwxyz"

	for i := range stringsAleatorias {
		tamanhoString := rng.Intn(10) + 1
		stringAleatoria := make([]byte, tamanhoString)
		for j := range stringAleatoria {
			stringAleatoria[j] = caracteres[rng.Intn(len(caracteres))]
		}
		stringsAleatorias[i] = string(stringAleatoria)
	}
	return stringsAleatorias
}
