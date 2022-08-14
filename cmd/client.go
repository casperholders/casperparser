// Package cmd define the client command
package cmd

import (
	"casperParser/db"
	"casperParser/tasks"
	"context"
	"github.com/Jeffail/gabs/v2"
	"github.com/hibiken/asynq"
	"github.com/r3labs/sse/v2"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var database db.DB
var client *asynq.Client
var pool int
var event string
var wg sync.WaitGroup

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a client",
	Long:  `The client will fetch the last block in the database and add all missing blocks to the queue to be parsed.`,
	Run: func(cmd *cobra.Command, args []string) {
		client = asynq.NewClient(getRedisConf(cmd))
		defer client.Close()
		pgPool, err := db.NewPGXPool(context.Background(), getDatabaseConnectionString(), pool)
		defer pgPool.Close()
		if err != nil {
			log.Fatal(err)
		}
		database = db.DB{Postgres: pgPool}
		go listenEvents()
		wg.Add(1)
		go startClient()
		wg.Add(1)
		wg.Wait()
	},
}

// init the command flags
func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().IntVarP(&pool, "pool", "p", 10, "Database connection pool max connections")
	clientCmd.Flags().StringVarP(&event, "event", "e", "http://127.0.0.1:9999/events/main", "Node main event endpoint")
}

// startClient and add all blocks to the queue
func startClient() {
	defer wg.Done()
	blocks, err := database.GetMissingBlocks(context.Background())
	if err != nil {
		log.Println("Unable to verify if there's missing blocks in the db.")
		log.Println(err)
	}
	for _, block := range blocks {
		addBlockTask(block)
	}
	lastBlock := getLastBlockInDatabase()
	rpcClient := getRpcClient()
	lastBlockHeight, err := rpcClient.GetLastBlockHeight()
	if err != nil {
		log.Println("Unable to determine the last block height on the blockchain.")
		return
	}
	for i := lastBlock; i <= lastBlockHeight; i++ {
		addBlockTask(i)
	}
}

// addBlockTask to the queue
func addBlockTask(height int) {
	task, err := tasks.NewBlockRawTask(height)
	if err != nil {
		log.Printf("could not create task: %v\n", err)
	}
	_, err = client.Enqueue(task, asynq.Queue("blocks"))
	if err != nil {
		log.Printf("could not enqueue task: %v\n", err)
	}
}

// getLastBlockInDatabase defined by the max height block in the db
func getLastBlockInDatabase() int {
	lastBlock := 0
	err := database.Postgres.QueryRow(context.Background(), `select max(height) from blocks`).Scan(&lastBlock)
	if err != nil {
		log.Println("Failed to find the last block in the DB. Starting from 0.")
		log.Println(err)
	}
	return lastBlock
}

func listenEvents() {
	defer wg.Done()
	clientSSE := sse.NewClient(event)

	err := clientSSE.Subscribe("", func(msg *sse.Event) {
		var transforms *gabs.Container
		transforms, _ = gabs.ParseJSON(msg.Data)
		height, ok := transforms.S("BlockAdded", "block", "header", "height").Data().(float64)
		if ok {
			addBlockTask(int(height))
		}
	})
	if err != nil {
		log.Println(err)
	}
}
