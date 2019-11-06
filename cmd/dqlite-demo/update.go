package main

import (
	"database/sql"
	"fmt"

	"github.com/CanonicalLtd/go-dqlite/driver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Return a new update key command.
func newUpdate() *cobra.Command {
	var cluster *[]string

	cmd := &cobra.Command{
		Use:   "update <key> <value>",
		Short: "Insert or update a key in the demo table.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]
			store := getStore(*cluster)
			driver, err := driver.New(store, driver.WithLogFunc(logFunc))
			if err != nil {
				return errors.Wrapf(err, "failed to create dqlite driver")
			}
			sql.Register("dqlite", driver)

			db, err := sql.Open("dqlite", "demo.db")
			if err != nil {
				return errors.Wrap(err, "can't open demo database")
			}
			defer db.Close()

			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS model (key TEXT, value TEXT)"); err != nil {
				return errors.Wrap(err, "can't create demo table")
			}

			if _, err := db.Exec("INSERT OR REPLACE INTO model(key, value) VALUES(?, ?)", key, value); err != nil {
				return errors.Wrap(err, "can't update key")
			}

			fmt.Println("done")

			return nil
		},
	}

	flags := cmd.Flags()
	defaultCluster := []string{
		"127.0.0.1:9181",
		"127.0.0.1:9182",
		"127.0.0.1:9183",
	}
	cluster = flags.StringSliceP("cluster", "c", defaultCluster, "addresses of existing cluster nodes")

	return cmd
}
