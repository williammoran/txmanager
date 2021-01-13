package txmanager

// Finalizer is a code-only implementation of a TxFinalizer
// It does not require any external dependencies and can
// be used to quickly implement support for finalizers for
// most types of storage
type Finalizer struct {
	name       string
	finalizers []func() error
	commits    []func() error
	aborts     []func()
}

// NewFinalizer simple constructor to make a TxFinalizer
// It's important to supply a unique name for purposes of
// tracking. The Transaction object will only manage one
// Finalizer with each name
func NewFinalizer(name string) *Finalizer {
	return &Finalizer{name: name}
}

// Register adds a data modification to the Finalizer
// by registering callback functions for the Finalize,
// Commit and Abort steps
func (m *Finalizer) Register(f, c func() error, a func()) {
	m.finalizers = append(m.finalizers, f)
	m.commits = append(m.commits, c)
	m.aborts = append(m.aborts, a)
}

// Finalize attempts to finalize all data modifications
// it aborts on the first failure with the error from
// that failure
func (m *Finalizer) Finalize() error {
	for _, ff := range m.finalizers {
		if ff != nil {
			err := ff()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Commit executes the commit step on all transactions
func (m *Finalizer) Commit() error {
	for _, cf := range m.commits {
		if cf != nil {
			err := cf()
			if err != nil {
				m.Abort()
				return err
			}
		}
	}
	return nil
}

// Abort this transaction by calling the abort function
// for all data modifications
func (m *Finalizer) Abort() {
	for _, af := range m.aborts {
		if af != nil {
			af()
		}
	}
}
