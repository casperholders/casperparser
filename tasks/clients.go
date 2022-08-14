// Package tasks Contains different worker pool/clients needed for the workers
package tasks

import (
	"casperParser/rpc"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4/pgxpool"
)

var WorkerPool *pgxpool.Pool
var WorkerAsyncClient *asynq.Client
var WorkerRpcClient *rpc.Client
