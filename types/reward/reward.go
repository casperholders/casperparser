// Package reward provide a struct for unmarshalling a json reward response from Casper RPC
package reward

type Result struct {
	EraSummary *struct {
		BlockHash   string `json:"block_hash"`
		EraId       int    `json:"era_id"`
		StoredValue struct {
			EraInfo struct {
				SeigniorageAllocations []struct {
					Delegator *struct {
						DelegatorPublicKey string `json:"delegator_public_key"`
						ValidatorPublicKey string `json:"validator_public_key"`
						Amount             string `json:"amount"`
					} `json:"Delegator"`
					Validator *struct {
						ValidatorPublicKey string `json:"validator_public_key"`
						Amount             string `json:"amount"`
					} `json:"Validator"`
				} `json:"seigniorage_allocations"`
			} `json:"EraInfo"`
		} `json:"stored_value"`
	} `json:"era_summary"`
}
