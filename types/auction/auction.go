package auction

type Result struct {
	AuctionState struct {
		StateRootHash string `json:"state_root_hash"`
		BlockHeight   int    `json:"block_height"`
		EraValidators []struct {
			EraID            int `json:"era_id"`
			ValidatorWeights []struct {
				PublicKey string `json:"public_key"`
				Weight    string `json:"weight"`
			} `json:"validator_weights"`
		} `json:"era_validators"`
		Bids []struct {
			PublicKey string `json:"public_key"`
			Bid       struct {
				BondingPurse   string `json:"bonding_purse"`
				StakedAmount   string `json:"staked_amount"`
				DelegationRate int    `json:"delegation_rate"`
				Delegators     []struct {
					PublicKey    string `json:"public_key"`
					StakedAmount string `json:"staked_amount"`
					BondingPurse string `json:"bonding_purse"`
					Delegatee    string `json:"delegatee"`
				} `json:"delegators"`
				Inactive bool `json:"inactive"`
			} `json:"bid"`
		} `json:"bids"`
	} `json:"auction_state"`
}
