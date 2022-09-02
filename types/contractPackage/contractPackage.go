// Package contractPackage provide a struct for unmarshalling a json contractPackage response from Casper RPC
package contractPackage

type Result struct {
	StoredValue struct {
		ContractPackage struct {
			AccessKey string `json:"access_key"`
			Versions  []struct {
				ProtocolVersionMajor int    `json:"protocol_version_major"`
				ContractVersion      int    `json:"contract_version"`
				ContractHash         string `json:"contract_hash"`
			} `json:"versions"`
			DisabledVersions []interface{} `json:"disabled_versions"`
			Groups           []interface{} `json:"groups"`
		} `json:"ContractPackage"`
	} `json:"stored_value"`
}
