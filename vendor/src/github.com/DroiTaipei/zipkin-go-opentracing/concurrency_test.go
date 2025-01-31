package zipkintracer

import (
	"strings"
	"sync"
	"testing"

	opentracing "github.com/DroiTaipei/opentracing-go"
)

const op = "test"

func TestDebugAssertSingleGoroutine(t *testing.T) {
	tracer, err := NewTracer(
		NewInMemoryRecorder(),
		EnableSpanPool(true),
		DebugAssertSingleGoroutine(true),
		TraceID128Bit(true),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}
	sp := tracer.StartSpan(op)
	sp.LogEvent("something on my goroutine")
	wait := make(chan struct{})
	var panicked bool
	go func() {
		defer func() {
			if r := recover(); r != nil {
				_, panicked = r.(*errAssertionFailed)
			}
			close(wait)
		}()
		sp.LogEvent("something on your goroutine")
	}()
	<-wait
	if !panicked {
		t.Fatal("expected a panic")
	}
}

func TestDebugAssertUseAfterFinish(t *testing.T) {
	tracer, err := NewTracer(
		NewInMemoryRecorder(),
		EnableSpanPool(true),
		DebugAssertUseAfterFinish(true),
		TraceID128Bit(true),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}
	const msg = "I shall be finished"
	for _, double := range []bool{false, true} {
		sp := tracer.StartSpan(op)
		sp.Log(opentracing.LogData{Event: msg})
		if double {
			sp.Finish()
		}
		var panicked bool
		func() {
			defer func() {
				r := recover()
				var assertionErr error
				assertionErr, panicked = r.(*errAssertionFailed)
				if !panicked && r != nil {
					panic(r)
				}
				if panicked && !strings.Contains(assertionErr.Error(), msg) {
					t.Fatalf("debug output did not contain log message '%s': %+v", msg, assertionErr)
				}
				spImpl := sp.(*spanImpl)
				// The panic should leave the Mutex unlocked.
				spImpl.Mutex.Lock()
				spImpl.Mutex.Unlock()
			}()
			sp.Finish()
		}()
		if panicked != double {
			t.Errorf("finished double = %t, but panicked = %t", double, panicked)
		}
	}
}

func TestConcurrentUsage(t *testing.T) {
	var cr CountingRecorder
	tracer, err := NewTracer(
		&cr,
		EnableSpanPool(true),
		DebugAssertSingleGoroutine(true),
		TraceID128Bit(true),
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}
	var wg sync.WaitGroup
	const num = 100
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < num; j++ {
				sp := tracer.StartSpan(op)
				sp.LogEvent("test event")
				sp.SetTag("foo", "bar")
				sp.SetBaggageItem("boo", "far")
				sp.SetOperationName("x")
				csp := tracer.StartSpan(
					"csp",
					opentracing.ChildOf(sp.Context()))
				csp.Finish()
				defer sp.Finish()
			}
		}()
	}
	wg.Wait()
}

func TestDisableSpanPool(t *testing.T) {
	var cr CountingRecorder
	tracer, err := NewTracer(
		&cr,
	)
	if err != nil {
		t.Fatalf("Unable to create Tracer: %+v", err)
	}

	parent := tracer.StartSpan("parent")
	parent.Finish()
	// This shouldn't panic.
	child := tracer.StartSpan(
		"child",
		opentracing.ChildOf(parent.Context()))
	child.Finish()
}
