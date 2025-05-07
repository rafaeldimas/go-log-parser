package pipeline

import (
	"sync"

	"github.com/rafaeldimas/go-log-parser/pkg/pipeline/stage"
	"github.com/rafaeldimas/go-log-parser/pkg/queue"
)

type Pipeline interface {
	Done() <-chan struct{}
	AddPipe(stage stage.StagePipe)
	Execute(executor stage.StageExecutor)
}

type pipeline struct {
	generator stage.StageGenerator
	stages    queue.Queue[stage.StagePipe]
	done      chan struct{}
	wg        *sync.WaitGroup
}

func New(generator stage.StageGenerator) Pipeline {
	return &pipeline{
		generator: generator,
		stages:    queue.New[stage.StagePipe](),
		done:      make(chan struct{}),
		wg:        &sync.WaitGroup{},
	}
}

func (p *pipeline) Done() <-chan struct{} {
	return p.done
}

func (p *pipeline) AddPipe(stage stage.StagePipe) {
	p.stages.Enqueue(stage)
}

func (p *pipeline) Execute(executor stage.StageExecutor) {
	p.wg.Add(p.stages.Length() + 2)
	out := p.generator.Execute(p.wg)

	for !p.stages.IsEmpty() {
		stage := p.stages.Dequeue()
		out = stage.Execute(out, p.wg)
	}

	executor.Execute(out, p.wg)

	p.wg.Wait()

	close(p.done)
}
