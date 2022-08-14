// Package rpc client wrapper for easy rpc call to a casper node
package rpc

import (
	"bytes"
	"casperParser/types/block"
	"casperParser/types/deploy"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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
