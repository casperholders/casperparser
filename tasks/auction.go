// Package tasks Define the auction task payload and handler
package tasks

import (
	"casperParser/db"
	"context"
	"github.com/hibiken/asynq"
)

// TypeAuction Task auction type
const (
	TypeAuction = "auction:raw"
)

// NewAuctionTask Used for auction
func NewAuctionTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeAuction, nil), nil
}

// HandleAuctionTask fetch auction from the rpc endpoint, parse it, and insert it in the database
func HandleAuctionTask(ctx context.Context, t *asynq.Task) error {
	auctionParsed, err := WorkerRpcClient.GetAuction()
	if err != nil {
		return err
	}
	var rowsToInsertBids [][]interface{}
	var rowsToInsertDelegators [][]interface{}

	for _, b := range auctionParsed.AuctionState.Bids {
		bidRow := []interface{}{b.PublicKey, b.Bid.BondingPurse, b.Bid.StakedAmount, b.Bid.DelegationRate, b.Bid.Inactive}
		rowsToInsertBids = append(rowsToInsertBids, bidRow)
		for _, d := range b.Bid.Delegators {
			delegatorRow := []interface{}{d.PublicKey, d.Delegatee, d.StakedAmount, d.BondingPurse}
			rowsToInsertDelegators = append(rowsToInsertDelegators, delegatorRow)
		}
	}

	var database = db.DB{Postgres: WorkerPool}

	err = database.InsertAuction(ctx, rowsToInsertBids, rowsToInsertDelegators)
	if err != nil {
		return err
	}

	return nil
}
