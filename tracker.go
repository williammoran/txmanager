package txmanager

import "potentialtech.com/erpg/pkg/report"

// tracker exists to make tracking down errors easier
type tracker struct {
	file string
	line int
}

func (t *tracker) Finalize() error {
	return nil
}

func (t *tracker) Commit() {}

// Abort reports the file/line the transaction was created
// on to make tracking down errors easier
func (t *tracker) Abort() {
	report.Msgf(
		report.Alert,
		"Aborting trasaction started %s:%d",
		t.file, t.line,
	)
}
