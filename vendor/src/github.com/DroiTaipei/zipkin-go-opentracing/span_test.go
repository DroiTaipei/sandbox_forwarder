package zipkintracer

import (
	"reflect"
	"strconv"
	"testing"

	opentracing "github.com/DroiTaipei/opentracing-go"
	"github.com/DroiTaipei/opentracing-go/ext"
	"github.com/DroiTaipei/opentracing-go/log"
	"github.com/stretchr/testify/assert"
)

func TestSpan_Baggage(t *testing.T) {
	recorder := NewInMemoryRecorder()
	tracer, err := NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return true }),
		WithLogger(&nopLogger{}),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	span := tracer.StartSpan("x")
	span.SetBaggageItem("x", "y")
	assert.Equal(t, "y", span.BaggageItem("x"))
	span.Finish()
	spans := recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, map[string]string{"x": "y"}, spans[0].Context.Baggage)

	recorder.Reset()
	span = tracer.StartSpan("x")
	span.SetBaggageItem("x", "y")
	baggage := make(map[string]string)
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return true
	})
	assert.Equal(t, map[string]string{"x": "y"}, baggage)

	span.SetBaggageItem("a", "b")
	baggage = make(map[string]string)
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return false // exit early
	})
	assert.Equal(t, 1, len(baggage))
	span.Finish()
	spans = recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, 2, len(spans[0].Context.Baggage))
}

func TestSpan_Sampling(t *testing.T) {
	recorder := NewInMemoryRecorder()
	tracer, err := NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return true }),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	span := tracer.StartSpan("x")
	span.Finish()
	assert.Equal(t, 1, len(recorder.GetSampledSpans()), "by default span should be sampled")

	recorder.Reset()
	span = tracer.StartSpan("x")
	ext.SamplingPriority.Set(span, 0)
	span.Finish()
	assert.Equal(t, 0, len(recorder.GetSampledSpans()), "SamplingPriority=0 should turn off sampling")

	tracer, err = NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return false }),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	recorder.Reset()
	span = tracer.StartSpan("x")
	span.Finish()
	assert.Equal(t, 0, len(recorder.GetSampledSpans()), "by default span should not be sampled")

	recorder.Reset()
	span = tracer.StartSpan("x")
	ext.SamplingPriority.Set(span, 1)
	span.Finish()
	assert.Equal(t, 1, len(recorder.GetSampledSpans()), "SamplingPriority=1 should turn on sampling")
}

func TestSpan_SingleLoggedTaggedSpan(t *testing.T) {
	recorder := NewInMemoryRecorder()
	tracer, err := NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return true }),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	span := tracer.StartSpan("x")
	span.LogEventWithPayload("event", "payload")
	span.LogFields(log.String("key_str", "value"), log.Uint32("32bit", 4294967295))
	span.SetTag("tag", "value")
	span.Finish()
	spans := recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "x", spans[0].Operation)
	assert.Equal(t, 2, len(spans[0].Logs))
	assert.Equal(t, opentracing.Tags{"tag": "value"}, spans[0].Tags)
	fv := NewLogFieldValidator(t, spans[0].Logs[0].Fields)
	fv.
		ExpectNextFieldEquals("event", reflect.String, "event").
		ExpectNextFieldEquals("payload", reflect.Interface, "payload")
	fv = NewLogFieldValidator(t, spans[0].Logs[1].Fields)
	fv.
		ExpectNextFieldEquals("key_str", reflect.String, "value").
		ExpectNextFieldEquals("32bit", reflect.Uint32, "4294967295")
}

func TestSpan_TrimUnsampledSpans(t *testing.T) {
	recorder := NewInMemoryRecorder()
	// Tracer that trims only unsampled but always samples
	tracer, err := NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return true }),
		TrimUnsampledSpans(true),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	span := tracer.StartSpan("x")
	span.LogFields(log.String("key_str", "value"), log.Uint32("32bit", 4294967295))
	span.SetTag("tag", "value")
	span.Finish()
	spans := recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, 1, len(spans[0].Logs))
	assert.Equal(t, opentracing.Tags{"tag": "value"}, spans[0].Tags)
	fv := NewLogFieldValidator(t, spans[0].Logs[0].Fields)
	fv.
		ExpectNextFieldEquals("key_str", reflect.String, "value").
		ExpectNextFieldEquals("32bit", reflect.Uint32, "4294967295")

	recorder.Reset()
	// Tracer that trims only unsampled and never samples
	tracer, err = NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return false }),
		TrimUnsampledSpans(true),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	span = tracer.StartSpan("x")
	span.LogFields(log.String("key_str", "value"), log.Uint32("32bit", 4294967295))
	span.SetTag("tag", "value")
	span.Finish()
	spans = recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, 0, len(spans[0].Logs))
	assert.Equal(t, 0, len(spans[0].Tags))
}

func TestSpan_DropAllLogs(t *testing.T) {
	recorder := NewInMemoryRecorder()
	// Tracer that drops logs
	tracer, err := NewTracer(
		recorder,
		WithSampler(func(_ uint64) bool { return true }),
		DropAllLogs(true),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	span := tracer.StartSpan("x")
	span.LogFields(log.String("key_str", "value"), log.Uint32("32bit", 4294967295))
	span.SetTag("tag", "value")
	span.Finish()
	spans := recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "x", spans[0].Operation)
	assert.Equal(t, opentracing.Tags{"tag": "value"}, spans[0].Tags)
	// Only logs are dropped
	assert.Equal(t, 0, len(spans[0].Logs))
}

func TestSpan_MaxLogSperSpan(t *testing.T) {
	for _, limit := range []int{5, 10, 15, 20, 30, 40, 50} {
		for _, numLogs := range []int{5, 10, 15, 20, 30, 40, 50, 60, 70, 80} {
			recorder := NewInMemoryRecorder()
			// Tracer that only retains the last <limit> logs.
			tracer, err := NewTracer(
				recorder,
				WithSampler(func(_ uint64) bool { return true }),
				WithMaxLogsPerSpan(limit),
			)
			if err != nil {
				t.Fatalf("Unable to create Tracer: %+v", err)
			}

			span := tracer.StartSpan("x")
			for i := 0; i < numLogs; i++ {
				span.LogKV("eventIdx", i)
			}
			span.Finish()

			spans := recorder.GetSpans()
			assert.Equal(t, 1, len(spans))
			assert.Equal(t, "x", spans[0].Operation)

			logs := spans[0].Logs
			var firstLogs, lastLogs []opentracing.LogRecord
			if numLogs <= limit {
				assert.Equal(t, numLogs, len(logs))
				firstLogs = logs
			} else {
				assert.Equal(t, limit, len(logs))
				if len(logs) > 0 {
					numOld := (len(logs) - 1) / 2
					firstLogs = logs[:numOld]
					lastLogs = logs[numOld+1:]

					fv := NewLogFieldValidator(t, logs[numOld].Fields)
					fv = fv.ExpectNextFieldEquals("event", reflect.String, "dropped Span logs")
					fv = fv.ExpectNextFieldEquals(
						"dropped_log_count", reflect.Int, strconv.Itoa(numLogs-limit+1),
					)
					fv.ExpectNextFieldEquals("component", reflect.String, "zipkintracer")
				}
			}

			for i, lr := range firstLogs {
				fv := NewLogFieldValidator(t, lr.Fields)
				fv.ExpectNextFieldEquals("eventIdx", reflect.Int, strconv.Itoa(i))
			}

			for i, lr := range lastLogs {
				fv := NewLogFieldValidator(t, lr.Fields)
				fv.ExpectNextFieldEquals("eventIdx", reflect.Int, strconv.Itoa(numLogs-len(lastLogs)+i))
			}
		}
	}
}
