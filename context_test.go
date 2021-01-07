package txmanager

import (
	"context"
	"testing"
	"time"
)

func TestAbort(t *testing.T) {
	abortCalled := false
	tm := Transaction{}
	f := MakeFinalizer("test")
	f.Register(nil, nil, func() {
		abortCalled = true
	})
	tm.Add("test", f)
	ctx, cancel := context.WithCancel(context.Background())
	WithTx(ctx, &tm)
	cancel()
	// Without the sleep, the check may happen before the
	// abort thread completes, thus creating a false failure
	time.Sleep(1 * time.Second)
	if !tm.aborted {
		t.Fatal("Transaction not aborted")
	}
	if !abortCalled {
		t.Fatal("Finalizer abort not called")
	}
}
