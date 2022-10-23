package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewContractPackageRawTask(t *testing.T) {
	task, err := NewContractPackageRawTask("3cb7d7849ebbd75b08d1883cc2642f846317fc5d86d5327c1102aff4ed9e1482", "03eb82b2e02c5880cd03fcc75580505571c69d476ce28d6cdbb0ee1930cf5950", "017fbbccf39a639a1a5f469e3fb210d9f355b532bd786f945409f0fc9a8c6313b1")
	if err != nil {
		t.Errorf("Unable to create a NewContractPackageRawTask : %s", err)
	}
	if task.Type() != TypeContractPackageRaw {
		t.Errorf("NewContractPackageRawTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeContractPackageRaw)
	}
}

func TestHandleContractPackageRawTask(t *testing.T) {
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
	task, err := NewBlockRawTask(981072)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
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
	task, err = NewContractPackageRawTask("3cb7d7849ebbd75b08d1883cc2642f846317fc5d86d5327c1102aff4ed9e1482", "03eb82b2e02c5880cd03fcc75580505571c69d476ce28d6cdbb0ee1930cf5950", "01ff85d8d335d2e5e1a8ba3554b447e2a61853971fc2a5bf9f1302557ef5eb2d4f")
	if err != nil {
		t.Errorf("Unable to create a NewContractPackageRawTask : %s", err)
	}
	err = HandleContractPackageRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleContractPackageRawTask : %s", err)
	}
}
