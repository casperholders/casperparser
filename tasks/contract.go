// Package tasks Define the contract task payload and handler
package tasks

import (
	"casperParser/db"
	"casperParser/types/contract"
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
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	contractParsed, err := WorkerRpcClient.GetContract(strings.ToLower(p.ContractHash))
	if err != nil {
		return err
	}
	var database = db.DB{Postgres: WorkerPool}
	namedKeys := retrieveNamedKeyValues(contractParsed)
	contractParsed.StoredValue.Contract.NamedKeys = []contract.NamedKey{}
	contractJsonString, err := json.Marshal(contractParsed.StoredValue)
	if err != nil {
		return err
	}
	contractType, score := contractParsed.GetContractTypeAndScore()
	err = database.InsertContract(ctx, p.ContractHash, strings.ReplaceAll(contractParsed.StoredValue.Contract.ContractPackageHash, "contract-package-wasm", ""), p.DeployHash, p.From, contractType, score, string(contractJsonString))
	if err != nil {
		return err
	}
	for _, namedKey := range namedKeys {
		err = database.InsertNamedKey(ctx, namedKey.Uref, namedKey.Name, namedKey.IsPurse, namedKey.InitialValue, p.ContractHash)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveNamedKeyValues(c contract.Result) []NamedKey {
	var namedKeys []NamedKey
	for _, namedKey := range c.StoredValue.Contract.NamedKeys {
		if strings.Contains(namedKey.Key, "account-hash-") {
			namedKeys = append(namedKeys, NamedKey{
				Uref:         namedKey.Key,
				Name:         namedKey.Name,
				IsPurse:      false,
				InitialValue: "null",
			})
		} else {
			value, isPurse, _ := WorkerRpcClient.GetUrefValue(namedKey.Key)
			namedKeys = append(namedKeys, NamedKey{
				Uref:         namedKey.Key,
				Name:         namedKey.Name,
				IsPurse:      isPurse,
				InitialValue: value,
			})
		}
	}
	return namedKeys
}

type NamedKey struct {
	Uref         string
	Name         string
	IsPurse      bool
	InitialValue string
}

type ContractRawPayload struct {
	ContractHash string
	DeployHash   string
	From         string
}
