// Package block provide a struct for unmarshalling a json Block response from Casper RPC
package block

import "math/big"

type Result struct {
	Block struct {
		Hash   string `json:"hash"`
		Header struct {
			ParentHash      string `json:"parent_hash"`
			StateRootHash   string `json:"state_root_hash"`
			BodyHash        string `json:"body_hash"`
			RandomBit       bool   `json:"random_bit"`
			AccumulatedSeed string `json:"accumulated_seed"`
			Timestamp       string `json:"timestamp"`
			EraID           int    `json:"era_id"`
			Height          int    `json:"height"`
			ProtocolVersion string `json:"protocol_version"`
			EraEnd          *struct {
				EraReport struct {
					Equivocators []string `json:"equivocators"`
					Rewards      []struct {
						Validator string  `json:"validator"`
						Amount    big.Int `json:"amount"`
					} `json:"rewards"`
					InactiveValidators []string `json:"inactiveValidators"`
				} `json:"era_report"`
				NextEraValidatorWeights []struct {
					Validator string `json:"validator"`
					Weight    string `json:"weight"`
				} `json:"next_era_validator_weights"`
			} `json:"era_end"`
		} `json:"header"`
		Body struct {
			Proposer       string   `json:"proposer"`
			DeployHashes   []string `json:"deploy_hashes"`
			TransferHashes []string `json:"transfer_hashes"`
		} `json:"body"`
		Proofs []struct {
			PublicKey string `json:"public_key"`
			Signature string `json:"signature"`
		} `json:"proofs"`
	} `json:"block"`
}
