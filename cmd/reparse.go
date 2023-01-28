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
	Use:   "reparse [all|era|deploys|moduleBytes|exceptTransfers]",
	Short: "Reparse all unknown deploys from the database without calling rpc",
	Long: `Reparse specifics items from the database

You must add at least one argument from those :

all: reparse every blocks from rpc, will ignore any other args
era: only reparse switch blocks
deploys: only reparse deploys, will ignore any other deploy args
moduleBytes: only reparse moduleBytes deploys
exceptTransfers: only reparse deploys except transfers deploys
systemPackageContracts: add system Packages Contracts. You must add the network type right after. Ex : reparse systemPackageContracts testnet
`,
	ValidArgs: []string{"all", "era", "deploys", "moduleBytes", "exceptTransfers", "accountPurses", "systemPackageContracts", "testnet", "mainnet"},
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		if arg == "all" {
			reparseAll(getRedisConf(cmd))
			return
		}
		if arg == "era" {
			reparseEraBlocks(getRedisConf(cmd))
		}
		if arg == "systemPackageContracts" {
			if len(args) > 1 {
				network := args[1]
				reparseSystemPackageContracts(getRedisConf(cmd), network)
			}
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
		if arg == "accountPurses" {
			startAccountPurses(getRedisConf(cmd))
			startAccountHashPurses(getRedisConf(cmd))
			startUrefPurses(getRedisConf(cmd))
			startPurses(getRedisConf(cmd))
		}
	},
}

// init the command flags
func init() {
	RootCmd.AddCommand(reparseCmd)
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

// reparseSystemPackageContracts reparse System Package Contracts
func reparseSystemPackageContracts(redis asynq.RedisConnOpt, network string) {
	if network == "testnet" {
		reparseClient = asynq.NewClient(redis)
		defer reparseClient.Close()
		//Handle payment testnet contract
		task, err := tasks.NewContractPackageRawTask("624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "", "")
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("contracts"))
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
	}
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

// startAccountHashPurses reparse account hash purses
func startAccountHashPurses(redis asynq.RedisConnOpt) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()

	findMissingAccountHashesSql := `WITH accounthashes AS (SELECT DISTINCT LOWER(metadata ->> 'target') as accounthash
FROM deploys
WHERE length(LOWER(metadata ->> 'target')) < 66 AND type = 'transfer' AND result is true)
SELECT accountHash from accounthashes
LEFT JOIN accounts ON accounthashes.accounthash = accounts.account_hash WHERE accounts.account_hash IS NULL;`

	rows, err := reparseDatabase.Postgres.Query(context.Background(), findMissingAccountHashesSql)
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
		task, err := tasks.NewAccountHashTask(missing.hash)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("accounts"))
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

// startAccountPurses reparse account purses
func startAccountPurses(redis asynq.RedisConnOpt) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()

	findMissingPublicKeysSql := `WITH keys AS (SELECT DISTINCT LOWER("from") as key
FROM deploys
WHERE ("from" LIKE '01%' AND length("from") = 66) OR ("from" LIKE '02%' AND length("from") = 68)
UNION
SELECT DISTINCT LOWER(metadata ->> 'target') as key
FROM deploys
WHERE length(LOWER(metadata ->> 'target')) >= 66 AND length(LOWER(metadata ->> 'target')) <= 68 AND type = 'transfer' AND result is true)
SELECT key from keys
LEFT JOIN accounts ON keys.key = accounts.public_key WHERE accounts.public_key IS NULL;`

	rows, err := reparseDatabase.Postgres.Query(context.Background(), findMissingPublicKeysSql)
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
		task, err := tasks.NewAccountTask(missing.hash)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("accounts"))
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

// startUrefPurses reparse uref purses
func startUrefPurses(redis asynq.RedisConnOpt) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()

	findMissingPursesSql := `WITH urefs as (WITH uref AS (SELECT jsonb_array_elements(data -> 'Contract' -> 'named_keys') as j
                             from contracts)
               SELECT DISTINCT j ->> 'key' as uref
               FROM uref
               WHERE (j -> 'is_purse')::bool is true
               UNION
               SELECT DISTINCT LOWER(metadata ->> 'target') as uref
               FROM deploys
               WHERE length(LOWER(metadata ->> 'target')) > 68
                 AND type = 'transfer'
                 AND result is true
               UNION
               SELECT DISTINCT LOWER(main_purse) as uref
               FROM accounts)
SELECT uref
from urefs
         LEFT JOIN purses ON urefs.uref = purses.purse
WHERE purses.purse IS NULL;`

	rows, err := reparseDatabase.Postgres.Query(context.Background(), findMissingPursesSql)
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
		task, err := tasks.NewPurseTask(missing.hash)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("accounts"))
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

// startPurses reparse purses
func startPurses(redis asynq.RedisConnOpt) {
	pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
	defer pgPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	reparseDatabase = db.DB{Postgres: pgPool}
	reparseClient = asynq.NewClient(redis)
	defer reparseClient.Close()

	findPursesSql := `SELECT purse from purses;`

	rows, err := reparseDatabase.Postgres.Query(context.Background(), findPursesSql)
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
		task, err := tasks.NewFetchPurseTask(missing.hash)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		_, err = reparseClient.Enqueue(task, asynq.Queue("accounts"))
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
