package grpc_prometheus

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-prometheus/packages/grpcstatus"
	prom "github.com/prometheus/client_golang/prometheus"

	"google.golang.org/grpc"
)

// ServerMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for a gRPC server.
type ServerMetrics struct {
	extension                      ServerExtension
	serverStartedCounter           *prom.CounterVec
	serverHandledCounter           *prom.CounterVec
	serverStreamMsgReceivedCounter *prom.CounterVec
	serverStreamMsgSentCounter     *prom.CounterVec
	serverHandledHistogramEnabled  bool
	serverHandledHistogramOpts     prom.HistogramOpts
	serverHandledHistogram         *prom.HistogramVec
}

// NewServerMetrics returns a ServerMetrics object. Use a new instance of
// ServerMetrics when not using the default Prometheus metrics registry, for
// example when wanting to control which metrics are added to a registry as
// opposed to automatically adding metrics via init functions.
func NewServerMetrics(counterOpts ...CounterOption) *ServerMetrics {
	return NewServerMetricsWithExtension(&DefaultExtension{}, counterOpts...)
}

func NewServerMetricsWithExtension(extension ServerExtension, counterOpts ...CounterOption) *ServerMetrics {
	opts := counterOptions(counterOpts)
	return &ServerMetrics{
		extension: extension,
		serverStartedCounter: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_started_total",
				Help: "Total number of RPCs started on the server.",
			}), []string{"grpc_type", "grpc_service", "grpc_method"}),
		serverHandledCounter: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_handled_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}), append(extension.ServerHandledCounterCustomLabels(), "grpc_type", "grpc_service", "grpc_method", "grpc_code")),
		serverStreamMsgReceivedCounter: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_msg_received_total",
				Help: "Total number of RPC stream messages received on the server.",
			}), append(extension.ServerStreamMsgReceivedCounterCustomLabels(), "grpc_type", "grpc_service", "grpc_method")),
		serverStreamMsgSentCounter: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_msg_sent_total",
				Help: "Total number of gRPC stream messages sent by the server.",
			}), append(extension.ServerStreamMsgSentCounterCustomLabels(), "grpc_type", "grpc_service", "grpc_method")),
		serverHandledHistogramEnabled: false,
		serverHandledHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_handling_seconds",
			Help:    "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
			Buckets: prom.DefBuckets,
		},
		serverHandledHistogram: nil,
	}
}

// EnableHandlingTimeHistogram enables histograms being registered when
// registering the ServerMetrics on a Prometheus registry. Histograms can be
// expensive on Prometheus servers. It takes options to configure histogram
// options such as the defined buckets.
func (m *ServerMetrics) EnableHandlingTimeHistogram(opts ...HistogramOption) {
	for _, o := range opts {
		o(&m.serverHandledHistogramOpts)
	}
	if !m.serverHandledHistogramEnabled {
		m.serverHandledHistogram = prom.NewHistogramVec(
			m.serverHandledHistogramOpts,
			[]string{"grpc_type", "grpc_service", "grpc_method"},
		)
	}
	m.serverHandledHistogramEnabled = true
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ServerMetrics) Describe(ch chan<- *prom.Desc) {
	m.serverStartedCounter.Describe(ch)
	m.serverHandledCounter.Describe(ch)
	m.serverStreamMsgReceivedCounter.Describe(ch)
	m.serverStreamMsgSentCounter.Describe(ch)
	if m.serverHandledHistogramEnabled {
		m.serverHandledHistogram.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerMetrics) Collect(ch chan<- prom.Metric) {
	m.serverStartedCounter.Collect(ch)
	m.serverHandledCounter.Collect(ch)
	m.serverStreamMsgReceivedCounter.Collect(ch)
	m.serverStreamMsgSentCounter.Collect(ch)
	if m.serverHandledHistogramEnabled {
		m.serverHandledHistogram.Collect(ch)
	}
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ServerMetrics) UnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		monitor := newServerReporter(m, Unary, info.FullMethod)
		monitor.ReceivedMessage(ctx)
		resp, err := handler(ctx, req)
		st, _ := grpcstatus.FromError(err)
		monitor.Handled(ctx, st.Code())
		if err == nil {
			monitor.SentMessage(ctx)
		}
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func (m *ServerMetrics) StreamServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		monitor := newServerReporter(m, streamRPCType(info), info.FullMethod)
		err := handler(srv, &monitoredServerStream{ss, monitor})
		st, _ := grpcstatus.FromError(err)
		monitor.Handled(ss.Context(), st.Code())
		return err
	}
}

// InitializeMetrics initializes all metrics, with their appropriate null
// value, for all gRPC methods registered on a gRPC server. This is useful, to
// ensure that all metrics exist when collecting and querying.
func (m *ServerMetrics) InitializeMetrics(server *grpc.Server) {
	serviceInfo := server.GetServiceInfo()
	for serviceName, info := range serviceInfo {
		for _, mInfo := range info.Methods {
			preRegisterMethod(m, serviceName, &mInfo)
		}
	}
}

func streamRPCType(info *grpc.StreamServerInfo) grpcType {
	if info.IsClientStream && !info.IsServerStream {
		return ClientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return ServerStream
	}
	return BidiStream
}

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct {
	grpc.ServerStream
	monitor *serverReporter
}

func (s *monitoredServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.monitor.SentMessage(s.ServerStream.Context())
	}
	return err
}

func (s *monitoredServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.monitor.ReceivedMessage(s.ServerStream.Context())
	}
	return err
}

// preRegisterMethod is invoked on Register of a Server, allowing all gRPC services labels to be pre-populated.
func preRegisterMethod(metrics *ServerMetrics, serviceName string, mInfo *grpc.MethodInfo) {
	methodName := mInfo.Name
	methodType := string(typeFromMethodInfo(mInfo))
	// These are just references (no increments), as just referencing will create the labels but not set values.
	metrics.serverStartedCounter.GetMetricWithLabelValues(methodType, serviceName, methodName)
	metrics.serverStreamMsgReceivedCounter.GetMetricWithLabelValues(methodType, serviceName, methodName)
	metrics.serverStreamMsgSentCounter.GetMetricWithLabelValues(methodType, serviceName, methodName)
	if metrics.serverHandledHistogramEnabled {
		metrics.serverHandledHistogram.GetMetricWithLabelValues(methodType, serviceName, methodName)
	}
	for _, code := range allCodes {
		metrics.serverHandledCounter.GetMetricWithLabelValues(methodType, serviceName, methodName, code.String())
	}
}
