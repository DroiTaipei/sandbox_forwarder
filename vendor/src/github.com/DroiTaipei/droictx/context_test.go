package droictx

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/DroiTaipei/droipkg"
)

func TestConcurrentKV(t *testing.T) {
	// no panic is good
	var wg sync.WaitGroup
	var ctx Context
	ctx = &DoneContext{}
	ctx.Set("whatever", 0)
	for i := 0; i < 20; i++ {
		go func(ctx Context, i int) {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			ctx.Set("whatever", i)
			_, ok := ctx.GetInt("whatever")
			if !ok {
				t.Error("concurrent fail")
			}
			wg.Done()
		}(ctx, i)
		wg.Add(1)
	}

	for i := 0; i < 20; i++ {
		go func(ctx Context) {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			_, ok := ctx.GetInt("whatever")
			if !ok {
				t.Error("concurrent fail")
			}
			wg.Done()
		}(ctx)
		wg.Add(1)
	}
	wg.Wait()
}

func TestDeadline(t *testing.T) {
	var ctx Context
	ctx = &DoneContext{}

	d, ok := ctx.Deadline()
	if ok {
		t.Error("deadline not init, ok should be false")
	}

	if ctx.IsTimeout() {
		t.Error("timeout not init, so never timeout")
	}

	ctx.SetTimeout(500*time.Millisecond, nil)
	d, ok = ctx.Deadline()
	if !ok {
		t.Error("deadline has init, ok should be true")
	}
	time.Sleep(500 * time.Millisecond)
	if !d.Before(time.Now()) {
		t.Error("deadline wrong")
	}

	if !ctx.IsTimeout() {
		t.Error("timeout error")
	}
}

func TestDone(t *testing.T) {
	// test Context type not *DoneContext
	var ctx Context
	ctx = &DoneContext{}
	ctx.SetTimeout(5000*time.Millisecond, nil)
	go func(ctx Context) {
		time.Sleep(500 * time.Millisecond)
		ctx.Finish()
	}(ctx)

	select {
	case <-ctx.Done():
		// check finish() have stopped timer
		time.Sleep(100 * time.Millisecond)
		if ctx.StopTimer() {
			t.Error("Finish() haven't stop timer")
		}
	case <-ctx.Timeout():
		t.Error("done doesn't work")
	}
}

func TestTimeout(t *testing.T) {
	err := droipkg.ConstDroiError("1000000 timeout")
	// test Context type not *DoneContext
	var ctx Context
	ctx = &DoneContext{}
	ctx.SetTimeout(50*time.Millisecond, err)

	select {
	case <-time.After(1 * time.Second):
		ctx.StopTimer()
		t.Error("context overslept")
	case <-ctx.Timeout():
		errTimeout := ctx.TimeoutErr()
		if errTimeout.Error() != err.Error() {
			t.Error("timeout error fail")
		}
	}

	if !ctx.IsTimeout() {
		t.Error("isTimeout doesn't work")
	}
}
