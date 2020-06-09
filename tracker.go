package txmanager

import "log"

// Tracker exists to make tracking down errors easier
type Tracker struct {
	File   string
	Line   int
	Logger log.Logger
}

// Finalize is a NOOP
func (t *Tracker) Finalize() error {
	return nil
}

// Commit is a NOOP
func (t *Tracker) Commit() {}

// Abort reports the file/line the transaction was created
// on to make tracking down errors easier
func (t *Tracker) Abort() {
	t.Logger.Printf(
		"Aborting trasaction started %s:%d",
		t.File, t.Line,
	)
}
