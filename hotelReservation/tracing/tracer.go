package tracing

import (
	"os"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go/config"
)

var (
	defaultSampleRatio float64 = 0.01
)

// Init returns a newly configured tracer
func Init(serviceName, host string) (opentracing.Tracer, error) {
	ratio := defaultSampleRatio
	if val, ok := os.LookupEnv("JAEGER_SAMPLE_RATIO"); ok {
		ratio, _ = strconv.ParseFloat(val, 64)
		if ratio > 1 {
			ratio = 1.0
		}
	}
	log.Info().Msgf("Jaeger client: adjusted sample ratio %f", ratio)
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "probabilistic",
			Param: ratio,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  host,
		},
		RPCMetrics: true,
	}

	tracer, _, err := cfg.New(serviceName)
	if err != nil {
		return nil, err
	}

	// otelTracer := otel.GetTracerProvider().Tracer("tracer_name")

	// bridgeTracer, wrapperTracerProvider := otelBridge.NewTracerPair(otelTracer)

	// otel.SetTracerProvider(wrapperTracerProvider)
	return tracer, nil
}
