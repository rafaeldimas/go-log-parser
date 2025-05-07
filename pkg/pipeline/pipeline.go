package pipeline

import (
	"context"
	"sync"
)

type WgContextKey struct{}

type Pipeline[ChannelType any] interface {
	Done() <-chan struct{}
	AddStage(stage Stage[ChannelType])
	Execute()
}

type pipeline[ChannelType any] struct {
	ctx    context.Context
	stages []Stage[ChannelType]
	done   chan struct{}
}

func NewPipeline[ChannelType any](ctx context.Context) Pipeline[ChannelType] {
	ctx = context.WithValue(ctx, WgContextKey{}, &sync.WaitGroup{})
	return &pipeline[ChannelType]{
		ctx:  ctx,
		done: make(chan struct{}),
	}
}

func (p *pipeline[ChannelType]) Done() <-chan struct{} {
	return p.done
}

func (p *pipeline[ChannelType]) AddStage(stage Stage[ChannelType]) {
	p.stages = append(p.stages, stage)
}

func (p *pipeline[ChannelType]) Execute() {
	wg := p.ctx.Value(WgContextKey{}).(*sync.WaitGroup)
	wg.Add(len(p.stages))

	var prevOutput *chan ChannelType
	for _, stage := range p.stages {
		prevOutput = stage.Execute(p.ctx, prevOutput)
	}
	wg.Wait()
	close(p.done)
}
