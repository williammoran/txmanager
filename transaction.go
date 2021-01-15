package txmanager

import (
	"log"
	"sync"
)

// Transaction is a manager to group multiple persistence
// changes together. To use it, create a Transaction
// object, then use the Add() method to add TxFinalizers
// to the Transaction object. When all data changes are
// staged, call Commit() to ensure that all changes are
// committed together or Abort() to roll back all changes.
type Transaction struct {
	mutex     sync.Mutex
	once      sync.Once
	committed bool
	aborted   bool
	handlers  map[string]TxFinalizer
}

// Add registers a transaction handler
func (tx *Transaction) Add(name string, handler TxFinalizer) {
	tx.once.Do(tx.setup)
	tx.mutex.Lock()
	defer tx.mutex.Unlock()
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

// Commit Finalize()s all backends and Commit()s them
// if all succeed. Abort()s all Finalizers if anything
// goes wrong.
// With properly implemented finalizers, this provides a
// guarantee that either all data was committed, or none.
func (tx *Transaction) Commit() error {
	if tx.aborted {
		return MakeError("Transaction already aborted")
	}
	if tx.committed {
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
		e := f.Commit()
		if e != nil {
			tx.Abort(e.Error())
			return e
		}
	}
	tx.committed = true
	return nil
}

// Abort rolls back all storage operations if an error
// is pending. If Commit() was successfully called,
// nothing is done.
func (tx *Transaction) Abort(msg string) {
	if !tx.committed && !tx.aborted {
		for _, f := range tx.handlers {
			f.Abort()
		}
		tx.aborted = true
	}
}

func (tx *Transaction) setup() {
	tx.handlers = make(map[string]TxFinalizer)
}
