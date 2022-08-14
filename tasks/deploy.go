// Package tasks Define the deploy task payload and handler
package tasks

import (
	"casperParser/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"strings"
)

// TypeDeployRaw Task deploy raw type
// TypeDeployKnown Task deploy known type
const (
	TypeDeployRaw   = "deploy:raw"
	TypeDeployKnown = "deploy:known"
)

// NewDeployRawTask Used for not yet parsed deploy
func NewDeployRawTask(hash string) (*asynq.Task, error) {
	payload, err := json.Marshal(DeployRawPayload{DeployHash: hash})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDeployRaw, payload), nil
}

// NewDeployKnownTask used for already parsed deploy
func NewDeployKnownTask(hash string) (*asynq.Task, error) {
	payload, err := json.Marshal(DeployKnownPayload{DeployHash: hash})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDeployKnown, payload), nil
}

// HandleDeployRawTask fetch a deploy from the rpc endpoint, parse it, and insert it in the database
func HandleDeployRawTask(ctx context.Context, t *asynq.Task) error {
	var p DeployRawPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	rpcDeploy, resp, err := WorkerRpcClient.GetDeploy(p.DeployHash)
	if err != nil {
		return err
	}

	result, cost, err := rpcDeploy.GetResultAndCost()
	if err != nil {
		return err
	}
	metadataDeployType, metadata := rpcDeploy.GetDeployMetadata()
	events := rpcDeploy.GetEvents()
	jsonString := strings.ReplaceAll(string(resp), "\\u0000", "")
	var database = db.DB{Postgres: WorkerPool}
	err = database.InsertDeploy(ctx, rpcDeploy.Deploy.Hash, rpcDeploy.Deploy.Header.Account, cost, result, rpcDeploy.Deploy.Header.Timestamp, rpcDeploy.ExecutionResults[0].BlockHash, rpcDeploy.GetType(), jsonString, metadataDeployType, metadata, events)
	if err != nil {
		return err
	}

	return nil
}

// HandleDeployKnownTask fetch a deploy from the database, parse it, and insert it in the database
func HandleDeployKnownTask(ctx context.Context, t *asynq.Task) error {
	var p DeployKnownPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	var database = db.DB{Postgres: WorkerPool}
	dbDeploy, err := database.GetDeploy(ctx, p.DeployHash)
	if err != nil {
		log.Printf("Can't find deploy %s\n", p.DeployHash)
		return err
	}

	result, cost, err := dbDeploy.GetResultAndCost()
	if err != nil {
		return err
	}
	metadataDeployType, metadata := dbDeploy.GetDeployMetadata()
	events := dbDeploy.GetEvents()
	if metadata != "" {
		log.Printf("New metadata found for %s of type : %s\n", p.DeployHash, metadataDeployType)
		err = database.UpdateDeploy(ctx, dbDeploy.Deploy.Hash, dbDeploy.Deploy.Header.Account, cost, result, dbDeploy.Deploy.Header.Timestamp, dbDeploy.ExecutionResults[0].BlockHash, dbDeploy.GetType(), metadataDeployType, metadata, events)
		if err != nil {
			return err
		}
	}
	return nil
}

type DeployRawPayload struct {
	DeployHash string
}

type DeployKnownPayload struct {
	DeployHash string
}
