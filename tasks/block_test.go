package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewBlockRawTask(t *testing.T) {
	task, err := NewBlockRawTask(1)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	if task.Type() != "block:raw" {
		t.Errorf("NewBlockRawTask has a bad name. Received : %s. Expected : %s", task.Type(), "block:raw")
	}
}

func TestNewBlockVerifyTask(t *testing.T) {
	task, err := NewBlockVerifyTask("1")
	if err != nil {
		t.Errorf("Unable to create a NewBlockVerifyTask : %s", err)
	}
	if task.Type() != "block:verify" {
		t.Errorf("NewBlockVerifyTask has a bad name. Received : %s. Expected : %s", task.Type(), "block:verify")
	}
}

func TestHandleBlockRawTask(t *testing.T) {
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
	task, err := NewBlockRawTask(1894586)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockVerifyTask("21fd6475128c11e71b45d4c88fa7f251cfec18ac2a481f39b8f88c405b140754")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockVerifyTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockRawTask(1894357)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockVerifyTask("8d6a98a977482af4eb308cfb4ebdf6981643afdc06f56d6589792808992f56fe")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockVerifyTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
}
