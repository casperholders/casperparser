package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewAuctionTask(t *testing.T) {
	task, err := NewAuctionTask()
	if err != nil {
		t.Errorf("Unable to create a NewAuctionTask : %s", err)
	}
	if task.Type() != "auction:raw" {
		t.Errorf("NewAuctionTask has a bad name. Received : %s. Expected : %s", task.Type(), "auction:raw")
	}
}

func TestHandleAuctionTask(t *testing.T) {
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
	task, err := NewAuctionTask()
	if err != nil {
		t.Errorf("Unable to create a NewAuctionTask : %s", err)
	}
	err = HandleAuctionTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleAuctionTask : %s", err)
	}
}
