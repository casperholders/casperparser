// Package rpc client wrapper for easy rpc call to a casper node
package rpc

import (
	"bytes"
	"casperParser/types/auction"
	"casperParser/types/block"
	"casperParser/types/contract"
	"casperParser/types/contractPackage"
	"casperParser/types/deploy"
	"casperParser/types/reward"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var cachedStateRootHash = ""
var cachedStateRootHashTime = time.Now()

type Client struct {
	endpoint string
}

// NewRpcClient return a rpc client
func NewRpcClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// RpcCall make a rpc call
func (c *Client) RpcCall(method string, params interface{}) (Response, error) {
	body, err := json.Marshal(Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
	})

	if err != nil {
		return Response{}, fmt.Errorf("failed to marshal json %w", err)
	}

	resp, err := http.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return Response{}, fmt.Errorf("failed to make request: %w", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to get response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Response{}, fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}

	var rpcResponse Response
	err = json.Unmarshal(b, &rpcResponse)
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to parse response body: %w", err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse, fmt.Errorf("rpc call failed, code - %d, message - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse, nil
}

// GetBlock from the casper blockchain
func (c *Client) GetBlock(height int) (block.Result, json.RawMessage, error) {
	resp, err := c.RpcCall("chain_get_block",
		blockParams{blockIdentifier{
			Height: uint64(height),
		}})

	if err != nil {
		return block.Result{}, json.RawMessage{}, err
	}
	var result block.Result
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return block.Result{}, json.RawMessage{}, fmt.Errorf("failed to get result: %w", err)
	}
	return result, resp.Result, nil
}

// GetLastBlockHeight from the casper blockchain
func (c *Client) GetLastBlockHeight() (int, error) {
	resp, err := c.RpcCall("chain_get_block", nil)

	if err != nil {
		return -1, err
	}

	var result block.Result
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return -1, fmt.Errorf("failed to get result: %w", err)
	}
	return result.Block.Header.Height, nil
}

// GetAuction from the casper blockchain
func (c *Client) GetAuction() (auction.Result, error) {
	resp, err := c.RpcCall("state_get_auction_info", nil)

	if err != nil {
		return auction.Result{}, err
	}

	var result auction.Result
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return auction.Result{}, fmt.Errorf("failed to get result: %w", err)
	}
	return result, nil
}

// GetDeploy from the casper blockchain
func (c *Client) GetDeploy(hash string) (deploy.Result, json.RawMessage, error) {
	resp, err := c.RpcCall("info_get_deploy", map[string]string{
		"deploy_hash": hash,
	})

	if err != nil {
		return deploy.Result{}, json.RawMessage{}, err
	}
	var result deploy.Result
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return deploy.Result{}, json.RawMessage{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, resp.Result, nil
}

// GetContractPackage from the casper blockchain
func (c *Client) GetContractPackage(hash string) (string, error) {
	srh, err := c.GetStateRootHash(false)
	if err != nil {
		return "", fmt.Errorf("failed to get result: %w", err)
	}

	resp, err := c.RpcCall("state_get_item", []string{srh, "hash-" + hash})
	if err != nil {
		return "", err
	}
	var cp contractPackage.Result
	err = json.Unmarshal(resp.Result, &cp)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(cp.StoredValue)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetStateRootHash from the casper blockchain
func (c *Client) GetStateRootHash(cache bool) (string, error) {
	now := time.Now()
	if cachedStateRootHash == "" || !cache || now.Sub(cachedStateRootHashTime).Seconds() > 32 {
		resp, err := c.RpcCall("chain_get_state_root_hash", nil)
		if err != nil {
			return "", err
		}
		var result stateRootHash
		err = json.Unmarshal(resp.Result, &result)
		if err != nil {
			return "", fmt.Errorf("failed to get result: %w", err)
		}

		cachedStateRootHash = result.StateRootHash
		cachedStateRootHashTime = time.Now()
	}
	return cachedStateRootHash, nil
}

// GetMainPurse from the casper blockchain
func (c *Client) GetMainPurse(hash string) (string, error) {
	srh, err := c.GetStateRootHash(false)
	if err != nil {
		return "", fmt.Errorf("failed to get result: %w", err)
	}
	resp, err := c.RpcCall("state_get_item", []string{srh, hash})
	if err != nil {
		return "", err
	}

	var result mainPurse
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return "", fmt.Errorf("failed to get result: %w", err)
	}
	log.Println(result)
	return result.StoredValue.Account.MainPurse, nil
}

// GetPurseBalance from the casper blockchain
func (c *Client) GetPurseBalance(hash string) (string, error) {
	srh, err := c.GetStateRootHash(false)
	if err != nil {
		return "", fmt.Errorf("failed to get result: %w", err)
	}
	resp, err := c.RpcCall("state_get_balance", []string{srh, hash})
	if err != nil {
		return "", err
	}
	var result purseBalance
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return "", fmt.Errorf("failed to get result: %w", err)
	}
	return result.BalanceValue, nil
}

// GetContract from the casper blockchain
func (c *Client) GetContract(hash string) (contract.Result, error) {
	srh, err := c.GetStateRootHash(false)
	if err != nil {
		return contract.Result{}, fmt.Errorf("failed to get result: %w", err)
	}

	resp, err := c.RpcCall("state_get_item", []string{srh, "hash-" + hash})
	if err != nil {
		return contract.Result{}, err
	}
	var contractParsed contract.Result
	err = json.Unmarshal(resp.Result, &contractParsed)
	if err != nil {
		return contract.Result{}, err
	}
	return contractParsed, nil
}

// GetEraInfo from the casper blockchain
func (c *Client) GetEraInfo(hash string) (reward.Result, error) {
	resp, err := c.RpcCall("chain_get_era_info_by_switch_block", []map[string]string{{
		"Hash": hash,
	}})
	if err != nil {
		return reward.Result{}, err
	}
	var rewardParsed reward.Result
	err = json.Unmarshal(resp.Result, &rewardParsed)
	if err != nil {
		return reward.Result{}, err
	}
	return rewardParsed, nil
}

// GetUrefValue from the casper blockchain
func (c *Client) GetUrefValue(hash string) (string, bool, error) {
	srh, err := c.GetStateRootHash(true)
	if err != nil {
		return "", false, fmt.Errorf("failed to get result: %w", err)
	}

	resp, err := c.RpcCall("state_get_item", []string{srh, hash})
	if err != nil {
		return "", false, err
	}
	var parsedUref uref
	err = json.Unmarshal(resp.Result, &parsedUref)
	if err != nil {
		return "", false, err
	}
	if parsedUref.StoredValue.CLValue.Parsed == nil {
		balance, errB := c.GetPurseBalance(hash)
		if errB == nil {
			return balance, true, nil
		}
	}
	b, err := json.Marshal(parsedUref.StoredValue.CLValue.Parsed)
	if err != nil {
		return "", false, err
	}
	return string(b), false, nil
}

type Request struct {
	Version string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Response struct {
	Version string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type blockParams struct {
	BlockIdentifier blockIdentifier `json:"block_identifier"`
}

type blockIdentifier struct {
	Hash   string `json:"Hash,omitempty"`
	Height uint64 `json:"Height,omitempty"`
}

type stateRootHash struct {
	StateRootHash string `json:"state_root_hash"`
}

type mainPurse struct {
	StoredValue struct {
		Account struct {
			MainPurse string `json:"main_purse"`
		} `json:"Account"`
	} `json:"stored_value"`
}

type uref struct {
	StoredValue struct {
		CLValue struct {
			Parsed interface{} `json:"parsed"`
		} `json:"CLValue"`
	} `json:"stored_value"`
}

type purseBalance struct {
	BalanceValue string `json:"balance_value"`
}

type urefValue struct {
	Value interface{} `json:"value"`
}
