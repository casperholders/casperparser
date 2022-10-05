// Package tasks Define the auction task payload and handler
package tasks

import (
	"casperParser/db"
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"math/big"
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
		bStaked, ok := new(big.Int).SetString(b.Bid.StakedAmount, 10)
		if ok {
			bidRow := []interface{}{b.PublicKey, b.Bid.BondingPurse, bStaked.Int64(), b.Bid.DelegationRate, b.Bid.Inactive}
			rowsToInsertBids = append(rowsToInsertBids, bidRow)
			for _, d := range b.Bid.Delegators {
				dStaked, dok := new(big.Int).SetString(d.StakedAmount, 10)
				if dok {
					delegatorRow := []interface{}{d.PublicKey, d.Delegatee, dStaked.Int64(), d.BondingPurse}
					rowsToInsertDelegators = append(rowsToInsertDelegators, delegatorRow)
				} else {
					return fmt.Errorf("cannot convert stake to bigint for %s / v: %s", d.PublicKey, d.Delegatee)
				}
			}
		} else {
			return fmt.Errorf("cannot convert stake to bigint for v: %s", b.PublicKey)
		}
	}

	var database = db.DB{Postgres: WorkerPool}

	err = database.InsertAuction(ctx, rowsToInsertBids, rowsToInsertDelegators)
	if err != nil {
		return err
	}

	return nil
}
