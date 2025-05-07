package stage

import "sync"

type StageGenerator interface {
	Send(data any)
	Execute(wg *sync.WaitGroup) *chan any
	Done()
}

type Generator func(stage StageGenerator)

type stageGenerator struct {
	out       *chan any
	generator Generator
	wg        *sync.WaitGroup
}

func (s *stageGenerator) Send(data any) {
	*s.out <- data
}

func (s *stageGenerator) Execute(wg *sync.WaitGroup) *chan any {
	s.wg = wg
	go s.generator(s)
	return s.out
}

func (s *stageGenerator) Done() {
	close(*s.out)
	s.wg.Done()
}

func NewGenerator(generator Generator, bufferSize int) StageGenerator {
	out := make(chan any, bufferSize)
	return &stageGenerator{
		generator: generator,
		out:       &out,
	}
}
