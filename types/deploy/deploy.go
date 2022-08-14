// Package deploy provide struct and object methods to interact with deploys from the Casper Blockchain
package deploy

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"log"
	"math/big"
	"strconv"
)

// ConfDeployTypes Represent all the types of deploys that will be parsed. This is initialized in the cmd/root.go file with Viper
var ConfDeployTypes ConfigDeployTypes

// GetDeployMetadata Retrieve deploy metadata
func (d Result) GetDeployMetadata() (string, string) {
	if d.Deploy.Session.Transfer != nil {
		return d.getTransferMetadata()
	}
	if d.Deploy.Session.StoredContractByHash != nil ||
		d.Deploy.Session.StoredContractByName != nil ||
		d.Deploy.Session.StoredVersionedContractByHash != nil ||
		d.Deploy.Session.StoredVersionedContractByName != nil {
		return d.ParseStoredContract()
	}
	if d.Deploy.Session.ModuleBytes != nil {
		return d.getModuleByteMetadata()
	}
	return "unknown", ""
}

// getTransferMetadata retrieve transfer metadata
func (d Result) getTransferMetadata() (string, string) {
	values := d.MapArgs()
	var metadata = TransferMetadata{}
	metadata.Hash = d.Deploy.Hash
	metadata.From = d.Deploy.Header.Account
	metadata.ID = ""
	if values["id"] != nil {
		metadata.ID = values["id"].(string)
	}
	metadata.Amount = ""
	if values["amount"] != nil {
		metadata.Amount = values["amount"].(string)
	}
	metadata.Target = ""
	if values["target"] != nil {
		metadata.Target = values["target"].(string)
	}
	metadataString, _ := json.Marshal(metadata)
	return "transfer", string(metadataString)
}

// getModuleByteMetadata retrieve module bytes metadata
func (d Result) getModuleByteMetadata() (string, string) {
	for deployType, argConf := range ConfDeployTypes.ModuleBytes {
		metadata := d.ParseArgs(argConf.StrictArgs, argConf.Args)
		if metadata != "" {
			resolvedDeployType := deployType
			if resolvedDeployType == "stackingOperation" {
				resolvedDeployType = "undelegate"
				_, cost, _ := d.GetResultAndCost()
				bigCost := new(big.Int)
				bigCost.SetString(cost, 10)
				if bigCost.Cmp(big.NewInt(1000000000)) == 1 {
					resolvedDeployType = "delegate"
				}
			}
			return deployType, metadata
		}
	}
	return "moduleBytes", ""
}

// GetType retrieve the deploy session type
func (d Result) GetType() string {
	if d.Deploy.Session.Transfer != nil {
		return "transfer"
	}
	if d.Deploy.Session.StoredContractByHash != nil {
		return "storedContractByHash"
	}
	if d.Deploy.Session.StoredContractByName != nil {
		return "storedContractByName"
	}
	if d.Deploy.Session.StoredVersionedContractByHash != nil {
		return "storedVersionedContractByHash"
	}
	if d.Deploy.Session.StoredVersionedContractByName != nil {
		return "storedVersionedContractByName"
	}
	if d.Deploy.Session.ModuleBytes != nil {
		return "moduleBytes"
	}
	return "unknown"
}

// GetArgs retrieve the args for the deploy
func (d Result) GetArgs() [][]interface{} {
	if d.Deploy.Session.Transfer != nil {
		return d.Deploy.Session.Transfer.Args
	}
	if d.Deploy.Session.StoredContractByHash != nil {
		return d.Deploy.Session.StoredContractByHash.Args
	}
	if d.Deploy.Session.StoredContractByName != nil {
		return d.Deploy.Session.StoredContractByName.Args
	}
	if d.Deploy.Session.StoredVersionedContractByHash != nil {
		return d.Deploy.Session.StoredVersionedContractByHash.Args
	}
	if d.Deploy.Session.StoredVersionedContractByName != nil {
		return d.Deploy.Session.StoredVersionedContractByName.Args
	}
	if d.Deploy.Session.ModuleBytes != nil {
		return d.Deploy.Session.ModuleBytes.Args
	}
	return nil
}

// GetEntrypoint retrieve the entrypoint of the deploy or return an error if none
func (d Result) GetEntrypoint() (string, error) {
	if d.Deploy.Session.StoredContractByHash != nil {
		return d.Deploy.Session.StoredContractByHash.EntryPoint, nil
	}
	if d.Deploy.Session.StoredContractByName != nil {
		return d.Deploy.Session.StoredContractByName.EntryPoint, nil
	}
	if d.Deploy.Session.StoredVersionedContractByHash != nil {
		return d.Deploy.Session.StoredVersionedContractByHash.EntryPoint, nil
	}
	if d.Deploy.Session.StoredVersionedContractByName != nil {
		return d.Deploy.Session.StoredVersionedContractByName.EntryPoint, nil
	}
	return "", fmt.Errorf("deploy %s doesn't have an entrypoint", d.Deploy.Hash)
}

// GetStoredContractHash get the contract hash or return an error if none
func (d Result) GetStoredContractHash() (string, error) {
	if d.Deploy.Session.StoredContractByHash != nil {
		return d.Deploy.Session.StoredContractByHash.Hash, nil
	}
	if d.Deploy.Session.StoredVersionedContractByHash != nil {
		return d.Deploy.Session.StoredVersionedContractByHash.Hash, nil
	}
	return "", fmt.Errorf("deploy %s doesn't have an hash", d.Deploy.Hash)
}

// GetStoredContractVersion get the contract version or return an error if none
func (d Result) GetStoredContractVersion() (int, error) {
	if d.Deploy.Session.StoredVersionedContractByHash != nil {
		return d.Deploy.Session.StoredVersionedContractByHash.Version, nil
	}
	if d.Deploy.Session.StoredVersionedContractByName != nil {
		return d.Deploy.Session.StoredVersionedContractByName.Version, nil
	}
	return 0, fmt.Errorf("deploy %s doesn't have a version", d.Deploy.Hash)
}

// GetName get the contract name or return an empty string if none
func (d Result) GetName() string {
	if d.Deploy.Session.StoredContractByName != nil {
		return d.Deploy.Session.StoredContractByName.Name
	}
	if d.Deploy.Session.StoredVersionedContractByName != nil {
		return d.Deploy.Session.StoredVersionedContractByName.Name
	}
	return ""
}

// GetResultAndCost retrieve the result and cost of a deploy. Return an error if no result found (this can happen when a node is not sync properly)
func (d Result) GetResultAndCost() (bool, string, error) {
	if len(d.ExecutionResults) > 0 {
		var cost string
		var result bool
		if d.ExecutionResults[0].Result.Success != nil {
			cost = d.ExecutionResults[0].Result.Success.Cost
			result = true
		} else {
			cost = d.ExecutionResults[0].Result.Failure.Cost
			result = false
		}
		return result, cost, nil
	}
	return false, "", fmt.Errorf("no result found for deploy : %s", d.Deploy.Hash)
}

// MapArgs maps the arguments of a deploy within a map
func (d Result) MapArgs() map[string]interface{} {
	args := d.GetArgs()
	values := make(map[string]interface{})
	for _, t := range args {
		var value interface{}
		name, ok := t[0].(string)
		if !ok {
			name = t[1].(string)
			value = getValue(t[0].(map[string]interface{})["parsed"])
		} else {
			value = getValue(t[1].(map[string]interface{})["parsed"])
		}
		values[name] = value
	}
	return values
}

// getValue return the string value of an argument value
func getValue(v interface{}) interface{} {
	if unboxed, ok := v.(map[string]interface{}); ok {
		datas := make(map[string]interface{})

		for key, value := range unboxed {
			datas[key] = getValue(value)
		}
		return datas
	}
	if unboxed, ok := v.([]interface{}); ok {
		return unboxed
	}
	if unboxed, ok := v.(map[int]interface{}); ok {
		datas := make(map[string]interface{})

		for key, value := range unboxed {
			datas[fmt.Sprint(key)] = getValue(value)
		}
		return datas
	}
	switch v.(type) {
	case nil:
		return ""
	case bool:
		return strconv.FormatBool(v.(bool))
	case float64:
		return strconv.Itoa(int(v.(float64)))
	case int:
		return strconv.Itoa(v.(int))
	case string:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ParseStoredContract parse a stored contract
func (d Result) ParseStoredContract() (string, string) {
	entrypoint, e := d.GetEntrypoint()
	if e != nil {
		log.Println(e)
		return "unknown", ""
	}
	for _, entrypoints := range ConfDeployTypes.StoredContracts {
		if val, ok := entrypoints[entrypoint]; ok {
			if val.HasName && val.ContractName != d.GetName() {
				log.Println("has name but doesn't equal the current one")
				return "unknown", ""
			}
			parsedArgs := d.ParseArgs(true, val.Args)
			if parsedArgs != "" {
				return val.DeployName, d.ParseArgs(true, val.Args)
			}
		}
	}
	return "unknown", ""
}

// ParseArgs parse the args of a deploy from the user defined configuration
func (d Result) ParseArgs(strict bool, args []string) string {
	deployArgs := d.MapArgs()
	retrievedArgs := make(map[string]interface{})
	if strict && len(deployArgs) != len(args) {
		return ""
	}
	for _, arg := range args {
		if val, ok := deployArgs[arg]; ok {
			retrievedArgs[arg] = val
		} else if strict {
			return ""
		}
	}
	metadataString, _ := json.Marshal(retrievedArgs)
	return string(metadataString)
}

// GetEvents retrieve deploy events
func (d Result) GetEvents() string {
	retrievedEvents := make(map[string]map[string]string)
	var transforms *gabs.Container
	if d.ExecutionResults[0].Result.Success != nil {
		transforms = gabs.Wrap(d.ExecutionResults[0].Result.Success.Effect)
	} else {
		transforms = gabs.Wrap(d.ExecutionResults[0].Result.Failure.Effect)
	}

	for _, child := range transforms.S("transforms").Children() {
		writeCLValue, ok := child.S("transform", "WriteCLValue", "parsed").Data().([]interface{})
		if ok {
			isEvent := false
			tempMap := make(map[string]string)
			for _, mapCLValue := range writeCLValue {
				mapCLValue, converted := mapCLValue.(map[string]interface{})
				if converted {
					key := mapCLValue["key"]
					value := mapCLValue["value"]
					if key != nil {
						key, converted := key.(string)
						if converted {
							if value == nil {
								tempMap[key] = ""
							}
							value, converted := value.(string)
							if converted {
								tempMap[key] = value
							} else {
								tempMap[key] = ""
							}
							if mapCLValue["key"] == "event_type" {
								isEvent = true
							}
						}
					}
				}
			}
			if isEvent {
				retrievedEvents[tempMap["event_type"]] = tempMap
			}
		}
	}
	if len(retrievedEvents) > 0 {
		metadataString, _ := json.Marshal(retrievedEvents)
		return string(metadataString)
	}
	return ""
}

type Result struct {
	Deploy           JsonDeploy        `json:"deploy"`
	ExecutionResults []ExecutionResult `json:"execution_results"`
}

type JsonDeploy struct {
	Hash   string `json:"hash"`
	Header struct {
		Account      string   `json:"account"`
		Timestamp    string   `json:"timestamp"`
		TTL          string   `json:"ttl"`
		GasPrice     int      `json:"gas_price"`
		BodyHash     string   `json:"body_hash"`
		Dependencies []string `json:"dependencies"`
		ChainName    string   `json:"chain_name"`
	} `json:"header"`
	Payment struct {
		ModuleBytes struct {
			ModuleBytes string          `json:"module_bytes"`
			Args        [][]interface{} `json:"args"`
		} `json:"ModuleBytes"`
	} `json:"payment"`
	Session struct {
		Transfer *struct {
			Args [][]interface{} `json:"args"`
		} `json:"Transfer"`
		StoredContractByHash *struct {
			Hash       string          `json:"hash"`
			EntryPoint string          `json:"entry_point"`
			Args       [][]interface{} `json:"args"`
		} `json:"StoredContractByHash"`
		StoredContractByName *struct {
			Name       string          `json:"name"`
			EntryPoint string          `json:"entry_point"`
			Args       [][]interface{} `json:"args"`
		} `json:"StoredContractByName"`
		StoredVersionedContractByHash *struct {
			Hash       string          `json:"hash"`
			Version    int             `json:"version"`
			EntryPoint string          `json:"entry_point"`
			Args       [][]interface{} `json:"args"`
		} `json:"StoredVersionedContractByHash"`
		StoredVersionedContractByName *struct {
			Name       string          `json:"name"`
			Version    int             `json:"version"`
			EntryPoint string          `json:"entry_point"`
			Args       [][]interface{} `json:"args"`
		} `json:"StoredVersionedContractByName"`
		ModuleBytes *struct {
			Args [][]interface{} `json:"args"`
		} `json:"ModuleBytes"`
	} `json:"session"`
	Approvals []struct {
		Signer    string `json:"signer"`
		Signature string `json:"signature"`
	} `json:"approvals"`
}

type ExecutionResult struct {
	BlockHash string `json:"block_hash"`
	Result    struct {
		Success *struct {
			Effect    interface{} `json:"effect"`
			Transfers []string    `json:"transfers"`
			Cost      string      `json:"cost"`
		} `json:"Success"`
		Failure *struct {
			Effect       interface{} `json:"effect"`
			Transfers    []string    `json:"transfers"`
			Cost         string      `json:"cost"`
			ErrorMessage string      `json:"error_message"`
		} `json:"Failure"`
	} `json:"result"`
}

type Effect struct {
	Operations []struct {
		Key  string `json:"key"`
		Kind string `json:"kind"`
	} `json:"operations"`
	Transforms []map[string]interface{} `json:"transforms"`
}

type Transforms []struct {
	Key                string      `json:"key"`
	TransformOperation interface{} `json:"transform"`
}

type TransferMetadata struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	Hash   string `json:"hash"`
	Amount string `json:"amount"`
	Target string `json:"target"`
}

type ConfigDeployTypes struct {
	StoredContracts map[string]map[string]StoredContract
	ModuleBytes     map[string]ModuleByte
}

type StoredContract struct {
	DeployName   string
	HasName      bool
	Args         []string
	ContractName string   `mapstructure:",omitempty"`
	Events       []string `mapstructure:",omitempty"`
}

type ModuleByte struct {
	StrictArgs bool
	Args       []string
	Events     []string `mapstructure:",omitempty"`
}
