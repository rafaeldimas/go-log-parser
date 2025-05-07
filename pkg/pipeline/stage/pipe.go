package stage

import "sync"

type StagePipe interface {
	Send(data any)
	Execute(in *chan any, wg *sync.WaitGroup) *chan any
	Done()
}

type Pipe func(stage StagePipe, in *chan any)

type stagePipe struct {
	out  *chan any
	pipe Pipe
	wg   *sync.WaitGroup
}

func (s *stagePipe) Send(data any) {
	*s.out <- data
}

func (s *stagePipe) Execute(in *chan any, wg *sync.WaitGroup) *chan any {
	s.wg = wg
	go s.pipe(s, in)
	return s.out
}

func (s *stagePipe) Done() {
	close(*s.out)
	s.wg.Done()
}

func NewPipe(pipe Pipe, bufferSize int) StagePipe {
	out := make(chan any, bufferSize)
	return &stagePipe{
		pipe: pipe,
		out:  &out,
	}
}
