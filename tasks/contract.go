// Package tasks Define the contract task payload and handler
package tasks

import (
	"casperParser/db"
	"casperParser/types/contract"
	"casperParser/types/deploy"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"regexp"
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
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	contractParsed, err := WorkerRpcClient.GetContract(strings.ToLower(p.ContractHash))
	if err != nil {
		return err
	}
	var database = db.DB{Postgres: WorkerPool}
	contractDeploy, err := database.GetDeploy(ctx, p.DeployHash)
	if err != nil {
		return err
	}
	retrieveNamedKeyValues(&contractParsed, contractDeploy)
	if err != nil {
		return err
	}
	contractJsonString, err := json.Marshal(contractParsed.StoredValue)
	if err != nil {
		return err
	}
	contractType, score := contractParsed.GetContractTypeAndScore()
	err = database.InsertContract(ctx, p.ContractHash, strings.ReplaceAll(contractParsed.StoredValue.Contract.ContractPackageHash, "contract-package-wasm", ""), p.DeployHash, p.From, contractType, score, string(contractJsonString))
	if err != nil {
		return err
	}

	return nil
}

func retrieveNamedKeyValues(c *contract.Result, contractDeploy deploy.Result) {
	urefsMap := contractDeploy.MapUrefs()
	accessRights := regexp.MustCompile(`-\d{3}$`)
	urefPrefix := regexp.MustCompile(`^uref-`)
	for index, namedKey := range c.StoredValue.Contract.NamedKeys {
		urefHash := accessRights.ReplaceAllString(namedKey.Key, "")
		balanceHash := urefPrefix.ReplaceAllString(urefHash, "balance-")
		c.StoredValue.Contract.NamedKeys[index].IsPurse = false
		if val, ok := urefsMap[urefHash]; ok {
			c.StoredValue.Contract.NamedKeys[index].InitialValue = val
		}
		if _, balance := urefsMap[balanceHash]; balance {
			c.StoredValue.Contract.NamedKeys[index].IsPurse = true
		}
	}
}

type ContractRawPayload struct {
	ContractHash string
	DeployHash   string
	From         string
}
