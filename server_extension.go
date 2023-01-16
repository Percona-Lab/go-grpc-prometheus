package grpc_prometheus

import "context"

type ServerExtension interface {
	MetricsNameAdjust(name string) string

	ServerHandledCounterCustomLabels() []string
	ServerHandledCounterValues(ctx context.Context) []string

	ServerStreamMsgReceivedCounterCustomLabels() []string
	ServerStreamMsgReceivedCounterValues(ctx context.Context) []string

	ServerStreamMsgSentCounterCustomLabels() []string
	ServerStreamMsgSentCounterValues(ctx context.Context) []string

	ServerHandledHistogramCustomLabels() []string
	ServerHandledHistogramValues(ctx context.Context) []string
}

type DefaultExtension struct {
}

var emptyExtension ServerExtension = DefaultExtension{}

func (e DefaultExtension) MetricsNameAdjust(name string) string {
	return name
}

func (DefaultExtension) ServerHandledCounterCustomLabels() []string {
	return nil
}

func (DefaultExtension) ServerHandledCounterValues(context.Context) []string {
	return nil
}

func (DefaultExtension) ServerStreamMsgReceivedCounterCustomLabels() []string {
	return nil
}

func (DefaultExtension) ServerStreamMsgReceivedCounterValues(context.Context) []string {
	return nil
}

func (DefaultExtension) ServerStreamMsgSentCounterCustomLabels() []string {
	return nil
}

func (DefaultExtension) ServerStreamMsgSentCounterValues(context.Context) []string {
	return nil
}

func (DefaultExtension) ServerHandledHistogramCustomLabels() []string {
	return nil
}

func (DefaultExtension) ServerHandledHistogramValues(context.Context) []string {
	return nil
}
