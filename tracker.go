package txmanager

import "log"

// Tracker is a TxFinalizer that makes tracking down
// errors easier
type Tracker struct {
	File string
	Line int
}

// Finalize is a NOOP
func (t *Tracker) Finalize() error {
	return nil
}

// Commit is a NOOP
func (t *Tracker) Commit() error { return nil }

// Abort reports the file/line the Tracker was created
// with to make tracking down errors easier
func (t *Tracker) Abort() {
	log.Printf(
		"Aborting trasaction started %s:%d",
		t.File, t.Line,
	)
}
