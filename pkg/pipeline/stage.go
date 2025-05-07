package pipeline

import (
	"context"
)

type funcRoutine[ChannelType any] func(ctx context.Context, input *chan ChannelType, output *chan ChannelType)

type Stage[ChannelType any] interface {
	DefineConcorrency(concurrency int) Stage[ChannelType]
	Execute(ctx context.Context, input *chan ChannelType) *chan ChannelType
}

type stage[ChannelType any] struct {
	concorrency int
	routine     funcRoutine[ChannelType]
}

func NewPipelineStage[ChannelType any](routine funcRoutine[ChannelType]) Stage[ChannelType] {
	return &stage[ChannelType]{
		routine: routine,
	}
}

func (s *stage[ChannelType]) DefineConcorrency(concurrency int) Stage[ChannelType] {
	s.concorrency = concurrency
	return s
}

func (s *stage[ChannelType]) Execute(ctx context.Context, prevOutput *chan ChannelType) *chan ChannelType {
	output := make(chan ChannelType, s.concorrency)

	go s.routine(ctx, prevOutput, &output)

	return &output
}
