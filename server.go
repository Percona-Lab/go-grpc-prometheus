// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for server-side gRPC.

package grpc_prometheus

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

// MustRegister when many servers with different const label
func (m *ServerMetrics) MustRegister() {
	prom.MustRegister(m.serverStartedCounter)
	prom.MustRegister(m.serverHandledCounter)
	prom.MustRegister(m.serverStreamMsgReceivedCounter)
	prom.MustRegister(m.serverStreamMsgSentCounter)
}
