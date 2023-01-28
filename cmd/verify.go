// Package cmd define the verify command
package cmd

import (
	"casperParser/db"
	"casperParser/tasks"
	"context"
	"github.com/hibiken/asynq"
	"log"

	"github.com/spf13/cobra"
)

var verifyDatabase db.DB
var verifyClient *asynq.Client
var verifyPool int

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that all deploys are present in the database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		startVerify(getRedisConf(cmd))
	},
}

// init the command flags
func init() {
	RootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().IntVarP(&verifyPool, "pool", "p", 10, "Database connection pool max connections")
}

// startVerify fetch all blocks not validated and check if all deploys are present in the db
func startVerify(redis asynq.RedisConnOpt) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	verifyDatabase = db.DB{Postgres: pgPool}
	verifyClient = asynq.NewClient(redis)
	defer verifyClient.Close()
	const sql = `SELECT hash FROM blocks WHERE validated = false;`

	rows, err := verifyDatabase.Postgres.Query(context.Background(), sql)
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
		task, err := tasks.NewBlockVerifyTask(missing.hash)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = verifyClient.Enqueue(task, asynq.Queue("blocks"))
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
