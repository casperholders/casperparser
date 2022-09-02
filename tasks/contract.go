// Package tasks Define the contract task payload and handler
package tasks

import (
	"casperParser/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"strings"
)

// TypeContractRaw Task contract  raw type
const (
	TypeContractRaw = "contract:raw"
)

// NewContractRawTask Used for not yet parsed contract
func NewContractRawTask(hash string, deployHash string, from string) (*asynq.Task, error) {
	payload, err := json.Marshal(ContractRawPayload{ContractHash: hash, DeployHash: deployHash, From: from})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeContractRaw, payload), nil
}

// HandleContractRawTask fetch a contract  from the rpc endpoint, parse it, and insert it in the database
func HandleContractRawTask(ctx context.Context, t *asynq.Task) error {
	var p ContractRawPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	contract, rawContractHash, err := WorkerRpcClient.GetContract(strings.ToLower(p.ContractHash))
	if err != nil {
		return err
	}

	var database = db.DB{Postgres: WorkerPool}
	err = database.InsertContract(ctx, p.ContractHash, strings.ReplaceAll(contract.StoredValue.Contract.ContractPackageHash, "contract-package-wasm", ""), p.DeployHash, p.From, contract.GetContractType(), rawContractHash)
	if err != nil {
		return err
	}

	return nil
}

type ContractRawPayload struct {
	ContractHash string
	DeployHash   string
	From         string
}
