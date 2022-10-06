package tasks

import (
	"casperParser/db"
	"casperParser/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"strings"
)

// TypeAccountHash Task account hash raw insert
// TypeAccountPublicKey Task account public key raw insert
// TypeAccountUref Task uref purse raw insert
// TypeAccountFetch Task fetch uref balance
const (
	TypeAccountHash      = "account:hash"
	TypeAccountPublicKey = "account:publickey"
	TypeAccountUref      = "account:uref"
	TypeAccountFetch     = "account:fetch"
)

// NewAccountHashTask used to create account from account hash
func NewAccountHashTask(accountHash string) (*asynq.Task, error) {
	payload, err := json.Marshal(Hash{Hash: accountHash})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAccountHash, payload), nil
}

// NewAccountTask used to create account
func NewAccountTask(publickey string) (*asynq.Task, error) {
	payload, err := json.Marshal(Hash{Hash: publickey})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAccountPublicKey, payload), nil
}

// NewPurseTask used create purse
func NewPurseTask(purse string) (*asynq.Task, error) {
	payload, err := json.Marshal(Hash{
		Hash: purse,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAccountUref, payload), nil
}

// NewFetchPurseTask used to fetch purse
func NewFetchPurseTask(purse string) (*asynq.Task, error) {
	payload, err := json.Marshal(Hash{
		Hash: purse,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAccountFetch, payload), nil
}

// HandleAccountHashTask fetch account hash main purse from the rpc endpoint, parse it, and insert it in the database
func HandleAccountHashTask(ctx context.Context, t *asynq.Task) error {
	var p Hash
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}
	purse, err := WorkerRpcClient.GetMainPurse("account-hash-" + p.Hash)
	if err != nil {
		return err
	}
	log.Printf("Hash: %s Purse: %s \n", p.Hash, purse)
	var database = db.DB{Postgres: WorkerPool}

	err = database.InsertAccountHash(ctx, p.Hash, purse)
	if err != nil {
		return err
	}

	addFetchPurseToQueue(purse)
	return nil
}

// HandleAccountTask fetch main purse from the rpc endpoint, parse it, and insert it in the database
func HandleAccountTask(ctx context.Context, t *asynq.Task) error {
	var p Hash
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	accountHash := utils.AccountHash(p.Hash)

	if accountHash == "" {
		return fmt.Errorf("unable to convert public key : %s into account hash", p.Hash)
	}

	purse, err := WorkerRpcClient.GetMainPurse("account-hash-" + accountHash)
	if err != nil {
		if strings.Contains(err.Error(), "ValueNotFound(\"Failed to find base key at path") {
			return fmt.Errorf("failed to retrieve account : %v: %w", err, asynq.SkipRetry)
		}
		return err
	}
	log.Printf("Hash: %s AccountHash: %s Purse: %s \n", p.Hash, accountHash, purse)
	var database = db.DB{Postgres: WorkerPool}

	err = database.InsertAccount(ctx, p.Hash, accountHash, purse)
	if err != nil {
		return err
	}

	addFetchPurseToQueue(purse)
	return nil
}

// HandlePurseTask insert purse in the database
func HandlePurseTask(ctx context.Context, t *asynq.Task) error {
	var p Hash
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	var database = db.DB{Postgres: WorkerPool}

	log.Printf("Purse: %s \n", p.Hash)
	err := database.InsertPurse(ctx, p.Hash)
	if err != nil {
		return err
	}
	addFetchPurseToQueue(p.Hash)
	return nil
}

// addFetchPurseToQueue
func addFetchPurseToQueue(hash string) {
	task, err := NewFetchPurseTask(hash)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	_, err = WorkerAsyncClient.Enqueue(task, asynq.Queue("accounts"))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
}

// HandleFetchPurseTask fetch purse from the rpc endpoint, parse it, and insert it in the database
func HandleFetchPurseTask(ctx context.Context, t *asynq.Task) error {
	var p Hash
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	balance, err := WorkerRpcClient.GetPurseBalance(p.Hash)

	var database = db.DB{Postgres: WorkerPool}

	log.Printf("Purse: %s Balance: %s \n", p.Hash, balance)
	err = database.InsertPurseBalance(ctx, p.Hash, balance)
	if err != nil {
		return err
	}

	return nil
}

type Hash struct {
	Hash string
}
