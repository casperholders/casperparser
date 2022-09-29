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
	switch _, err := db.Postgres.Exec(ctx, sql, hash, era, timestamp, height, eraEnd, false); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create block on database: %v\n", err)
		return errors.New("cannot create block on database")
	}
	return nil
}

// InsertRawBlock in the database
func (db *DB) InsertRawBlock(ctx context.Context, hash string, json string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO raw_blocks ("hash", "data") 
	VALUES ($1, $2)
	ON CONFLICT (hash)
	DO UPDATE
	SET data = $2;`
	switch _, err := db.Postgres.Exec(ctx, sql, hash, json); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create raw block on database: %v\n", err)
		return errors.New("cannot create raw block on database")
	}
	return nil
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
	switch _, err := db.Postgres.Exec(ctx, sql, hash, from, cost, result, timestamp, block, deployType, metadataType, ch, cn, ep, m, e); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create/update deploy on database: %v\n", err)
		return errors.New("cannot create/update deploy on database")
	}
	return nil
}

// InsertRawDeploy in the database
func (db *DB) InsertRawDeploy(ctx context.Context, hash string, json string) error {
	hash = strings.ToLower(hash)
	const sql = `INSERT INTO raw_deploys ("hash", "data")
	VALUES ($1, $2)
	ON CONFLICT (hash)
	DO UPDATE
	SET data = $2;`
	switch _, err := db.Postgres.Exec(ctx, sql, hash, json); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create raw deploy on database: %v\n", err)
		return errors.New("cannot create raw deploy on database")
	}
	return nil
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
	switch _, err := db.Postgres.Exec(ctx, sql, hash, deploy, from, data); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create contract package on database: %v\n", err)
		return errors.New("cannot create contract package on database")
	}
	return nil
}

// InsertContract in the database
func (db *DB) InsertContract(ctx context.Context, hash string, packageHash string, deploy string, from string, contractType string, score float64, data string) error {
	hash = strings.ToLower(hash)
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
	switch _, err := db.Postgres.Exec(ctx, sql, hash, packageHash, deploy, from, contractType, score, data); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create contract on database: %v\n", err)
		return errors.New("cannot create contract on database")
	}
	return nil
}

// InsertReward in the database
func (db *DB) InsertReward(ctx context.Context, block string, era int, delegator_public_key string, validator_public_key string, amount string) error {
	block = strings.ToLower(block)
	const sql = `INSERT INTO rewards ("block", "era", "delegator_public_key", "validator_public_key", "amount")
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (block, era, delegator_public_key, validator_public_key)
	DO UPDATE
	SET amount = $5;`
	var dpk *string
	dpk = nil
	if delegator_public_key != "" {
		dpk = &delegator_public_key
	}
	switch _, err := db.Postgres.Exec(ctx, sql, block, era, dpk, validator_public_key, amount); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create contract on database: %v\n", err)
		return errors.New("cannot create contract on database")
	}
	return nil
}

// GetMissingBlocks from the database
func (db *DB) GetMissingBlocks(ctx context.Context) ([]int, error) {
	const sql = `SELECT all_ids AS missing_ids FROM generate_series((SELECT MIN(height) FROM blocks), (SELECT MAX(height) FROM blocks)) all_ids EXCEPT SELECT height FROM blocks;`
	rows, err := db.Postgres.Query(ctx, sql)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()
	var missingDeploys []int
	for rows.Next() {
		missing := struct {
			id int
		}{}
		err := rows.Scan(&missing.id)
		if err != nil {
			return []int{}, err
		}

		log.Printf("Missing block found: id=%d ", missing.id)
		missingDeploys = append(missingDeploys, missing.id)
	}
	// check rows.Err() after the last rows.Next() :
	if err := rows.Err(); err != nil {
		// on top of errors triggered by bad conditions on the 'rows.Scan()' call,
		// there could also be some bad things like a truncated response because
		// of some network error, etc ...
		return []int{}, err
	}
	return missingDeploys, nil
}

// GetMissingMetadataDeploysHash from the database
func (db *DB) GetMissingMetadataDeploysHash(ctx context.Context) ([]string, error) {
	const sql = `SELECT hash FROM deploys WHERE metadata IS NULL;`

	rows, err := db.Postgres.Query(ctx, sql)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	var missingDeploys []string
	for rows.Next() {
		missing := struct {
			hash string
		}{}
		err := rows.Scan(&missing.hash)
		if err != nil {
			return []string{}, err
		}
		missingDeploys = append(missingDeploys, missing.hash)
	}
	// check rows.Err() after the last rows.Next() :
	if err := rows.Err(); err != nil {
		// on top of errors triggered by bad conditions on the 'rows.Scan()' call,
		// there could also be some bad things like a truncated response because
		// of some network error, etc ...
		return []string{}, err
	}
	return missingDeploys, nil
}

// GetDeploy from the database
func (db *DB) GetDeploy(ctx context.Context, hash string) (deploy.Result, error) {
	const sql = `SELECT data FROM raw_deploys WHERE hash = $1;`
	rows, err := db.Postgres.Query(ctx, sql, hash)
	if err != nil {
		log.Println(err)
		return deploy.Result{}, err
	}
	defer rows.Close()
	var d deploy.Result
	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			log.Println(err)
			return deploy.Result{}, err
		}
	}
	return d, nil
}

// GetRawBlock from the database
func (db *DB) GetRawBlock(ctx context.Context, hash string) (block.Result, error) {
	const sql = `SELECT data FROM raw_blocks WHERE hash = $1;`
	rows, err := db.Postgres.Query(ctx, sql, hash)
	if err != nil {
		log.Println(err)
		return block.Result{}, err
	}
	defer rows.Close()
	var d block.Result
	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			log.Println(err)
			return block.Result{}, err
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
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()
	var d int
	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	}
	return d, nil
}

// ValidateBlock from the database
func (db *DB) ValidateBlock(ctx context.Context, hash string) error {
	const sql = `UPDATE blocks SET validated = true WHERE hash = $1;`
	switch _, err := db.Postgres.Exec(ctx, sql, hash); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.blockPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot validate block on database: %v\n", err)
		return errors.New("cannot validate block on database")
	}
	return nil
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
