package txmanager

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// MakePostgresTxFinalizer is a constructor for a Postgres
// transaction driver
func MakePostgresTxFinalizer(
	ctx context.Context, name string, cPool *sql.DB,
) *PostgresTxFinalizer {
	tx, err := cPool.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	return &PostgresTxFinalizer{name: name, TX: tx}
}

// PostgresTxFinalizer manages transactions on a PostgreSQL
// server
type PostgresTxFinalizer struct {
	name string
	TX   *sql.Tx
	id   string
}

// Finalize sets up a prepared transaction
func (m *PostgresTxFinalizer) Finalize() error {
	m.id = uuid.New().String()
	_, err := m.TX.Exec(fmt.Sprintf("PREPARE TRANSACTION '%s'", m.id))
	return err
}

// Commit finishes the transaction
func (m *PostgresTxFinalizer) Commit() {
	_, err := m.TX.Exec(fmt.Sprintf("COMMIT PREPARED '%s'", m.id))
	if err != nil {
		panic(err)
	}
}

// Abort rolls back the transaction
func (m *PostgresTxFinalizer) Abort() {
	_, err := m.TX.Exec(fmt.Sprintf("ROLLBACK PREPARED '%s'", m.id))
	if err != nil {
		panic(err)
	}
}
