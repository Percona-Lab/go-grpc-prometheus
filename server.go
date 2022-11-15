// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for server-side gRPC.

package grpc_prometheus

import (
	"context"
	"sync"

	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// defaultServerMetrics is the default instance of ServerMetrics. It is
	// intended to be used in conjunction the default Prometheus metrics
	// registry.
	defaultServerMetrics     *ServerMetrics
	defaultServerMetricsOnce sync.Once
)

func DefaultServerMetrics() *ServerMetrics {
	defaultServerMetricsOnce.Do(func() {
		defaultServerMetrics = NewServerMetrics()

		PrometheusMustRegister(defaultServerMetrics)
	})

	return defaultServerMetrics
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return DefaultServerMetrics().UnaryServerInterceptor(ctx, req, info, handler)
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return DefaultServerMetrics().StreamServerInterceptor(srv, ss, info, handler)
}

func PrometheusMustRegister(serverMetrics *ServerMetrics) {
	prom.MustRegister(serverMetrics.serverStartedCounter)
	prom.MustRegister(serverMetrics.serverHandledCounter)
	prom.MustRegister(serverMetrics.serverStreamMsgReceivedCounter)
	prom.MustRegister(serverMetrics.serverStreamMsgSentCounter)
}

// Register takes a gRPC server and pre-initializes all counters to 0. This
// allows for easier monitoring in Prometheus (no missing metrics), and should
// be called *after* all services have been registered with the server. This
// function acts on the DefaultServerMetrics variable.
func Register(server *grpc.Server) {
	DefaultServerMetrics().InitializeMetrics(server)
}

// DefaultEnableHandlingTimeHistogram turns on recording of handling time
// of RPCs. Histogram metrics can be very expensive for Prometheus
// to retain and query. This function acts on the DefaultServerMetrics
// variable and the default Prometheus metrics registry.
func EnableHandlingTimeHistogram(opts ...HistogramOption) {
	CustomEnableHandlingTimeHistogram(DefaultServerMetrics(), opts...)
}

func CustomEnableHandlingTimeHistogram(serverMetrics *ServerMetrics, opts ...HistogramOption) {
	serverMetrics.EnableHandlingTimeHistogram(opts...)
	prom.Register(serverMetrics.serverHandledHistogram)
}
