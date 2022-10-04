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
	Use:       "reparse",
	Short:     "Reparse all unknown deploys from the database without calling rpc",
	Long:      ``,
	ValidArgs: []string{"all", "era", "deploys", "moduleBytes", "exceptTransfers"},
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			if arg == "all" {
				reparseAll(getRedisConf(cmd))
				return
			}
			if arg == "era" {
				reparseEraBlocks(getRedisConf(cmd))
			}
			if arg == "deploys" {
				reparseDeploys(getRedisConf(cmd))
				return
			}
			if arg == "moduleBytes" {
				reparseModuleBytes(getRedisConf(cmd))
			}
			if arg == "exceptTransfers" {
				reparseExceptTransfers(getRedisConf(cmd))
			}
		}
	},
}

// init the command flags
func init() {
	rootCmd.AddCommand(reparseCmd)
	reparseCmd.Flags().IntVarP(&reparsePool, "pool", "p", 10, "Database connection pool max connections")
}

func reparseAll(redis asynq.RedisConnOpt) {
	const sql = `SELECT height FROM blocks;`
	startReparseBlocks(redis, sql)
}

func reparseEraBlocks(redis asynq.RedisConnOpt) {
	const sql = `SELECT height FROM blocks WHERE era_end is true;`
	startReparseBlocks(redis, sql)
}

func reparseDeploys(redis asynq.RedisConnOpt) {
	const sql = `SELECT hash FROM deploys;`
	startReparseDeploys(redis, sql)
}

func reparseModuleBytes(redis asynq.RedisConnOpt) {
	const sql = `SELECT hash FROM deploys WHERE type = 'moduleBytes';`
	startReparseDeploys(redis, sql)
}

func reparseExceptTransfers(redis asynq.RedisConnOpt) {
	const sql = `SELECT hash FROM deploys WHERE type != 'transfer';`
	startReparseDeploys(redis, sql)
}

// startReparseBlocks reparse blocks for a given sql query
func startReparseBlocks(redis asynq.RedisConnOpt, sql string) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()

	rows, err := reparseDatabase.Postgres.Query(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		missing := struct {
			height int
		}{}
		err := rows.Scan(&missing.height)
		if err != nil {
			log.Fatal(err)
		}
		task, err := tasks.NewBlockRawTask(missing.height)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("blocks"))
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

// startReparseDeploys reparse deploys for a given sql query
func startReparseDeploys(redis asynq.RedisConnOpt, sql string) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()

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
