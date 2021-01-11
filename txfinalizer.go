package txmanager

// TxFinalizer does the job of committing
type TxFinalizer interface {
	// Finalize checks to ensure that the transaction on
	// this finalizer is guaranteed to succeed on Commit
	Finalize() error
	// Commit actually finishes this finalizer with changes
	// saved
	Commit() error
	// Abort does anything necessary to abort this finalizer
	Abort()
}
