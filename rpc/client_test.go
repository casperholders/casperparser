package rpc

import (
	"math"
	"testing"
)

var rpcClient = NewRpcClient("http://rpc.testnet.casperholders.com/rpc")

func TestClient_GetBlock(t *testing.T) {
	_, _, err := rpcClient.GetBlock(1)
	if err != nil {
		t.Errorf("Unable to retrieve block : %s", err)
	}
	_, _, err = rpcClient.GetBlock(math.MaxInt)
	if err == nil {
		t.Errorf("Should have thrown an error")
	}
}

func TestClient_GetLastBlockHeight(t *testing.T) {
	_, err := rpcClient.GetLastBlockHeight()
	if err != nil {
		t.Errorf("Unable to retrieve last block height %s", err)
	}
}

func TestClient_GetDeploy(t *testing.T) {
	_, _, err := rpcClient.GetDeploy("999eebd07739ca44945a102cbfa31e97a486979b83111e1af8ac050c73cc4872")
	if err != nil {
		t.Errorf("Unable to retrieve deploy %s", err)
	}
	_, _, err = rpcClient.GetDeploy("wrongdeploy")
	if err == nil {
		t.Errorf("Should have thrown an error")
	}
}
