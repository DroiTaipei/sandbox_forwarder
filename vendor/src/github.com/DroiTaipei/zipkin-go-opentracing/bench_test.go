package zipkintracer

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	opentracing "github.com/DroiTaipei/opentracing-go"
)

var tags []string

func init() {
	tags = make([]string, 1000)
	for j := 0; j < len(tags); j++ {
		tags[j] = fmt.Sprintf("%d", randomID())
	}
}

func executeOps(sp opentracing.Span, numEvent, numTag, numItems int) {
	for j := 0; j < numEvent; j++ {
		sp.LogEvent("event")
	}
	for j := 0; j < numTag; j++ {
		sp.SetTag(tags[j], nil)
	}
	for j := 0; j < numItems; j++ {
		sp.SetBaggageItem(tags[j], tags[j])
	}
}

func benchmarkWithOps(b *testing.B, numEvent, numTag, numItems int) {
	var r CountingRecorder
	t, err := NewTracer(&r)
	if err != nil {
		b.Fatalf("Unable to create Tracer: %+v", err)
	}
	benchmarkWithOpsAndCB(b, func() opentracing.Span {
		return t.StartSpan("test")
	}, numEvent, numTag, numItems)
	if int(r) != b.N {
		b.Fatalf("missing traces: expected %d, got %d", b.N, r)
	}
}

func benchmarkWithOpsAndCB(b *testing.B, create func() opentracing.Span,
	numEvent, numTag, numItems int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sp := create()
		executeOps(sp, numEvent, numTag, numItems)
		sp.Finish()
	}
	b.StopTimer()
}

func BenchmarkSpan_Empty(b *testing.B) {
	benchmarkWithOps(b, 0, 0, 0)
}

func BenchmarkSpan_100Events(b *testing.B) {
	benchmarkWithOps(b, 100, 0, 0)
}

func BenchmarkSpan_1000Events(b *testing.B) {
	benchmarkWithOps(b, 100, 0, 0)
}

func BenchmarkSpan_100Tags(b *testing.B) {
	benchmarkWithOps(b, 0, 100, 0)
}

func BenchmarkSpan_1000Tags(b *testing.B) {
	benchmarkWithOps(b, 0, 100, 0)
}

func BenchmarkSpan_100BaggageItems(b *testing.B) {
	benchmarkWithOps(b, 0, 0, 100)
}

func BenchmarkTrimmedSpan_100Events_100Tags_100BaggageItems(b *testing.B) {
	var r CountingRecorder
	t, err := NewTracer(
		&r,
		TrimUnsampledSpans(true),
		WithSampler(neverSample),
		TraceID128Bit(true),
	)
	if err != nil {
		b.Fatalf("Unable to create Tracer: %+v", err)
	}
	benchmarkWithOpsAndCB(b, func() opentracing.Span {
		sp := t.StartSpan("test")
		return sp
	}, 100, 100, 100)
	if int(r) != b.N {
		b.Fatalf("missing traces: expected %d, got %d", b.N, r)
	}
}

func benchmarkInject(b *testing.B, format opentracing.BuiltinFormat, numItems int) {
	var r CountingRecorder
	tracer, err := NewTracer(&r)
	if err != nil {
		b.Fatalf("Unable to create Tracer: %+v", err)
	}
	sp := tracer.StartSpan("testing")
	executeOps(sp, 0, 0, numItems)
	var carrier interface{}
	switch format {
	case opentracing.TextMap, opentracing.HTTPHeaders:
		carrier = opentracing.HTTPHeadersCarrier(http.Header{})
	case opentracing.Binary:
		carrier = &bytes.Buffer{}
	default:
		b.Fatalf("unhandled format %d", format)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := tracer.Inject(sp.Context(), format, carrier)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkExtract(b *testing.B, format opentracing.BuiltinFormat, numItems int) {
	var r CountingRecorder
	tracer, err := NewTracer(&r)
	if err != nil {
		b.Fatalf("Unable to create Tracer: %+v", err)
	}
	sp := tracer.StartSpan("testing")
	executeOps(sp, 0, 0, numItems)
	var carrier interface{}
	switch format {
	case opentracing.TextMap, opentracing.HTTPHeaders:
		carrier = opentracing.HTTPHeadersCarrier(http.Header{})
	case opentracing.Binary:
		carrier = &bytes.Buffer{}
	default:
		b.Fatalf("unhandled format %d", format)
	}
	if err := tracer.Inject(sp.Context(), format, carrier); err != nil {
		b.Fatal(err)
	}

	// We create a new bytes.Buffer every time for tracer.Extract() to keep
	// this benchmark realistic.
	var rawBinaryBytes []byte
	if format == opentracing.Binary {
		rawBinaryBytes = carrier.(*bytes.Buffer).Bytes()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if format == opentracing.Binary {
			carrier = bytes.NewBuffer(rawBinaryBytes)
		}
		_, err := tracer.Extract(format, carrier)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInject_TextMap_Empty(b *testing.B) {
	benchmarkInject(b, opentracing.TextMap, 0)
}

func BenchmarkInject_TextMap_100BaggageItems(b *testing.B) {
	benchmarkInject(b, opentracing.TextMap, 100)
}

func BenchmarkInject_Binary_Empty(b *testing.B) {
	benchmarkInject(b, opentracing.Binary, 0)
}

func BenchmarkInject_Binary_100BaggageItems(b *testing.B) {
	benchmarkInject(b, opentracing.Binary, 100)
}

func BenchmarkExtract_TextMap_Empty(b *testing.B) {
	benchmarkExtract(b, opentracing.TextMap, 0)
}

func BenchmarkExtract_TextMap_100BaggageItems(b *testing.B) {
	benchmarkExtract(b, opentracing.TextMap, 100)
}

func BenchmarkExtract_Binary_Empty(b *testing.B) {
	benchmarkExtract(b, opentracing.Binary, 0)
}

func BenchmarkExtract_Binary_100BaggageItems(b *testing.B) {
	benchmarkExtract(b, opentracing.Binary, 100)
}
