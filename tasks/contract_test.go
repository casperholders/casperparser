package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewContractRawTask(t *testing.T) {
	task, err := NewContractRawTask("a37f861cad9bb577d6062512b85695083579056bfcb3c5db56650cc7687e7f17", "03eb82b2e02c5880cd03fcc75580505571c69d476ce28d6cdbb0ee1930cf5950", "017fbbccf39a639a1a5f469e3fb210d9f355b532bd786f945409f0fc9a8c6313b1")
	if err != nil {
		t.Errorf("Unable to create a NewContractRawTask : %s", err)
	}
	if task.Type() != TypeContractRaw {
		t.Errorf("NewContractRawTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeContractRaw)
	}
}

func TestHandleContractRawTask(t *testing.T) {
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
	task, err := NewBlockRawTask(1894266)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewDeployRawTask("38f988ec7dc18d0af88066a72f037ad242155364bf01ce036e028e1f0e7dbba0")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleDeployRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewContractRawTask("29c8584098d57238f78475df60facdd00514ec2aec9da5afdb75ad4fc070e64e", "38f988ec7dc18d0af88066a72f037ad242155364bf01ce036e028e1f0e7dbba0", "015b2d9fa4eb0a7832d4082627f0eee02eefde45e5549c08406a17980ce2455ab7")
	if err != nil {
		t.Errorf("Unable to create a NewContractRawTask : %s", err)
	}
	err = HandleContractRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleContractRawTask : %s", err)
	}
}
