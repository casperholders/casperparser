package tasks

import (
	"casperParser/db"
	"casperParser/rpc"
	"context"
	"github.com/hibiken/asynq"
	"os"
	"testing"
)

func TestNewRewardTask(t *testing.T) {
	task, err := NewRewardTask("fc204a0bc7788604fd0ded0ac19a73b687d12a8d735ccf57f3c65ce58d6f4d1f")
	if err != nil {
		t.Errorf("Unable to create a NewRewardTask : %s", err)
	}
	if task.Type() != TypeReward {
		t.Errorf("NewRewardTask has a bad name. Received : %s. Expected : %s", task.Type(), TypeReward)
	}
}

func TestHandleRewardTask(t *testing.T) {
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
	task, err := NewRewardTask("8d6a98a977482af4eb308cfb4ebdf6981643afdc06f56d6589792808992f56fe")
	if err != nil {
		t.Errorf("Unable to create a NewRewardTask : %s", err)
	}
	err = HandleRewardTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unable to run HandleRewardTask : %s", err)
	}
}
