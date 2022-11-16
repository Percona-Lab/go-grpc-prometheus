// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for server-side gRPC.

package grpc_prometheus

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

// PrometheusMustRegister when many servers with different const label
func PrometheusMustRegister(serverMetrics *ServerMetrics) {
	prom.MustRegister(serverMetrics.serverStartedCounter)
	prom.MustRegister(serverMetrics.serverHandledCounter)
	prom.MustRegister(serverMetrics.serverStreamMsgReceivedCounter)
	prom.MustRegister(serverMetrics.serverStreamMsgSentCounter)
}

// EnableHandlingTimeHistogram turns on recording of handling time
// of RPCs. Histogram metrics can be very expensive for Prometheus
// to retain and query. This function acts on the DefaultServerMetrics
// variable and the default Prometheus metrics registry.
func EnableHandlingTimeHistogram(serverMetrics *ServerMetrics, opts ...HistogramOption) {
	serverMetrics.EnableHandlingTimeHistogram(opts...)
	prom.Register(serverMetrics.serverHandledHistogram)
}
