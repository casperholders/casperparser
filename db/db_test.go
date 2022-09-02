package db

import (
	"context"
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	dbconstring := os.Getenv("CASPER_PARSER_DATABASE")
	pool, err := NewPGXPool(context.Background(), dbconstring, 10)
	if err != nil {
		t.Errorf("Unable to init the database pool : %s", err)
	}
	db := DB{Postgres: pool}
	t.Run("Should insert raw block", func(t *testing.T) {
		err = db.InsertRawBlock(context.Background(), "224db16e646f1455daa0a753031516d98533be4b3a5180c949743a7e689aa22c", `{"block": {"body": {"proposer": "0106ca7c39cd272dbf21a86eeb3b36b7c26e2e9b94af64292419f7862936bca2ca", "deploy_hashes": ["00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "45b0bdcbf3cf0b8d86b17996ff1ebb68d025b0ac355dd7a4736a91a11e70c935", "5d09796fea53aa1612e8073351d9788b5631a92dd198326b4d8e20b9b8918767"], "transfer_hashes": []}, "hash": "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731", "header": {"era_id": 0, "height": 64, "era_end": null, "body_hash": "1322de564ba4bbeccab3cb9c7fc52f0eb2ade809408c85eea582b0cc7dc11d85", "timestamp": "2021-04-08T18:10:51.008Z", "random_bit": true, "parent_hash": "f14e4b402ca3ad80ac445cfa37b3a78d344f47319bc036012982f8c3669c2c57", "state_root_hash": "af60ed13e7998ec3b291096177b6a840311a0232655cadb70a4ec577ceef5c43", "accumulated_seed": "20177f34b8fb5d439cce93c5c249aef0495370f6291fc937fd9e2dc452386ccf", "protocol_version": "1.0.0"}, "proofs": []}, "api_version": "1.4.6"}`)
		if err != nil {
			t.Errorf("Unable to insert raw block : %s", err)
		}
	})
	t.Run("Should get raw block", func(t *testing.T) {
		_, err = db.GetRawBlock(context.Background(), "224db16e646f1455daa0a753031516d98533be4b3a5180c949743a7e689aa22c")
		if err != nil {
			t.Errorf("Unable to insert raw block : %s", err)
		}
	})
	t.Run("Should insert block", func(t *testing.T) {
		err = db.InsertBlock(context.Background(), "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731", 0, "2021-04-08 18:10:51.008000 +00:00", 64, false, `{"block": {"body": {"proposer": "0106ca7c39cd272dbf21a86eeb3b36b7c26e2e9b94af64292419f7862936bca2ca", "deploy_hashes": ["00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "45b0bdcbf3cf0b8d86b17996ff1ebb68d025b0ac355dd7a4736a91a11e70c935", "5d09796fea53aa1612e8073351d9788b5631a92dd198326b4d8e20b9b8918767"], "transfer_hashes": []}, "hash": "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731", "header": {"era_id": 0, "height": 64, "era_end": null, "body_hash": "1322de564ba4bbeccab3cb9c7fc52f0eb2ade809408c85eea582b0cc7dc11d85", "timestamp": "2021-04-08T18:10:51.008Z", "random_bit": true, "parent_hash": "f14e4b402ca3ad80ac445cfa37b3a78d344f47319bc036012982f8c3669c2c57", "state_root_hash": "af60ed13e7998ec3b291096177b6a840311a0232655cadb70a4ec577ceef5c43", "accumulated_seed": "20177f34b8fb5d439cce93c5c249aef0495370f6291fc937fd9e2dc452386ccf", "protocol_version": "1.0.0"}, "proofs": []}, "api_version": "1.4.6"}`)
		if err != nil {
			t.Errorf("Unable to insert block : %s", err)
		}
	})
	t.Run("Should validate block", func(t *testing.T) {
		err = db.ValidateBlock(context.Background(), "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731")
		if err != nil {
			t.Errorf("Unable to insert block : %s", err)
		}
	})
	t.Run("Should InsertDeploy", func(t *testing.T) {
		err = db.InsertDeploy(context.Background(), "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "232824230", false, "2021-04-08 18:10:32.115000 +00:00", "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731", "moduleBytes", `{"deploy": {"hash": "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "header": {"ttl": "1h", "account": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "body_hash": "a14865e1d6016f4c83a240a7df64a538798b2bb392d86674a4e33a59901a354b", "gas_price": 1, "timestamp": "2021-04-08T18:10:32.115Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "05003ad0b814", "parsed": "89000000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"ModuleBytes": {"args": [["public_key", {"bytes": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "parsed": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "cl_type": "PublicKey"}], ["amount", {"bytes": "0500282e8cd1", "parsed": "900000000000", "cl_type": "U512"}], ["delegation_rate", {"bytes": "0a", "parsed": 10, "cl_type": "U8"}]], "module_bytes": ""}}, "approvals": [{"signer": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "signature": "01d911dea193709debe86db026bb2b9a43e9d0985a3a6f4cde749b55f08b83c12ac9ec7a47b8ae6e7a12cf5dbd96a0314846febf1c9fa95f5aea5a480a170bcc0e"}]}, "api_version": "1.4.7", "execution_results": [{"result": {"Failure": {"cost": "232824230", "effect": {"operations": [{"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "kind": "Read"}, {"key": "balance-2c4bac63bc01ddc6f76e2bc2bcc6af61d6efa9cd65b22785135adea57f98b24c", "kind": "Write"}, {"key": "balance-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1", "kind": "Write"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "kind": "Read"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "kind": "Read"}], "transforms": [{"key": "balance-2c4bac63bc01ddc6f76e2bc2bcc6af61d6efa9cd65b22785135adea57f98b24c", "transform": {"WriteCLValue": {"bytes": "04808f9b95", "parsed": "2510000000", "cl_type": "U512"}}}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "balance-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1", "transform": {"AddUInt512": "89000000000"}}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}]}, "transfers": [], "error_message": "ApiError::AuctionError(4) [64516]"}}, "block_hash": "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731"}]}`, "moduleBytes", "", "", "", "", "")
		if err != nil {
			t.Errorf("Unable to InsertDeploy : %s", err)
		}
	})
	t.Run("Should CountDeploys", func(t *testing.T) {
		count, err := db.CountDeploys(context.Background(), []string{"00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2"})
		if err != nil {
			t.Errorf("Unable to InsertDeploy : %s", err)
		}
		if count != 1 {
			t.Errorf("Count should be 1, received : %d", count)
		}
	})
	t.Run("Should UpdateDeploy", func(t *testing.T) {
		err = db.UpdateDeploy(context.Background(), "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "232824230", false, "2021-04-08 18:10:32.115000 +00:00", "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731", "moduleBytes", "moduleBytes", "", "", "", "", "")
		if err != nil {
			t.Errorf("Unable to UpdateDeploy : %s", err)
		}
	})
	t.Run("Should InsertRawDeploy", func(t *testing.T) {
		err = db.InsertRawDeploy(context.Background(), "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", `{"deploy": {"hash": "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "header": {"ttl": "1h", "account": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "body_hash": "a14865e1d6016f4c83a240a7df64a538798b2bb392d86674a4e33a59901a354b", "gas_price": 1, "timestamp": "2021-04-08T18:10:32.115Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "05003ad0b814", "parsed": "89000000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"ModuleBytes": {"args": [["public_key", {"bytes": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "parsed": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "cl_type": "PublicKey"}], ["amount", {"bytes": "0500282e8cd1", "parsed": "900000000000", "cl_type": "U512"}], ["delegation_rate", {"bytes": "0a", "parsed": 10, "cl_type": "U8"}]], "module_bytes": ""}}, "approvals": [{"signer": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "signature": "01d911dea193709debe86db026bb2b9a43e9d0985a3a6f4cde749b55f08b83c12ac9ec7a47b8ae6e7a12cf5dbd96a0314846febf1c9fa95f5aea5a480a170bcc0e"}]}, "api_version": "1.4.7", "execution_results": [{"result": {"Failure": {"cost": "232824230", "effect": {"operations": [{"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "kind": "Read"}, {"key": "balance-2c4bac63bc01ddc6f76e2bc2bcc6af61d6efa9cd65b22785135adea57f98b24c", "kind": "Write"}, {"key": "balance-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1", "kind": "Write"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "kind": "Read"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "kind": "Read"}], "transforms": [{"key": "balance-2c4bac63bc01ddc6f76e2bc2bcc6af61d6efa9cd65b22785135adea57f98b24c", "transform": {"WriteCLValue": {"bytes": "04808f9b95", "parsed": "2510000000", "cl_type": "U512"}}}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "balance-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1", "transform": {"AddUInt512": "89000000000"}}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}]}, "transfers": [], "error_message": "ApiError::AuctionError(4) [64516]"}}, "block_hash": "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731"}]}`)
		if err != nil {
			t.Errorf("Unable to InsertRawDeploy : %s", err)
		}
	})
	t.Run("Should GetMissingBlocks", func(t *testing.T) {
		_, err = db.GetMissingBlocks(context.Background())
		if err != nil {
			t.Errorf("Unable to GetMissingBlocks : %s", err)
		}
	})
	t.Run("Should GetMissingMetadataDeploysHash", func(t *testing.T) {
		_, err = db.GetMissingMetadataDeploysHash(context.Background())
		if err != nil {
			t.Errorf("Unable to GetMissingMetadataDeploysHash : %s", err)
		}
	})
	t.Run("Should GetDeploy", func(t *testing.T) {
		_, err = db.GetDeploy(context.Background(), "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2")
		if err != nil {
			t.Errorf("Unable to GetDeploy : %s", err)
		}
	})
}
