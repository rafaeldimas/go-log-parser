package stage

import "sync"

type StageExecutor interface {
	Execute(in *chan any, wg *sync.WaitGroup)
	Done()
}

type Executor func(stage StageExecutor, in *chan any)

type stageExecutor struct {
	executor Executor
	wg       *sync.WaitGroup
}

func (s *stageExecutor) Execute(in *chan any, wg *sync.WaitGroup) {
	s.wg = wg
	go s.executor(s, in)
}

func (s *stageExecutor) Done() {
	s.wg.Done()
}

func NewExecutor(executor Executor) StageExecutor {
	return &stageExecutor{
		executor: executor,
	}
}
