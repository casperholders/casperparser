// Package contractPackage provide a struct for unmarshalling a json contractPackage response from Casper RPC
package contract

import (
	"casperParser/types/config"
	"golang.org/x/exp/slices"
	"log"
)

func (c Result) GetContractType() string {
	contractType := "unknown"
	previousCount := 0
	count := 0
	perfectScore := 0
	flatEntrypointArgs := make(map[string][]string)
	for _, entrypoint := range c.StoredValue.Contract.EntryPoints {
		var flatArgs []string
		for _, arg := range entrypoint.Args {
			flatArgs = append(flatArgs, arg.Name)
		}
		flatEntrypointArgs[entrypoint.Name] = flatArgs
	}

	for contractName, contractDefinitions := range config.ConfigParsed.ContractTypes {
		perfectScore = perfectScore + len(contractDefinitions.Entrypoints)
		for _, contractDefEntrypoint := range contractDefinitions.Entrypoints {
			perfectScore = perfectScore + len(contractDefEntrypoint.Args)
			if _, ok := flatEntrypointArgs[contractDefEntrypoint.Name]; ok {
				count++
				for _, arg := range flatEntrypointArgs[contractDefEntrypoint.Name] {
					if slices.Contains(contractDefEntrypoint.Args, arg) {
						count++
					}
				}
			}
		}
		perfectScore = perfectScore + len(contractDefinitions.NamedKeys)
		count = count + c.GetNamedKeysScore(contractDefinitions.NamedKeys)
		if count > previousCount && count >= (perfectScore/4) {
			contractType = contractName
			previousCount = count
		}
		if count > 0 {
			log.Printf("Contract name: %s Score: %d Perfect Score: %d Score NamedKeys: %d Perfect Score NamedKeys: %d\n", contractName, count, perfectScore, c.GetNamedKeysScore(contractDefinitions.NamedKeys), len(contractDefinitions.NamedKeys))
		}
		perfectScore = 0
		count = 0
	}
	return contractType
}

func (c Result) GetNamedKeysScore(namedKeys []string) int {
	count := 0
	for _, namedKey := range c.StoredValue.Contract.NamedKeys {
		if slices.Contains(namedKeys, namedKey.Name) {
			count++
		}
	}
	return count
}

type Result struct {
	StoredValue struct {
		Contract struct {
			ContractPackageHash string       `json:"contract_package_hash"`
			ContractWasmHash    string       `json:"contract_wasm_hash"`
			NamedKeys           []NamedKey   `json:"named_keys"`
			EntryPoints         []Entrypoint `json:"entry_points"`
			ProtocolVersion     string       `json:"protocol_version"`
		} `json:"Contract"`
	} `json:"stored_value"`
}

type NamedKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Entrypoint struct {
	Name           string      `json:"name"`
	Args           []Arg       `json:"args"`
	Ret            interface{} `json:"ret"`
	Access         interface{} `json:"access"`
	EntryPointType string      `json:"entry_point_type"`
}

type Arg struct {
	Name   string      `json:"name"`
	ClType interface{} `json:"cl_type"`
}
