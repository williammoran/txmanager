package txmanager

import (
	"context"
	"runtime"
)

type txManagerCtxName string

const thName txManagerCtxName = "txmanager.ContextName"

// WithTx adds a transaction manager to a context
// It also adds a Tracker that holds information about
// the file and line # where this was called from
func WithTx(
	ctx context.Context, t *Transaction,
) context.Context {
	_, f, l, _ := runtime.Caller(2)
	t.Add(
		"Context transaction tracker",
		&Tracker{File: f, Line: l},
	)
	go func() {
		<-ctx.Done()
		t.Abort("context complete abort check")
	}()
	nc := context.WithValue(ctx, thName, t)
	return nc
}

// GetTx gets a transaction manager from a context
func GetTx(ctx context.Context) *Transaction {
	rv := ctx.Value(thName)
	if rv == nil {
		panic("no transaction in context")
	}
	return rv.(*Transaction)
}
