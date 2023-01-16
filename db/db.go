// Package db wrapper around the PGX library. Used to insert blocks and deploys and query some data easily
package db

import (
	"casperParser/types/block"
	"casperParser/types/deploy"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"strings"
)

type DB struct {
	// Postgres database.PGX
	Postgres *pgxpool.Pool
}

// NewPGXPool create a postgresql database connexion pool
func NewPGXPool(ctx context.Context, connString string, maxConnection int) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(connString) // Using environment variables instead of a connection string.
	if err != nil {
		return nil, err
	}

	// pgx, by default, does some I/O operation on initialization of a pool to check if the database is reachable.
	// Comment the following line if you don't want pgx to try to connect pool once the Connect function is called,
	//
	// If you comment it, and your application seems stuck, you probably forgot to set up PGCONNECT_TIMEOUT,
	// and your code is hanging waiting for a connection to be established.
	conf.LazyConnect = true

	// pgxpool default max number of connections is the number of CPUs on your machine returned by runtime.NumCPU().
	// This number is very conservative, and you might be able to improve performance for highly concurrent applications
	// by increasing it.
	conf.MaxConns = int32(maxConnection)

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	return pool, nil
}

// InsertBlock in the database
func (db *DB) InsertBlock(ctx context.Context, hash string, era int, timestamp string, height int, eraEnd bool, json string) error {
	hash = strings.ToLower(hash)
	err := db.InsertRawBlock(ctx, hash, json)
	if err != nil {
		return err
	}
	const sql = `INSERT INTO blocks ("hash", "era", "timestamp", "height", "era_end", "validated") 
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (hash)
	DO UPDATE
	SET era = $2,
	timestamp = $3,
	height = $4,
	era_end = $5,
	validated = $6;`
	_, err = db.Postgres.Exec(ctx, sql, hash, era, timestamp, height, eraEnd, false)
	return db.checkErr(err)
}

// InsertRawBlock in the database
func (db *DB) InsertRawBlock(ctx context.Context, hash string, json string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO raw_blocks ("hash", "data") 
	VALUES ($1, $2)
	ON CONFLICT (hash)
	DO UPDATE
	SET data = $2;`
	_, err := db.Postgres.Exec(ctx, sql, hash, json)
	return db.checkErr(err)
}

// InsertDeploy in the database
func (db *DB) InsertDeploy(ctx context.Context, hash string, from string, cost string, result bool, timestamp string, block string, deployType string, json string, metadataType string, contractHash string, contractName string, entrypoint string, metadata string, events string) error {
	hash = strings.ToLower(hash)
	err := db.InsertRawDeploy(ctx, hash, json)
	if err != nil {
		return err
	}
	return db.UpdateDeploy(ctx, hash, from, cost, result, timestamp, block, deployType, metadataType, contractHash, contractName, entrypoint, metadata, events)
}

func (db *DB) InsertAuction(ctx context.Context, rowsToInsertBids [][]interface{}, rowsToInsertDelegators [][]interface{}) error {
	const sql = `TRUNCATE bids; TRUNCATE delegators;`
	_, err := db.Postgres.Exec(ctx, sql)
	if db.checkErr(err) != nil {
		return db.checkErr(err)
	}

	_, err = db.Postgres.CopyFrom(
		ctx,
		pgx.Identifier{"bids"},
		[]string{"public_key", "bonding_purse", "staked_amount", "delegation_rate", "inactive"},
		pgx.CopyFromRows(rowsToInsertBids),
	)
	if db.checkErr(err) != nil {
		return db.checkErr(err)
	}
	_, err = db.Postgres.CopyFrom(
		ctx,
		pgx.Identifier{"delegators"},
		[]string{"public_key", "delegatee", "staked_amount", "bonding_purse"},
		pgx.CopyFromRows(rowsToInsertDelegators),
	)
	return db.checkErr(err)
}

// UpdateDeploy in the database
func (db *DB) UpdateDeploy(ctx context.Context, hash string, from string, cost string, result bool, timestamp string, block string, deployType string, metadataType string, contractHash string, contractName string, entrypoint string, metadata string, events string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO deploys ("hash", "from", "cost", "result", "timestamp", "block", "type", "metadata_type", "contract_hash", "contract_name", "entrypoint", "metadata", "events")
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT (hash)
	DO UPDATE
	SET "from" = $2,
	cost = $3,
	result = $4,
	timestamp = $5,
	block = $6,
	type = $7,
	metadata_type = $8,
	contract_hash = $9,
	contract_name = $10,
	entrypoint = $11,
	metadata = $12,
	events = $13;`
	var m *string
	m = nil
	if metadata != "" {
		m = &metadata
	}
	var e *string
	e = nil
	if events != "" {
		e = &events
	}
	var ch *string
	ch = nil
	if contractHash != "" {
		ch = &contractHash
	}
	var cn *string
	cn = nil
	if contractName != "" {
		cn = &contractName
	}
	var ep *string
	ep = nil
	if entrypoint != "" {
		ep = &entrypoint
	}
	_, err := db.Postgres.Exec(ctx, sql, hash, from, cost, result, timestamp, block, deployType, metadataType, ch, cn, ep, m, e)
	return db.checkErr(err)
}

// InsertRawDeploy in the database
func (db *DB) InsertRawDeploy(ctx context.Context, hash string, json string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO raw_deploys ("hash", "data")
	VALUES ($1, $2)
	ON CONFLICT (hash)
	DO UPDATE
	SET data = $2;`
	_, err := db.Postgres.Exec(ctx, sql, hash, json)
	return db.checkErr(err)
}

// InsertContractPackage in the database
func (db *DB) InsertContractPackage(ctx context.Context, hash string, deploy string, from string, data string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO contract_packages ("hash", "deploy", "from", "data")
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (hash)
	DO UPDATE
	SET deploy = $2,
	"from" = $3,
	data = $4;`
	_, err := db.Postgres.Exec(ctx, sql, hash, deploy, from, data)
	return db.checkErr(err)
}

// InsertContract in the database
func (db *DB) InsertContract(ctx context.Context, hash string, packageHash string, deploy string, from string, contractType string, score float64, data string) error {
	hash = strings.ToLower(hash)
	packageHash = strings.ToLower(packageHash)
	const sql = `INSERT INTO contracts ("hash", "package", "deploy", "from", "type", "score", "data")
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (hash)
	DO UPDATE
	SET package = $2,
	deploy = $3,    
	"from" = $4,
	type = $5,
	score = $6,
	data = $7;`
	_, err := db.Postgres.Exec(ctx, sql, hash, packageHash, deploy, from, contractType, score, data)
	if err != nil {
		log.Printf(data)
	}
	return db.checkErr(err)
}

// InsertNamedKey in the database
func (db *DB) InsertNamedKey(ctx context.Context, uref string, name string, isPurse bool, initialValue string, contractHash string) error {
	contractHash = strings.ToLower(contractHash)
	const sql = `INSERT INTO named_keys ("uref", "name", "is_purse", "initial_value")
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (uref)
	DO UPDATE
	SET name = $2,
	is_purse = $3,    
	initial_value = $4;`
	_, err := db.Postgres.Exec(ctx, sql, uref, name, isPurse, initialValue)
	if err != nil {
		log.Printf("Uref: %s Name: %s Initial Value: %s \n", uref, name, initialValue)
		return db.checkErr(err)
	}
	const joinsql = `INSERT INTO contracts_named_keys (contract_hash, named_key_uref) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	_, err = db.Postgres.Exec(ctx, joinsql, contractHash, uref)
	if err != nil {
		log.Printf("Uref: %s Name: %s Initial Value: %s \n", uref, name, initialValue)
	}
	return db.checkErr(err)
}

// InsertAccountHash in the database
func (db *DB) InsertAccountHash(ctx context.Context, hash string, purse string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO accounts ("public_key", "account_hash", "main_purse")
	VALUES ($1, $2, $3)
	ON CONFLICT (account_hash)
	DO UPDATE
	SET public_key = $1,
	main_purse = $3;`
	_, err := db.Postgres.Exec(ctx, sql, nil, hash, purse)
	return db.checkErr(err)
}

// InsertAccount in the database
func (db *DB) InsertAccount(ctx context.Context, publicKey string, hash string, purse string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO accounts ("public_key", "account_hash", "main_purse")
	VALUES ($1, $2, $3)
	ON CONFLICT (account_hash)
	DO UPDATE
	SET public_key = $1,
	main_purse = $3;`
	_, err := db.Postgres.Exec(ctx, sql, publicKey, hash, purse)
	return db.checkErr(err)
}

// InsertPurse in the database
func (db *DB) InsertPurse(ctx context.Context, hash string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO purses ("purse", "balance")
	VALUES ($1, $2)
	ON CONFLICT (purse)
	DO UPDATE
	SET balance = $2;`
	_, err := db.Postgres.Exec(ctx, sql, hash, nil)
	return db.checkErr(err)
}

// InsertPurseBalance in the database
func (db *DB) InsertPurseBalance(ctx context.Context, hash string, balance string) error {
	hash = strings.ToLower(hash)

	const sql = `INSERT INTO purses ("purse", "balance")
	VALUES ($1, $2)
	ON CONFLICT (purse)
	DO UPDATE
	SET balance = $2;`
	_, err := db.Postgres.Exec(ctx, sql, hash, balance)
	return db.checkErr(err)
}

// InsertRewards in the database
func (db *DB) InsertRewards(ctx context.Context, rowsToInsert [][]interface{}) error {
	count, err := db.Postgres.CopyFrom(
		ctx,
		pgx.Identifier{"rewards"},
		[]string{"block", "era", "delegator_public_key", "validator_public_key", "amount"},
		pgx.CopyFromRows(rowsToInsert),
	)
	if err != nil || count < int64(len(rowsToInsert)) {
		const sql = `DELETE FROM rewards WHERE block = $1;`
		_, errD := db.Postgres.Query(ctx, sql, rowsToInsert[0][0])
		if db.checkErr(errD) != nil {
			return db.checkErr(errD)
		}
	}
	return db.checkErr(err)
}

// GetMissingBlocks from the database
func (db *DB) GetMissingBlocks(ctx context.Context) ([]int, error) {
	const sql = `SELECT all_ids AS missing_ids FROM generate_series((SELECT MIN(height) FROM blocks), (SELECT MAX(height) FROM blocks)) all_ids EXCEPT SELECT height FROM blocks;`
	rows, err := db.Postgres.Query(ctx, sql)
	if db.checkErr(err) != nil {
		return []int{}, db.checkErr(err)
	}
	defer rows.Close()
	var missingDeploys []int
	for rows.Next() {
		missing := struct {
			id int
		}{}
		err = rows.Scan(&missing.id)
		if db.checkErr(err) != nil {
			return []int{}, db.checkErr(err)
		}

		log.Printf("Missing block found: id=%d ", missing.id)
		missingDeploys = append(missingDeploys, missing.id)
	}
	// check rows.Err() after the last rows.Next() :
	if err = rows.Err(); db.checkErr(err) != nil {
		return []int{}, db.checkErr(err)
		// on top of errors triggered by bad conditions on the 'rows.Scan()' call,
		// there could also be some bad things like a truncated response because
		// of some network error, etc ...
	}
	return missingDeploys, nil
}

// GetMissingMetadataDeploysHash from the database
func (db *DB) GetMissingMetadataDeploysHash(ctx context.Context) ([]string, error) {
	const sql = `SELECT hash FROM deploys WHERE metadata IS NULL;`

	rows, err := db.Postgres.Query(ctx, sql)
	if db.checkErr(err) != nil {
		return []string{}, db.checkErr(err)
	}
	defer rows.Close()
	var missingDeploys []string
	for rows.Next() {
		missing := struct {
			hash string
		}{}
		err = rows.Scan(&missing.hash)
		if db.checkErr(err) != nil {
			return []string{}, db.checkErr(err)
		}
		missingDeploys = append(missingDeploys, missing.hash)
	}
	// check rows.Err() after the last rows.Next() :
	if err = rows.Err(); db.checkErr(err) != nil {
		return []string{}, db.checkErr(err)
		// on top of errors triggered by bad conditions on the 'rows.Scan()' call,
		// there could also be some bad things like a truncated response because
		// of some network error, etc ...
	}
	return missingDeploys, nil
}

// GetDeploy from the database
func (db *DB) GetDeploy(ctx context.Context, hash string) (deploy.Result, error) {
	const sql = `SELECT data FROM raw_deploys WHERE hash = $1;`
	rows, err := db.Postgres.Query(ctx, sql, hash)
	if db.checkErr(err) != nil {
		return deploy.Result{}, db.checkErr(err)
	}
	defer rows.Close()
	var d deploy.Result
	for rows.Next() {
		err = rows.Scan(&d)
		if db.checkErr(err) != nil {
			return deploy.Result{}, db.checkErr(err)
		}
	}
	return d, nil
}

// GetRawBlock from the database
func (db *DB) GetRawBlock(ctx context.Context, hash string) (block.Result, error) {
	const sql = `SELECT data FROM raw_blocks WHERE hash = $1;`
	rows, err := db.Postgres.Query(ctx, sql, hash)
	if db.checkErr(err) != nil {
		return block.Result{}, db.checkErr(err)
	}
	defer rows.Close()
	var d block.Result
	for rows.Next() {
		err = rows.Scan(&d)
		if db.checkErr(err) != nil {
			return block.Result{}, db.checkErr(err)
		}
	}
	return d, nil
}

// CountDeploys from the database
func (db *DB) CountDeploys(ctx context.Context, hashes []string) (int, error) {
	var paramrefs string
	for i := range hashes {
		paramrefs += `$` + strconv.Itoa(i+1) + `,`
	}
	paramrefs = paramrefs[:len(paramrefs)-1] // remove last ","
	sql := `SELECT count(*) FROM deploys WHERE hash IN (` + paramrefs + `)`
	genericHashes := make([]interface{}, len(hashes))
	for i, v := range hashes {
		genericHashes[i] = v
	}
	rows, err := db.Postgres.Query(ctx, sql, genericHashes...)
	if db.checkErr(err) != nil {
		return 0, db.checkErr(err)
	}
	defer rows.Close()
	var d int
	for rows.Next() {
		err = rows.Scan(&d)
		if db.checkErr(err) != nil {
			return 0, db.checkErr(err)
		}
	}
	return d, nil
}

// ValidateBlock from the database
func (db *DB) ValidateBlock(ctx context.Context, hash string) error {
	const sql = `UPDATE blocks SET validated = true WHERE hash = $1;`
	_, err := db.Postgres.Exec(ctx, sql, hash)
	return db.checkErr(err)
}

func (db *DB) blockPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}
	if pgErr.Code == pgerrcode.UniqueViolation {
		return errors.New("product already exists")
	}
	if pgErr.Code == pgerrcode.CheckViolation {
		switch pgErr.ConstraintName {
		case "product_id_check":
			return errors.New("invalid product ID")
		case "product_name_check":
			return errors.New("invalid product name")
		case "product_price_check":
			return errors.New("invalid price")
		}
	}
	return nil
}

type PGXStdLogger struct{}

func (l *PGXStdLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

func (db *DB) checkErr(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	if err != nil {
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("db error : %v\n", err)
		return errors.New("db error")
	}
	return nil
}
