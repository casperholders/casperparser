package deploy

import (
	"casperParser/types/config"
	"casperParser/utils"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"testing"
)

var transferDeploy = `{"deploy": {"hash": "2d0e59821d67125ab7a07ac719ed6696ce4dd4498ef6f3c283ac7d02f9de7259", "header": {"ttl": "1h", "account": "017717a9bb1f07cbb1b6c3afaaad9ff3b8a5b75ea13e5aae6ce33b4b74676c647c", "body_hash": "36088013e89250b53807fd562a047f40b076f97db0c557eec80d3b2a006c7400", "gas_price": 1, "timestamp": "2021-04-08T18:28:13.061Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "0400ca9a3b", "parsed": "1000000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"Transfer": {"args": [["amount", {"bytes": "05007c6f5de8", "parsed": "998000000000", "cl_type": "U512"}], ["target", {"bytes": "e70b850efb68c64e2443da2386452b0d8e4e799362edef0ff56eea8efb114815", "parsed": "e70b850efb68c64e2443da2386452b0d8e4e799362edef0ff56eea8efb114815", "cl_type": {"ByteArray": 32}}], ["id", {"bytes": "00", "parsed": null, "cl_type": {"Option": "U64"}}]]}}, "approvals": [{"signer": "017717a9bb1f07cbb1b6c3afaaad9ff3b8a5b75ea13e5aae6ce33b4b74676c647c", "signature": "0115b52bfc9632454713987c42a75f85d615b3e035a82894ed112f75779ed44cba0694941c1e1cc68cac0f902d6dbd0524f2f81667969c333dfe9a2e0d143edd06"}]}, "api_version": "1.4.6", "execution_results": [{"result": {"Success": {"cost": "10000", "effect": {"operations": [{"key": "balance-f3abd4d174755d6127e2145876b54e14a1280694cfd1f3924565397808cb7c3f", "kind": "Write"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "kind": "Read"}, {"key": "balance-62f7fe1cecb1a4c600ffa791479ce52fb8cbda408815f4dd1b1e0d82e704579a", "kind": "Write"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "kind": "Read"}, {"key": "balance-0a24ef56971d46bfefbd5590afe20e5f3482299aba74e1a0fc33a55008cf9453", "kind": "Write"}, {"key": "transfer-f8d688cfc73ca7a06925b0b1226d1b902ea09cb799aecee8c3dcbb69c23acd5a", "kind": "Write"}, {"key": "deploy-2d0e59821d67125ab7a07ac719ed6696ce4dd4498ef6f3c283ac7d02f9de7259", "kind": "Write"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "kind": "Read"}, {"key": "account-hash-e70b850efb68c64e2443da2386452b0d8e4e799362edef0ff56eea8efb114815", "kind": "Read"}], "transforms": [{"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "balance-0a24ef56971d46bfefbd5590afe20e5f3482299aba74e1a0fc33a55008cf9453", "transform": {"AddUInt512": "998000000000"}}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "deploy-2d0e59821d67125ab7a07ac719ed6696ce4dd4498ef6f3c283ac7d02f9de7259", "transform": {"WriteDeployInfo": {"gas": "10000", "from": "account-hash-7a697fb8fcbb706d0c42fe97c2b0b2472c634569bf71e71d07ad1b158c01f839", "source": "uref-f3abd4d174755d6127e2145876b54e14a1280694cfd1f3924565397808cb7c3f-007", "transfers": ["transfer-f8d688cfc73ca7a06925b0b1226d1b902ea09cb799aecee8c3dcbb69c23acd5a"], "deploy_hash": "2d0e59821d67125ab7a07ac719ed6696ce4dd4498ef6f3c283ac7d02f9de7259"}}}, {"key": "account-hash-e70b850efb68c64e2443da2386452b0d8e4e799362edef0ff56eea8efb114815", "transform": "Identity"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "balance-f3abd4d174755d6127e2145876b54e14a1280694cfd1f3924565397808cb7c3f", "transform": {"WriteCLValue": {"bytes": "04f06c3577", "parsed": "1999990000", "cl_type": "U512"}}}, {"key": "transfer-f8d688cfc73ca7a06925b0b1226d1b902ea09cb799aecee8c3dcbb69c23acd5a", "transform": {"WriteTransfer": {"id": null, "to": "account-hash-e70b850efb68c64e2443da2386452b0d8e4e799362edef0ff56eea8efb114815", "gas": "0", "from": "account-hash-7a697fb8fcbb706d0c42fe97c2b0b2472c634569bf71e71d07ad1b158c01f839", "amount": "998000000000", "source": "uref-f3abd4d174755d6127e2145876b54e14a1280694cfd1f3924565397808cb7c3f-007", "target": "uref-0a24ef56971d46bfefbd5590afe20e5f3482299aba74e1a0fc33a55008cf9453-004", "deploy_hash": "2d0e59821d67125ab7a07ac719ed6696ce4dd4498ef6f3c283ac7d02f9de7259"}}}, {"key": "balance-62f7fe1cecb1a4c600ffa791479ce52fb8cbda408815f4dd1b1e0d82e704579a", "transform": {"AddUInt512": "10000"}}]}, "transfers": ["transfer-f8d688cfc73ca7a06925b0b1226d1b902ea09cb799aecee8c3dcbb69c23acd5a"]}}, "block_hash": "da006f14c000b23977c5b40c7163a5d5a993c08b156751c30b096522362e1a16"}]}`
var storedContractByHashDeploy = `{"deploy": {"hash": "4cb068f987a91f5b60dc496dff632597922b01022c0e313048a1b72685a156ce", "header": {"ttl": "30m", "account": "01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681", "body_hash": "82af02c1176bc14839b9ce8e4719f9862fe74a80de865cc607fe1f9ffe266d52", "gas_price": 1, "timestamp": "2021-09-01T00:36:51.846Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "0410200395", "parsed": "2500010000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"StoredContractByHash": {"args": [["delegator", {"bytes": "01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681", "parsed": "01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681", "cl_type": "PublicKey"}], ["validator", {"bytes": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5", "parsed": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5", "cl_type": "PublicKey"}], ["amount", {"bytes": "0400ca9a3b", "parsed": "1000000000", "cl_type": "U512"}]], "hash": "93d923e336b20a4c4ca14d592b60e5bd3fe330775618290104f9beb326db7ae2", "entry_point": "delegate"}}, "approvals": [{"signer": "01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681", "signature": "01e6e0aad356b4b0c7011c97237ef2891b601d432d3f6f05f9784c479bbdd52dca984bca0d850f1d047158e9dcbf862d3ed15923b33b8dad14da24dc5093c69a0d"}]}, "api_version": "1.4.6", "execution_results": [{"result": {"Success": {"cost": "2500000000", "effect": {"operations": [{"key": "balance-b8c08e6c203ec937f34897bb3c1bfcf046437bfc70f0c2f76f90e3544c30c7d3", "kind": "Write"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "kind": "Read"}, {"key": "deploy-4cb068f987a91f5b60dc496dff632597922b01022c0e313048a1b72685a156ce", "kind": "Write"}, {"key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e", "kind": "Read"}, {"key": "transfer-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8", "kind": "Write"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "kind": "Read"}, {"key": "bid-72a08a1d2ea95c1b975a6e211e9468075d16ced5d4160250c3d0b1bf0c132314", "kind": "Write"}, {"key": "balance-f6bb43cd45bfc9e9990673b6580c203bdb908ee30f60a54a0e866d0e7abd18c4", "kind": "Write"}, {"key": "balance-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8", "kind": "Write"}, {"key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "kind": "Read"}, {"key": "uref-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8-000", "kind": "Write"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "kind": "Read"}], "transforms": [{"key": "transfer-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8", "transform": {"WriteTransfer": {"id": null, "to": "account-hash-6174cf2e6f8fed1715c9a3bace9c50bfe572eecb763b0ed3f644532616452008", "gas": "0", "from": "account-hash-21834cffb261e94569db30595210525b9a52fa7d0ce97c650fd985b65835263f", "amount": "1000000000", "source": "uref-b8c08e6c203ec937f34897bb3c1bfcf046437bfc70f0c2f76f90e3544c30c7d3-007", "target": "uref-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8-007", "deploy_hash": "4cb068f987a91f5b60dc496dff632597922b01022c0e313048a1b72685a156ce"}}}, {"key": "uref-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8-000", "transform": {"WriteCLValue": {"bytes": "", "parsed": null, "cl_type": "Unit"}}}, {"key": "deploy-4cb068f987a91f5b60dc496dff632597922b01022c0e313048a1b72685a156ce", "transform": {"WriteDeployInfo": {"gas": "2500000000", "from": "account-hash-21834cffb261e94569db30595210525b9a52fa7d0ce97c650fd985b65835263f", "source": "uref-b8c08e6c203ec937f34897bb3c1bfcf046437bfc70f0c2f76f90e3544c30c7d3-007", "transfers": ["transfer-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8"], "deploy_hash": "4cb068f987a91f5b60dc496dff632597922b01022c0e313048a1b72685a156ce"}}}, {"key": "balance-f6bb43cd45bfc9e9990673b6580c203bdb908ee30f60a54a0e866d0e7abd18c4", "transform": {"AddUInt512": "2500010000"}}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "balance-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8", "transform": {"WriteCLValue": {"bytes": "0400ca9a3b", "parsed": "1000000000", "cl_type": "U512"}}}, {"key": "balance-b8c08e6c203ec937f34897bb3c1bfcf046437bfc70f0c2f76f90e3544c30c7d3", "transform": {"WriteCLValue": {"bytes": "06b0a15a291001", "parsed": "1168924910000", "cl_type": "U512"}}}, {"key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e", "transform": "Identity"}, {"key": "bid-72a08a1d2ea95c1b975a6e211e9468075d16ced5d4160250c3d0b1bf0c132314", "transform": {"WriteBid": {"inactive": false, "delegators": {"0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5": {"bonding_purse": "uref-8f300b2f28c264762d1b457f0b54b78e1587c498d32a0a0f000a3ad865832492-007", "staked_amount": "2000000000", "vesting_schedule": null, "delegator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5", "validator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}, "01871be6fab07ead1c7114731143239127527c64d670d1cac86e4a582768f6549d": {"bonding_purse": "uref-811f3014c25ca30d157995b397bc96d410bee64f231d3b0f61b592e2f17972f0-007", "staked_amount": "1000000000", "vesting_schedule": null, "delegator_public_key": "01871be6fab07ead1c7114731143239127527c64d670d1cac86e4a582768f6549d", "validator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}, "018c092e4fb70c03e4de3749beae4cbb3bc9b03d2b854c96e2348d32387351081e": {"bonding_purse": "uref-b1dbc2d5588959510f1a8a14c78c2f223abbe3ae8a4478058fc6833ef923cb1a-007", "staked_amount": "7000000000", "vesting_schedule": null, "delegator_public_key": "018c092e4fb70c03e4de3749beae4cbb3bc9b03d2b854c96e2348d32387351081e", "validator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}, "01c534307e2c7a4839e01ebefae81517fb26d928e3a86802c48b9d47454625bf14": {"bonding_purse": "uref-74447455f516ecad74e3c37040a9833f965bd14300ea8705a51e320f7b2c5687-007", "staked_amount": "102000000000", "vesting_schedule": null, "delegator_public_key": "01c534307e2c7a4839e01ebefae81517fb26d928e3a86802c48b9d47454625bf14", "validator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}, "01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681": {"bonding_purse": "uref-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8-007", "staked_amount": "1000000000", "vesting_schedule": null, "delegator_public_key": "01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681", "validator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}}, "bonding_purse": "uref-a3847d12c1f40fec49c78d4f068ed33bdcdbee5faf69c53e117ec0eb2307d9d2-007", "staked_amount": "98000000000", "delegation_rate": 1, "vesting_schedule": null, "validator_public_key": "0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}}}, {"key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "transform": "Identity"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}]}, "transfers": ["transfer-2977d3aee4b1ae390f5f60014d33f0ec6095cf6343e18cefbedba351fad01eb8"]}}, "block_hash": "5ed8d37ac14bf26c811732fd003424e71835ac78de6049e49fa9c6b3b40bf5f1"}]}`
var storedContractByNameDeploy = `{"deploy": {"hash": "08bb9f11da011d5ec339db7ff3cec2ee6cc6fcbf6aadb76701714a87652f0b6b", "header": {"ttl": "1h", "account": "018afa98ca4be12d613617f7339a2d576950a2f9a92102ca4d6508ee31b54d2c02", "body_hash": "7463a03a4d0b14ca385d41be425c8bed7f8001941a6337d2b276a295ea03aad3", "gas_price": 1, "timestamp": "2021-04-08T18:08:56.394Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "0400c2eb0b", "parsed": "200000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"StoredContractByName": {"args": [["target", {"bytes": "b497711627a79370e1b779dbae5970171c5fcccd3785f1e6593cea0ad6ec7bee", "parsed": "b497711627a79370e1b779dbae5970171c5fcccd3785f1e6593cea0ad6ec7bee", "cl_type": {"ByteArray": 32}}], ["amount", {"bytes": "050010a5d4e8", "parsed": "1000000000000", "cl_type": "U512"}]], "name": "faucet", "entry_point": "call_faucet"}}, "approvals": [{"signer": "018afa98ca4be12d613617f7339a2d576950a2f9a92102ca4d6508ee31b54d2c02", "signature": "0199cd95909e8fdce16520e7cb7e0babac306d37768655c4f25393374c4f7ec527be78230bab3c8b4add5717205ccf669f5308c535d2cab6923383a02c68da8409"}]}, "api_version": "1.4.6", "execution_results": [{"result": {"Failure": {"cost": "11406830", "effect": {"operations": [{"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "kind": "Read"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "kind": "Read"}, {"key": "balance-62f7fe1cecb1a4c600ffa791479ce52fb8cbda408815f4dd1b1e0d82e704579a", "kind": "Write"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "kind": "Read"}, {"key": "balance-b06a1ab0cfb52b5d4f9a08b68a5dbe78e999de0b0484c03e64f5c03897cf637b", "kind": "Write"}], "transforms": [{"key": "balance-b06a1ab0cfb52b5d4f9a08b68a5dbe78e999de0b0484c03e64f5c03897cf637b", "transform": {"WriteCLValue": {"bytes": "0800a4754315b9c68a", "parsed": "9999883523600000000", "cl_type": "U512"}}}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "balance-62f7fe1cecb1a4c600ffa791479ce52fb8cbda408815f4dd1b1e0d82e704579a", "transform": {"AddUInt512": "200000000"}}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}]}, "transfers": [], "error_message": "User error: 1"}}, "block_hash": "f14e4b402ca3ad80ac445cfa37b3a78d344f47319bc036012982f8c3669c2c57"}]}`
var storedVersionedContractByNameDeploy = `{"deploy": {"hash": "c64a086d2ee622d35b4c1d23a35b66500f0e986d65e1fb0c712d9ddd6fc2ba73", "header": {"ttl": "30m", "account": "01ffac94de6c1db1f4f85d9ab2cf5f9eedf2fc6f8e8d27e40e2265bf8e19433193", "body_hash": "331772bd14bcf6bbf498e4bd3fb779fe2f9fbe7ed622dbd7b62276f37551f73d", "gas_price": 1, "timestamp": "2022-07-13T17:01:45.665Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "0500282e8cd1", "parsed": "900000000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"StoredVersionedContractByName": {"args": [["token_contract", {"bytes": "521f4f0319b3f46e325fc09810644e8107c1886540ef996f79e94981f47b31b9", "parsed": "521f4f0319b3f46e325fc09810644e8107c1886540ef996f79e94981f47b31b9", "cl_type": {"ByteArray": 32}}], ["account", {"bytes": "00a3aa343007ff3951f564a754cd3e2f31b2f7332788d02efe874584fbc49ca56d", "parsed": {"Account": "account-hash-a3aa343007ff3951f564a754cd3e2f31b2f7332788d02efe874584fbc49ca56d"}, "cl_type": "Key"}], ["id", {"bytes": "06000000746f6b656e32", "parsed": "token2", "cl_type": "String"}]], "name": "erc1155_test_call", "version": null, "entry_point": "check_balance_of"}}, "approvals": [{"signer": "01ffac94de6c1db1f4f85d9ab2cf5f9eedf2fc6f8e8d27e40e2265bf8e19433193", "signature": "01d52ba6445faed7f36341a574173ddb722d68c2640c84fd58db332c351d3cd8c31b08a6e6ad16be2695361a234fc02efd3bff2c02c6cfccbd3b28e45ff5ef090c"}]}, "api_version": "1.4.6", "execution_results": [{"result": {"Success": {"cost": "32364610", "effect": {"operations": [], "transforms": [{"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "transform": "Identity"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e", "transform": "Identity"}, {"key": "balance-3bb60b6e4a5615b830f64cea5dfbf9a05a403903ef69e7ec709c4da00010fb92", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "balance-3bb60b6e4a5615b830f64cea5dfbf9a05a403903ef69e7ec709c4da00010fb92", "transform": {"WriteCLValue": {"bytes": "076c37e34c7cd809", "parsed": "2771303167899500", "cl_type": "U512"}}}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": {"AddUInt512": "900000000000"}}, {"key": "hash-c29a016dd963e752086ff07c6c89acdcc06f8fc45a45a54f20144d8ffca987db", "transform": "Identity"}, {"key": "hash-45cf6c9653bd088e1b4aac1603026017144fddedec09e8e71bc5a08c4de2340a", "transform": "Identity"}, {"key": "hash-65037f70b60aed5489cdffc2efe71524250354ab171d9b74b311c7080d1c3be7", "transform": "Identity"}, {"key": "hash-521f4f0319b3f46e325fc09810644e8107c1886540ef996f79e94981f47b31b9", "transform": "Identity"}, {"key": "hash-92bf06710a08c6fab9e9a84cac07d10d23cd0455138510ebad9e81f452480464", "transform": "Identity"}, {"key": "hash-60de10fb518044cca0904d4a6ae0932968f09a4c8c492c37c599a5cad4b92859", "transform": "Identity"}, {"key": "dictionary-ee6be9bed3cad4fd1e74327e33cb84a5ab9db56acf75ed1dffce2d06a472fb2b", "transform": "Identity"}, {"key": "uref-d9640273dc5da8f33d8c0523e7b282d939da7caaf07d90093d0b146ca87ce100-000", "transform": {"WriteCLValue": {"bytes": "03400d03", "parsed": "200000", "cl_type": "U256"}}}, {"key": "deploy-c64a086d2ee622d35b4c1d23a35b66500f0e986d65e1fb0c712d9ddd6fc2ba73", "transform": {"WriteDeployInfo": {"gas": "32364610", "from": "account-hash-a3aa343007ff3951f564a754cd3e2f31b2f7332788d02efe874584fbc49ca56d", "source": "uref-3bb60b6e4a5615b830f64cea5dfbf9a05a403903ef69e7ec709c4da00010fb92-007", "transfers": [], "deploy_hash": "c64a086d2ee622d35b4c1d23a35b66500f0e986d65e1fb0c712d9ddd6fc2ba73"}}}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "balance-65791c5e4cc22b504aa6a98846c00bfa452cd89982657d03bf1118da9616819f", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": {"WriteCLValue": {"bytes": "00", "parsed": "0", "cl_type": "U512"}}}, {"key": "balance-65791c5e4cc22b504aa6a98846c00bfa452cd89982657d03bf1118da9616819f", "transform": {"AddUInt512": "900000000000"}}]}, "transfers": []}}, "block_hash": "4cc306cc8b2d0da89307883e678ca7c5642897a1f84814259861438e501a3520"}]}`
var storedVersionedContractByHashDeploy = `{"deploy": {"hash": "140944a190e7c8375b4535d35bf9f3863c32d18f15f17f2cd91261706f1bc952", "header": {"ttl": "30m", "account": "02033ad1cd00f637c3b4593a721194c76c224a84ed7c37f91f36110ccf12a8c24706", "body_hash": "c40a4e7b2a60843dd82a6c5fba5e204c24edf485fb9cdc2f3f50227be53635f9", "gas_price": 1, "timestamp": "2022-08-09T17:16:54.157Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "0400ca9a3b", "parsed": "1000000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"StoredVersionedContractByHash": {"args": [["owner", {"bytes": "00f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76", "parsed": {"Account": "account-hash-f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76"}, "cl_type": "Key"}], ["token_ids", {"bytes": "01000000011a", "parsed": ["26"], "cl_type": {"List": "U256"}}]], "hash": "6ca070c78d4eb468b4db4cbc5cadd815c35e15019a841c137372a88d7e247d1d", "version": null, "entry_point": "burn"}}, "approvals": [{"signer": "02033ad1cd00f637c3b4593a721194c76c224a84ed7c37f91f36110ccf12a8c24706", "signature": "028b2d299bc134b56bdae0e7664ef82822574b48c836d8c33c4a4e4ccd0d82960f51d4894ac51ea33304ef14b36287bc97d951af65225c9fc99815035b5ad5e6ea"}]}, "api_version": "1.4.7", "execution_results": [{"result": {"Success": {"cost": "601040500", "effect": {"operations": [], "transforms": [{"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "transform": "Identity"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e", "transform": "Identity"}, {"key": "balance-91eee3d01fb1276824ebfaba66a8e3aed82333de3f5d52527bd38bc1390b7abf", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "balance-91eee3d01fb1276824ebfaba66a8e3aed82333de3f5d52527bd38bc1390b7abf", "transform": {"WriteCLValue": {"bytes": "0509dcd97545", "parsed": "298329955337", "cl_type": "U512"}}}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": {"AddUInt512": "1000000000"}}, {"key": "hash-56455aea23a154412f438d88b728ace4c6b727be76978c650d655980cd8fad82", "transform": "Identity"}, {"key": "hash-6ca070c78d4eb468b4db4cbc5cadd815c35e15019a841c137372a88d7e247d1d", "transform": "Identity"}, {"key": "hash-a60ea460a5a456baf5d6bae6c73157fa6cd7fe14cb37afaae8c1efe30dffa716", "transform": "Identity"}, {"key": "dictionary-bc841f38e3bfa6bbb4f5759ad19acd37abb8eb2d3b5bab31491837946a102d97", "transform": "Identity"}, {"key": "dictionary-4e907c764f30be91287ed23bd7aef72d83f18a7971c3d0a151637b109ace1bbc", "transform": "Identity"}, {"key": "dictionary-6ec0d66bdea37104b6f74b337ad0da0ddaef03175292145a0da622f7b8557860", "transform": "Identity"}, {"key": "dictionary-85d2318eb5b8ca1700c42890d2f386ae002d220737fed4fe80acaabfc360f376", "transform": {"WriteCLValue": {"bytes": "01000000000d0720000000f48bc7112bb3dceb4f521b8c73a2ee8eef4087d6242f03693ad7e8465c2cb7c74000000036383261633334303337613837663530306537353132376633343638616230333235383063636665366236316137626637313837303530376639353664626466", "parsed": null, "cl_type": "Any"}}}, {"key": "dictionary-4e907c764f30be91287ed23bd7aef72d83f18a7971c3d0a151637b109ace1bbc", "transform": {"WriteCLValue": {"bytes": "0200000001000d07200000007087eeb135c26a288a08778df3c659142a2ff84bb5145140edbd2598eb2fd3d24000000066326433323738613864323438333766323331353662383132643732616237616535656138313436376566653266623731386532393237353663383863643736", "parsed": null, "cl_type": "Any"}}}, {"key": "dictionary-6ec0d66bdea37104b6f74b337ad0da0ddaef03175292145a0da622f7b8557860", "transform": {"WriteCLValue": {"bytes": "01000000000d0720000000f404bd549be7f491c0137351d8c3bf8e64af55eeb71ff65b0a751d2de934b2b64000000065336663323866336133373366633765346233336232616333623034663438643639363065376165393333653335326531613935353235626332326163636463", "parsed": null, "cl_type": "Any"}}}, {"key": "dictionary-bc56d1a2c8b2e4551cea5b56f5b0291cb7801687de291ed46a6bf4b9bc15aa02", "transform": {"WriteCLValue": {"bytes": "01000000000d110a0a20000000b184478ec03598d539c328b3e2931829d20afb3560d98ebdcf8278334ac14e2c020000003236", "parsed": null, "cl_type": "Any"}}}, {"key": "dictionary-bc841f38e3bfa6bbb4f5759ad19acd37abb8eb2d3b5bab31491837946a102d97", "transform": {"WriteCLValue": {"bytes": "01000000000d0b20000000512bd1d406f87de79160fc29f3d73c83650095cf79caa6a3290645ce05355385020000003236", "parsed": null, "cl_type": "Any"}}}, {"key": "dictionary-aa3cc5cf6009313a05120d37688022c92d7d6af7b493524bedaa777068e9f53e", "transform": {"WriteCLValue": {"bytes": "01000000000d0b20000000be0a47d9daff791d4e39932c71b185125f47c47475b6b939f6b40cc4694f8e0b4000000031356339353732616137386336663731356231393963363031356138306633346636643134346363626366366264373330393433663632306138336430393461", "parsed": null, "cl_type": "Any"}}}, {"key": "uref-326441ef0f9c225701cfe73c200332651c6249d17394ab9b67e5b96f52f3c279-000", "transform": "Identity"}, {"key": "uref-326441ef0f9c225701cfe73c200332651c6249d17394ab9b67e5b96f52f3c279-000", "transform": {"WriteCLValue": {"bytes": "0109", "parsed": "9", "cl_type": "U256"}}}, {"key": "uref-2b62aac14fa07969e1597afba8efd0b0e0d168c6012b64705503956008e747b5-000", "transform": {"WriteCLValue": {"bytes": "0400000015000000636f6e74726163745f7061636b6167655f6861736840000000366361303730633738643465623436386234646234636263356361646438313563333565313530313961383431633133373337326138386437653234376431640a0000006576656e745f747970650e00000063657034375f6275726e5f6f6e65050000006f776e65724e0000004b65793a3a4163636f756e7428663264333237386138643234383337663233313536623831326437326162376165356561383134363765666532666237313865323932373536633838636437362908000000746f6b656e5f6964020000003236", "parsed": [{"key": "contract_package_hash", "value": "6ca070c78d4eb468b4db4cbc5cadd815c35e15019a841c137372a88d7e247d1d"}, {"key": "event_type", "value": "cep47_burn_one"}, {"key": "owner", "value": "Key::Account(f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76)"}, {"key": "token_id", "value": "26"}], "cl_type": {"Map": {"key": "String", "value": "String"}}}}}, {"key": "deploy-140944a190e7c8375b4535d35bf9f3863c32d18f15f17f2cd91261706f1bc952", "transform": {"WriteDeployInfo": {"gas": "601040500", "from": "account-hash-f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76", "source": "uref-91eee3d01fb1276824ebfaba66a8e3aed82333de3f5d52527bd38bc1390b7abf-007", "transfers": [], "deploy_hash": "140944a190e7c8375b4535d35bf9f3863c32d18f15f17f2cd91261706f1bc952"}}}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "balance-f6bb43cd45bfc9e9990673b6580c203bdb908ee30f60a54a0e866d0e7abd18c4", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": {"WriteCLValue": {"bytes": "00", "parsed": "0", "cl_type": "U512"}}}, {"key": "balance-f6bb43cd45bfc9e9990673b6580c203bdb908ee30f60a54a0e866d0e7abd18c4", "transform": {"AddUInt512": "1000000000"}}]}, "transfers": []}}, "block_hash": "4e7c0aaf3c3cf8d6ddc3e03847c23170b0ba14cffa31f3411968d222b6b5b40e"}]}`
var moduleBytesDeploy = `{"deploy": {"hash": "00a7445d3be6c6b89308daf62bd055e01d3e96f1a2f6e3efe586dfb915e3dfe2", "header": {"ttl": "1h", "account": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "body_hash": "a14865e1d6016f4c83a240a7df64a538798b2bb392d86674a4e33a59901a354b", "gas_price": 1, "timestamp": "2021-04-08T18:10:32.115Z", "chain_name": "casper-test", "dependencies": []}, "payment": {"ModuleBytes": {"args": [["amount", {"bytes": "05003ad0b814", "parsed": "89000000000", "cl_type": "U512"}]], "module_bytes": ""}}, "session": {"ModuleBytes": {"args": [["public_key", {"bytes": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "parsed": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "cl_type": "PublicKey"}], ["amount", {"bytes": "0500282e8cd1", "parsed": "900000000000", "cl_type": "U512"}], ["delegation_rate", {"bytes": "0a", "parsed": 10, "cl_type": "U8"}]], "module_bytes": ""}}, "approvals": [{"signer": "01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231", "signature": "01d911dea193709debe86db026bb2b9a43e9d0985a3a6f4cde749b55f08b83c12ac9ec7a47b8ae6e7a12cf5dbd96a0314846febf1c9fa95f5aea5a480a170bcc0e"}]}, "api_version": "1.4.6", "execution_results": [{"result": {"Failure": {"cost": "232824230", "effect": {"operations": [{"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "kind": "Read"}, {"key": "balance-2c4bac63bc01ddc6f76e2bc2bcc6af61d6efa9cd65b22785135adea57f98b24c", "kind": "Write"}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "kind": "Read"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "kind": "Read"}, {"key": "balance-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1", "kind": "Write"}], "transforms": [{"key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7", "transform": "Identity"}, {"key": "balance-2c4bac63bc01ddc6f76e2bc2bcc6af61d6efa9cd65b22785135adea57f98b24c", "transform": {"WriteCLValue": {"bytes": "04808f9b95", "parsed": "2510000000", "cl_type": "U512"}}}, {"key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401", "transform": "Identity"}, {"key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6", "transform": "Identity"}, {"key": "balance-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1", "transform": {"AddUInt512": "89000000000"}}]}, "transfers": [], "error_message": "ApiError::AuctionError(4) [64516]"}}, "block_hash": "96b82d76f04b36ba1a83e004e03d862568dec5618620155ca8b53177d415f731"}]}`

func TestResult_GetDeployMetadata(t *testing.T) {
	err := utils.InitViper()
	if err != nil {
		t.Errorf("Unable to init viper : %s", err)
	}
	dt := viper.Get("config")
	log.Println(dt)
	err = mapstructure.Decode(dt, &config.ConfigParsed)
	if err != nil {
		t.Errorf("Unable to init deploy types from the config : %s", err)
	}
	t.Run("Should parse a transfer deploy", func(t *testing.T) {
		var transferResult Result
		err := json.Unmarshal([]byte(transferDeploy), &transferResult)
		if err != nil {
			t.Errorf("Unable to unmarshal transfer deploy : %s", err)
		}
		deployType, metadata := transferResult.GetDeployMetadata()
		if deployType != "transfer" || metadata != `{"id":"","from":"017717a9bb1f07cbb1b6c3afaaad9ff3b8a5b75ea13e5aae6ce33b4b74676c647c","hash":"2d0e59821d67125ab7a07ac719ed6696ce4dd4498ef6f3c283ac7d02f9de7259","amount":"998000000000","target":"e70b850efb68c64e2443da2386452b0d8e4e799362edef0ff56eea8efb114815"}` {
			t.Errorf("Transfer metadata bad parsing detected. Received : %s %s. Expected: %s %s", deployType, metadata, "transfer", ``)
		}
	})
	t.Run("Should parse a storedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByHashDeploy deploy : %s", err)
		}
		deployType, metadata := deployResult.GetDeployMetadata()
		if deployType != "delegate" || metadata != `{"amount":"1000000000","delegator":"01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681","validator":"0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}` {
			t.Errorf("Transfer metadata bad parsing detected. Received : %s %s. Expected: %s %s", deployType, metadata, "delegate", `{"amount":"1000000000","delegator":"01d1379629980bc0be3f2e43e72fe8310cf67c879c47b805c1545ddb76db877681","validator":"0124bfdae2ed128fa5e4057bc398e4933329570e47240e57fc92f5611a6178eba5"}`)
		}
	})
	t.Run("Should parse a storedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByNameDeploy deploy : %s", err)
		}
		deployType, metadata := deployResult.GetDeployMetadata()
		if deployType != "call_faucet" || metadata != `{"amount":"1000000000000","target":"b497711627a79370e1b779dbae5970171c5fcccd3785f1e6593cea0ad6ec7bee"}` {
			t.Errorf("storedContractByNameDeploy metadata bad parsing detected. Received : %s %s. Expected: %s %s", deployType, metadata, "call_faucet", `{"amount":"1000000000000","target":"b497711627a79370e1b779dbae5970171c5fcccd3785f1e6593cea0ad6ec7bee"}`)
		}
	})
	t.Run("Should parse a storedVersionedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByHashDeploy deploy : %s", err)
		}
		deployType, metadata := deployResult.GetDeployMetadata()
		if deployType != "burn" || metadata != `{"owner":{"Account":"account-hash-f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76"},"token_ids":["26"]}` {
			t.Errorf("storedVersionedContractByHashDeploy metadata bad parsing detected. Received : %s %s. Expected: %s %s", deployType, metadata, "burn", `{"owner":"map[Account:account-hash-f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76]","token_ids":"[26]"}`)
		}
	})
	t.Run("Should parse a storedVersionedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByNameDeploy deploy : %s", err)
		}
		deployType, metadata := deployResult.GetDeployMetadata()
		if deployType != "check_balance_of" || metadata != `{"account":{"Account":"account-hash-a3aa343007ff3951f564a754cd3e2f31b2f7332788d02efe874584fbc49ca56d"},"id":"token2","token_contract":"521f4f0319b3f46e325fc09810644e8107c1886540ef996f79e94981f47b31b9"}` {
			t.Errorf("storedVersionedContractByNameDeploy metadata bad parsing detected. Received : %s %s. Expected: %s %s", deployType, metadata, "check_balance_of", `{"account":{"Account":"account-hash-a3aa343007ff3951f564a754cd3e2f31b2f7332788d02efe874584fbc49ca56d"},"id":"token2","token_contract":"521f4f0319b3f46e325fc09810644e8107c1886540ef996f79e94981f47b31b9"}`)
		}
	})
	t.Run("Should parse a moduleBytesDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(moduleBytesDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal moduleBytesDeploy deploy : %s", err)
		}
		deployType, metadata := deployResult.GetDeployMetadata()
		if deployType != "addbid" || metadata != `{"amount":"900000000000","delegation_rate":"10","public_key":"01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231"}` {
			t.Errorf("moduleBytesDeploy metadata bad parsing detected. Received : %s %s. Expected: %s %s", deployType, metadata, "addbid", `{"amount":"900000000000","delegation_rate":"10","public_key":"01624b4b573e42137c9e379ad130c296a46b7e08c1cef7a5c54e0e9ab4f11d0231"}`)
		}
	})
}

func TestResult_GetEvents(t *testing.T) {
	t.Run("Should parse a storedVersionedContractByHashDeploy events", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByHashDeploy deploy : %s", err)
		}
		events := deployResult.GetEvents()
		if events != `{"cep47_burn_one":{"contract_package_hash":"6ca070c78d4eb468b4db4cbc5cadd815c35e15019a841c137372a88d7e247d1d","event_type":"cep47_burn_one","owner":"Key::Account(f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76)","token_id":"26"}}` {
			t.Errorf("deploy events bad parsing detected. Received : %s. Expected: %s ", events, `{"cep47_burn_one":{"contract_package_hash":"6ca070c78d4eb468b4db4cbc5cadd815c35e15019a841c137372a88d7e247d1d","event_type":"cep47_burn_one","owner":"Key::Account(f2d3278a8d24837f23156b812d72ab7ae5ea81467efe2fb718e292756c88cd76)","token_id":"26"}}`)
		}
	})
	t.Run("Should not parse a transferDeploy events", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(transferDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal transferDeploy deploy : %s", err)
		}
		events := deployResult.GetEvents()
		if events != `` {
			t.Errorf("deploy events bad parsing detected. Received : %s. Expected: %s ", events, ``)
		}
	})
}

func TestResult_GetResultAndCost(t *testing.T) {
	t.Run("Should parse result and cost of transferDeploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(transferDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal transferDeploy deploy : %s", err)
		}
		result, cost, err := deployResult.GetResultAndCost()
		if err != nil || cost != "10000" || result != true {
			t.Errorf("deploy cost and result bad parsing detected. Received : %t %s. Expected: %t %s ", result, cost, true, "10000")
		}
	})
	t.Run("Should parse result and cost of storedContractByNameDeploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByNameDeploy deploy : %s", err)
		}
		result, cost, err := deployResult.GetResultAndCost()
		if err != nil || cost != "11406830" || result != false {
			t.Errorf("deploy cost and result bad parsing detected. Received : %t %s. Expected: %t %s ", result, cost, false, "11406830")
		}
	})
}

func TestResult_GetType(t *testing.T) {
	t.Run("Should parse a transfer deploy", func(t *testing.T) {
		var transferResult Result
		err := json.Unmarshal([]byte(transferDeploy), &transferResult)
		if err != nil {
			t.Errorf("Unable to unmarshal transfer deploy : %s", err)
		}
		deployType := transferResult.GetType()
		if deployType != "transfer" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "transfer")
		}
	})
	t.Run("Should parse a storedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByHashDeploy deploy : %s", err)
		}
		deployType := deployResult.GetType()
		if deployType != "storedContractByHash" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "storedContractByHash")
		}
	})
	t.Run("Should parse a storedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByNameDeploy deploy : %s", err)
		}
		deployType := deployResult.GetType()
		if deployType != "storedContractByName" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "storedContractByName")
		}
	})
	t.Run("Should parse a storedVersionedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByHashDeploy deploy : %s", err)
		}
		deployType := deployResult.GetType()
		if deployType != "storedVersionedContractByHash" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "storedVersionedContractByHash")
		}
	})
	t.Run("Should parse a storedVersionedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByNameDeploy deploy : %s", err)
		}
		deployType := deployResult.GetType()
		if deployType != "storedVersionedContractByName" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "storedVersionedContractByName")
		}
	})
	t.Run("Should parse a moduleBytesDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(moduleBytesDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal moduleBytesDeploy deploy : %s", err)
		}
		deployType := deployResult.GetType()
		if deployType != "moduleBytes" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "moduleBytes")
		}
	})
}

func TestResult_GetName(t *testing.T) {
	t.Run("Should parse a storedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByNameDeploy deploy : %s", err)
		}
		deployType := deployResult.GetName()
		if deployType != "faucet" {
			t.Errorf("Deploy type bad parsing detected. Received : %s. Expected: %s", deployType, "faucet")
		}
	})
	t.Run("Should parse a storedVersionedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByHashDeploy deploy : %s", err)
		}
		deployType := deployResult.GetName()
		if deployType != "" {
			t.Errorf("Deploy name bad parsing detected. Received : %s. Expected: %s", deployType, "")
		}
	})
	t.Run("Should parse a storedVersionedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByNameDeploy deploy : %s", err)
		}
		deployType := deployResult.GetName()
		if deployType != "erc1155_test_call" {
			t.Errorf("Deploy name bad parsing detected. Received : %s. Expected: %s", deployType, "erc1155_test_call")
		}
	})
}

func TestResult_GetStoredContractVersion(t *testing.T) {
	t.Run("Should parse a storedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByNameDeploy deploy : %s", err)
		}
		_, err = deployResult.GetStoredContractVersion()
		if err == nil {
			t.Errorf("Received a version on a non versionned deploy")
		}
	})
	t.Run("Should parse a storedVersionedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByHashDeploy deploy : %s", err)
		}
		_, err = deployResult.GetStoredContractVersion()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Should parse a storedVersionedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByNameDeploy deploy : %s", err)
		}
		_, err = deployResult.GetStoredContractVersion()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestResult_GetEntrypoint(t *testing.T) {
	t.Run("Should parse a transfer deploy", func(t *testing.T) {
		var transferResult Result
		err := json.Unmarshal([]byte(transferDeploy), &transferResult)
		if err != nil {
			t.Errorf("Unable to unmarshal transfer deploy : %s", err)
		}
		_, err = transferResult.GetEntrypoint()
		if err == nil {
			t.Errorf("Received entrypoint on a transfer deploy")
		}
	})
	t.Run("Should parse a storedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByHashDeploy deploy : %s", err)
		}
		_, err = deployResult.GetEntrypoint()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestResult_GetStoredContractHash(t *testing.T) {
	t.Run("Should parse a storedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByHashDeploy deploy : %s", err)
		}
		_, err = deployResult.GetStoredContractHash()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Should parse a storedContractByNameDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedContractByNameDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedContractByNameDeploy deploy : %s", err)
		}
		_, err = deployResult.GetStoredContractHash()
		if err == nil {
			t.Errorf("Received hash on named deploy")
		}
	})
	t.Run("Should parse a storedVersionedContractByHashDeploy deploy", func(t *testing.T) {
		var deployResult Result
		err := json.Unmarshal([]byte(storedVersionedContractByHashDeploy), &deployResult)
		if err != nil {
			t.Errorf("Unable to unmarshal storedVersionedContractByHashDeploy deploy : %s", err)
		}
		_, err = deployResult.GetStoredContractHash()
		if err != nil {
			t.Error(err)
		}
	})
}
