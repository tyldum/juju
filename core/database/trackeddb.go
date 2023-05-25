// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package database

import (
	"context"
	"database/sql"

	"github.com/canonical/sqlair"
)

// TrackedDB defines an interface for keeping track of sql.DB. This is useful
// knowing if the underlying DB can be reused after an error has occurred.
type TrackedDB interface {
	// Txn executes the input function against the tracked database, using
	// the sqlair package. The sqlair package provides a mapping library for
	// SQL queries and statements.
	// Retry semantics are applied automatically based on transient failures.
	// This is the function that almost all downstream database consumers
	// should use.
	Txn(context.Context, func(context.Context, *sqlair.TX) error) error

	// StdTxn executes the input function against the tracked database,
	// within a transaction that depends on the input context.
	// Retry semantics are applied automatically based on transient failures.
	// This is the function that almost all downstream database consumers
	// should use.
	StdTxn(context.Context, func(context.Context, *sql.Tx) error) error
}
