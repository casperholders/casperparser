// Package cmd define the reparse command
package cmd

import (
	"casperParser/db"
	"casperParser/tasks"
	"context"
	"github.com/hibiken/asynq"
	"log"

	"github.com/spf13/cobra"
)

var reparseDatabase db.DB
var reparseClient *asynq.Client
var reparsePool int

// reparseCmd represents the reparse command
var reparseCmd = &cobra.Command{
	Use:   "reparse",
	Short: "Reparse all unknown deploys from the database without calling rpc",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		startReparse(getRedisConf(cmd))
	},
}

// init the command flags
func init() {
	rootCmd.AddCommand(reparseCmd)
	reparseCmd.Flags().IntVarP(&reparsePool, "pool", "p", 10, "Database connection pool max connections")
}

// startReparse fetch all deploys missing metadata to be parsed. Won't change anything if the config define the deploys type didn't changed
func startReparse(redis asynq.RedisConnOpt) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()
	const sql = `SELECT hash FROM deploys WHERE type = 'moduleBytes';`

	rows, err := reparseDatabase.Postgres.Query(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		missing := struct {
			hash string
		}{}
		err := rows.Scan(&missing.hash)
		if err != nil {
			log.Fatal(err)
		}
		task, err := tasks.NewDeployKnownTask(missing.hash)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("deploys"))
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
	}
	// check rows.Err() after the last rows.Next() :
	if err := rows.Err(); err != nil {
		// on top of errors triggered by bad conditions on the 'rows.Scan()' call,
		// there could also be some bad things like a truncated response because
		// of some network error, etc ...
		log.Fatal(err)
	}
}
