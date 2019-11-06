package client

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/CanonicalLtd/go-dqlite/internal/protocol"
	_ "github.com/mattn/go-sqlite3" // Go SQLite bindings
)

// NodeStore is used by a dqlite client to get an initial list of candidate
// dqlite nodes that it can dial in order to find a leader dqlite node to use.
type NodeStore = protocol.NodeStore

// NodeInfo holds information about a single server.
type NodeInfo = protocol.NodeInfo

// InmemNodeStore keeps the list of target dqlite nodes in memory.
type InmemNodeStore = protocol.InmemNodeStore

// NewInmemNodeStore creates NodeStore which stores its data in-memory.
var NewInmemNodeStore = protocol.NewInmemNodeStore

// DatabaseNodeStore persists a list addresses of dqlite nodes in a SQL table.
type DatabaseNodeStore struct {
	db     *sql.DB // Database handle to use.
	schema string  // Name of the schema holding the servers table.
	table  string  // Name of the servers table.
	column string  // Column name in the servers table holding the server address.
}

// DefaultNodeStore creates a new NodeStore using the given filename to
// open a SQLite database, with default names for the schema, table and column
// parameters.
//
// It also creates the table if it doesn't exist yet.
func DefaultNodeStore(filename string) (*DatabaseNodeStore, error) {
	// Open the database.
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	// Since we're setting SQLite single-thread mode, we need to have one
	// connection at most.
	db.SetMaxOpenConns(1)

	// Create the servers table if it does not exist yet.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS servers (address TEXT, UNIQUE(address))")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create servers table")
	}

	store := NewNodeStore(db, "main", "servers", "address")

	return store, nil
}

// NewNodeStore creates a new NodeStore.
func NewNodeStore(db *sql.DB, schema, table, column string) *DatabaseNodeStore {
	return &DatabaseNodeStore{
		db:     db,
		schema: schema,
		table:  table,
		column: column,
	}
}

// Get the current servers.
func (d *DatabaseNodeStore) Get(ctx context.Context) ([]NodeInfo, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	query := fmt.Sprintf("SELECT %s FROM %s.%s", d.column, d.schema, d.table)
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query servers table")
	}
	defer rows.Close()

	servers := make([]NodeInfo, 0)
	for rows.Next() {
		var address string
		err := rows.Scan(&address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch server address")
		}
		servers = append(servers, NodeInfo{ID: 1, Address: address})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "result set failure")
	}

	return servers, nil
}

// Set the servers addresses.
func (d *DatabaseNodeStore) Set(ctx context.Context, servers []NodeInfo) error {
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	query := fmt.Sprintf("DELETE FROM %s.%s", d.schema, d.table)
	if _, err := tx.ExecContext(ctx, query); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete existing servers rows")
	}

	query = fmt.Sprintf("INSERT INTO %s.%s(%s) VALUES (?)", d.schema, d.table, d.column)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to prepare insert statement")
	}
	defer stmt.Close()

	for _, server := range servers {
		if _, err := stmt.ExecContext(ctx, server.Address); err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "failed to insert server %s", server.Address)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
