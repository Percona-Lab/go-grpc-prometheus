package grpc_prometheus

import "context"

type ServerExtension interface {
	ServerStartedCounterCustomLabels() []string
	ServerStartedCounterPreRegisterValues() [][]string
	ServerStartedCounterValues(ctx context.Context) []string

	ServerHandledCounterCustomLabels() []string
	ServerHandledCounterPreRegisterValues() [][]string
	ServerHandledCounterValues(ctx context.Context) []string

	ServerStreamMsgReceivedCounterCustomLabels() []string
	ServerStreamMsgReceivedCounterPreRegisterValues() [][]string
	ServerStreamMsgReceivedCounterValues(ctx context.Context) []string

	ServerStreamMsgSentCounterCustomLabels() []string
	ServerStreamMsgSentCounterPreRegisterValues() [][]string
	ServerStreamMsgSentCounterValues(ctx context.Context) []string

	ServerHandledHistogramCustomLabels() []string
	ServerHandledHistogramPreRegisterValues() [][]string
	ServerHandledHistogramValues(ctx context.Context) []string
}

type NullExtension struct {
}

var _ ServerExtension = NullExtension{}

func (NullExtension) ServerStartedCounterCustomLabels() []string {
	return nil
}

func (NullExtension) ServerStartedCounterPreRegisterValues() [][]string {
	return nil
}

func (NullExtension) ServerStartedCounterValues(context.Context) []string {
	return nil
}

func (NullExtension) ServerHandledCounterCustomLabels() []string {
	return nil
}

func (NullExtension) ServerHandledCounterPreRegisterValues() [][]string {
	return nil
}

func (NullExtension) ServerHandledCounterValues(context.Context) []string {
	return nil
}

func (NullExtension) ServerStreamMsgReceivedCounterCustomLabels() []string {
	return nil
}

func (NullExtension) ServerStreamMsgReceivedCounterPreRegisterValues() [][]string {
	return nil
}

func (NullExtension) ServerStreamMsgReceivedCounterValues(context.Context) []string {
	return nil
}

func (NullExtension) ServerStreamMsgSentCounterCustomLabels() []string {
	return nil
}

func (NullExtension) ServerStreamMsgSentCounterPreRegisterValues() [][]string {
	return nil
}

func (NullExtension) ServerStreamMsgSentCounterValues(context.Context) []string {
	return nil
}

func (NullExtension) ServerHandledHistogramCustomLabels() []string {
	return nil
}

func (NullExtension) ServerHandledHistogramPreRegisterValues() [][]string {
	return nil
}

func (NullExtension) ServerHandledHistogramValues(context.Context) []string {
	return nil
}
