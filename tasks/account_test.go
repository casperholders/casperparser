package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewAccountHashTask(t *testing.T) {
	task, err := NewAccountHashTask("test")
	if err != nil {
		t.Errorf("Unable to create a NewAccountHashTask : %s", err)
	}
	if task.Type() != TypeAccountHash {
		t.Errorf("NewAccountHashTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeAccountHash)
	}
}

func TestNewAccountTask(t *testing.T) {
	task, err := NewAccountTask("test")
	if err != nil {
		t.Errorf("Unable to create a NewAccountTask : %s", err)
	}
	if task.Type() != TypeAccountPublicKey {
		t.Errorf("NewAccountTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeAccountPublicKey)
	}
}

func TestNewPurseTask(t *testing.T) {
	task, err := NewPurseTask("test")
	if err != nil {
		t.Errorf("Unable to create a NewPurseTask : %s", err)
	}
	if task.Type() != TypeAccountUref {
		t.Errorf("NewPurseTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeAccountUref)
	}
}

func TestNewFetchPurseTask(t *testing.T) {
	task, err := NewFetchPurseTask("test")
	if err != nil {
		t.Errorf("Unable to create a NewFetchPurseTask : %s", err)
	}
	if task.Type() != TypeAccountFetch {
		t.Errorf("NewFetchPurseTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeAccountFetch)
	}
}

func TestHandleAccountHashTask(t *testing.T) {
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
	task, err := NewAccountHashTask("fa12d2dd5547714f8c2754d418aa8c9d59dc88780350cb4254d622e2d4ef7e69")
	if err != nil {
		t.Errorf("Unable to create a NewAccountHashTask : %s", err)
	}
	err = HandleAccountHashTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleAccountHashTask : %s", err)
	}
}

func TestHandleAccountTask(t *testing.T) {
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
	task, err := NewAccountTask("0106ca7c39cd272dbf21a86eeb3b36b7c26e2e9b94af64292419f7862936bca2ca")
	if err != nil {
		t.Errorf("Unable to create a NewAccountTask : %s", err)
	}
	err = HandleAccountTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleAccountTask : %s", err)
	}
}

func TestHandlePurseTask(t *testing.T) {
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
	task, err := NewPurseTask("uref-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1-007")
	if err != nil {
		t.Errorf("Unable to create a NewPurseTask : %s", err)
	}
	err = HandlePurseTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandlePurseTask : %s", err)
	}
}

func TestHandleFetchPurseTask(t *testing.T) {
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
	task, err := NewFetchPurseTask("uref-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1-007")
	if err != nil {
		t.Errorf("Unable to create a NewFetchPurseTask : %s", err)
	}
	err = HandleFetchPurseTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleFetchPurseTask : %s", err)
	}
}
