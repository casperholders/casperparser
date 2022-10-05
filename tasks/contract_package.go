// Package tasks Define the contract package task payload and handler
package tasks

import (
	"casperParser/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"strings"
)

// TypeContractPackageRaw Task contract package raw type
const (
	TypeContractPackageRaw = "contract_package:raw"
)

// NewContractPackageRawTask Used for not yet parsed contract package
func NewContractPackageRawTask(hash string, deployHash string, from string) (*asynq.Task, error) {
	payload, err := json.Marshal(ContractPackageRawPayload{ContractPackageHash: hash, DeployHash: deployHash, From: from})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeContractPackageRaw, payload), nil
}

// HandleContractPackageRawTask fetch a contract package from the rpc endpoint, parse it, and insert it in the database
func HandleContractPackageRawTask(ctx context.Context, t *asynq.Task) error {
	var p ContractPackageRawPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	rawContractPackageHash, err := WorkerRpcClient.GetContractPackage(strings.ToLower(p.ContractPackageHash))
	if err != nil {
		return err
	}

	var database = db.DB{Postgres: WorkerPool}
	err = database.InsertContractPackage(ctx, p.ContractPackageHash, p.DeployHash, p.From, rawContractPackageHash)
	if err != nil {
		return err
	}

	return nil
}

type ContractPackageRawPayload struct {
	ContractPackageHash string
	DeployHash          string
	From                string
}
