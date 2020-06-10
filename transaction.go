package core

import (
	"log"
	"sync"
)

// Transaction is a manager to group multiple persistance
// changes together
type Transaction struct {
	sync.Mutex
	sync.Once
	comitted bool
	aborted  bool
	handlers map[string]TxFinalizer
}

// Add registers a transaction handler
func (tx *Transaction) Add(name string, handler TxFinalizer) {
	tx.Once.Do(tx.setup)
	tx.Lock()
	defer tx.Unlock()
	_, ok := tx.handlers[name]
	if ok {
		log.Panicf(
			"transaction handler %s already registered",
			name,
		)
	}
	tx.handlers[name] = handler
}

// Finalizer returns the finalizer for the given name
func (tx *Transaction) Finalizer(name string) TxFinalizer {
	return tx.handlers[name]
}

// Commit prepares commits on all backends and finalizes
// them if all succeed
func (tx *Transaction) Commit() error {
	return tx.commit()
}

func (tx *Transaction) commit() error {
	if tx.aborted {
		return MakeError("Transaction already aborted")
	}
	if tx.comitted {
		return MakeError("Transaction already comitted")
	}
	for _, f := range tx.handlers {
		e := f.Finalize()
		if e != nil {
			tx.Abort(e.Error())
			return e
		}
	}
	for _, f := range tx.handlers {
		f.Commit()
	}
	tx.comitted = true
	return nil
}

// Abort rolls back all storage operations if an error
// is pending. If Commit() was successfully called,
// nothing is done.
func (tx *Transaction) Abort(msg string) {
	if !tx.comitted && !tx.aborted {
		for _, f := range tx.handlers {
			f.Abort()
		}
		tx.aborted = true
	}
}

func (tx *Transaction) setup() {
	tx.handlers = make(map[string]TxFinalizer)
}
