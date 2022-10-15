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
	task, err := NewBlockRawTask(84)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockVerifyTask("d9dd87b06db708800036da57f1acf9302f51dde2a57b548ad4804ceb2377bdff")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockVerifyTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockRawTask(1153698)
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockRawTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
	task, err = NewBlockVerifyTask("fc204a0bc7788604fd0ded0ac19a73b687d12a8d735ccf57f3c65ce58d6f4d1f")
	if err != nil {
		t.Errorf("Unable to create a NewBlockRawTask : %s", err)
	}
	err = HandleBlockVerifyTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleBlockRawTask : %s", err)
	}
}
