// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prometheus

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
)

type serverReporter struct {
	metrics     *ServerMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
}

func newServerReporter(ctx context.Context, m *ServerMetrics, rpcType grpcType, fullMethod string) *serverReporter {
	r := &serverReporter{
		metrics: m,
		rpcType: rpcType,
	}
	if r.metrics.serverHandledHistogramEnabled {
		r.startTime = time.Now()
	}
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	r.metrics.serverStartedCounter.WithLabelValues(append(
		r.metrics.extension.ServerStartedCounterValues(ctx),
		string(r.rpcType), r.serviceName, r.methodName)...,
	).Inc()
	return r
}

func (r *serverReporter) ReceivedMessage(ctx context.Context) {
	r.metrics.serverStreamMsgReceivedCounter.WithLabelValues(append(
		r.metrics.extension.ServerStreamMsgReceivedCounterValues(ctx),
		string(r.rpcType), r.serviceName, r.methodName)...,
	).Inc()
}

func (r *serverReporter) SentMessage(ctx context.Context) {
	r.metrics.serverStreamMsgSentCounter.WithLabelValues(append(
		r.metrics.extension.ServerStreamMsgSentCounterValues(ctx),
		string(r.rpcType), r.serviceName, r.methodName)...,
	).Inc()
}

func (r *serverReporter) Handled(ctx context.Context, code codes.Code) {
	r.metrics.serverHandledCounter.WithLabelValues(append(
		r.metrics.extension.ServerHandledCounterValues(ctx),
		string(r.rpcType), r.serviceName, r.methodName, code.String())...,
	).Inc()

	if r.metrics.serverHandledHistogramEnabled {
		r.metrics.serverHandledHistogram.WithLabelValues(append(
			r.metrics.extension.ServerHandledHistogramValues(ctx),
			string(r.rpcType), r.serviceName, r.methodName)...,
		).Observe(time.Since(r.startTime).Seconds())
	}
}
