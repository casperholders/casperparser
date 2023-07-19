package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewDeployRawTask(t *testing.T) {
	task, err := NewDeployRawTask("test")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	if task.Type() != "deploy:raw" {
		t.Errorf("NewBlockRawTask has a bad name. Received : %s. Expected : %s", task.Type(), "block:raw")
	}
}

func TestHandleDeployKnownTask(t *testing.T) {
	task, err := NewDeployKnownTask("test")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	if task.Type() != "deploy:known" {
		t.Errorf("NewBlockRawTask has a bad name. Received : %s. Expected : %s", task.Type(), "block:raw")
	}
}

func TestHandleDeployRawTask(t *testing.T) {
	dbconstring := os.Getenv("CASPER_PARSER_DATABASE")
	redis := os.Getenv("CASPER_PARSER_REDIS")
	redisConf := asynq.RedisClientOpt{
		Addr: redis,
	}
	rpcendpointurl := os.Getenv("CASPER_PARSER_RPC")
	WorkerPool, _ = db.NewPGXPool(context.Background(), dbconstring, 10)
	WorkerRpcClient = rpc.NewRpcClient(rpcendpointurl)
	WorkerAsyncClient = asynq.NewClient(redisConf)
	defer WorkerAsyncClient.Close()
	defer WorkerPool.Close()
	task, err := NewBlockRawTask(64)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockVerifyTask("96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockVerifyTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockRawTask(981072)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockVerifyTask("b7c2082f9af73771b6e360cb762f80a0177235f05fbae52f312d0471d2678bb2")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockVerifyTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewDeployRawTask("00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleDeployRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewDeployKnownTask("00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleDeployKnownTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewDeployRawTask("03eb82b2e02c5880cd03fcc75580505571c69d476ce28d6cdbb0ee1930cf5950")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleDeployRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewDeployKnownTask("03eb82b2e02c5880cd03fcc75580505571c69d476ce28d6cdbb0ee1930cf5950")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleDeployKnownTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
}
