package grpc_prometheus

import "context"

type Extension interface {
	ServerSentMessageCustomLabels() []string
	ServerSentMessageValues(ctx context.Context) []string
	ServerReceivedMessageCustomLabels() []string
	ServerReceivedMessageValues(ctx context.Context) []string
}

type defaultExtension struct {
}

func (d *defaultExtension) ServerSentMessageCustomLabels() []string {
	return nil
}

func (d *defaultExtension) ServerSentMessageValues(ctx context.Context) []string {
	return nil
}

func (d *defaultExtension) ServerReceivedMessageCustomLabels() []string {
	return nil
}

func (d *defaultExtension) ServerReceivedMessageValues(ctx context.Context) []string {
	return nil
}
