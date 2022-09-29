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

// TypeReward Task reward type
const (
	TypeReward = "reward:raw"
)

// NewRewardTask Used for reward
func NewRewardTask(hash string) (*asynq.Task, error) {
	payload, err := json.Marshal(RewardPayload{BlockHash: hash})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReward, payload), nil
}

// HandleRewardTask fetch a contract  from the rpc endpoint, parse it, and insert it in the database
func HandleRewardTask(ctx context.Context, t *asynq.Task) error {
	var p RewardPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	eraParsed, err := WorkerRpcClient.GetEraInfo(strings.ToLower(p.BlockHash))
	if err != nil {
		return err
	}
	var database = db.DB{Postgres: WorkerPool}
	for _, s := range eraParsed.EraSummary.StoredValue.EraInfo.SeigniorageAllocations {
		dpk := ""
		vpk := ""
		amount := ""
		if s.Delegator != nil {
			dpk = s.Delegator.DelegatorPublicKey
			vpk = s.Delegator.ValidatorPublicKey
			amount = s.Delegator.Amount
		}
		if s.Validator != nil {
			vpk = s.Validator.ValidatorPublicKey
			amount = s.Validator.Amount
		}
		err = database.InsertReward(ctx, eraParsed.EraSummary.BlockHash, eraParsed.EraSummary.EraId, dpk, vpk, amount)
		if err != nil {
			return err
		}
	}

	return nil
}

type RewardPayload struct {
	BlockHash string
}
