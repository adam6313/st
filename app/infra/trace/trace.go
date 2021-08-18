package trace

import (
	"storage/app/infra/config"

	"github.com/tyr-tech-team/hawk/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
)

// NewTrace -
func NewTrace(c config.Config) {
	exp, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(c.Trace.URL),
		),
	)
	if err != nil {
		panic(err)
	}

	trace.TracerProvider(trace.Config{
		Exporter:    exp,
		Service:     c.Info.Name,
		Environment: c.Trace.Environment,
	})
}
