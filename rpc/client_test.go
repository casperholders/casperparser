package rpc

import (
	"log"
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

func TestClient_GetContractPackageHash(t *testing.T) {
	_, err := rpcClient.GetContractPackage("394d973ba79e37a455f85fd2e4d7c138c4a6c1a4145fa087d25de59b4a088c6b")
	if err != nil {
		t.Errorf("Unable to retrieve contract package %s", err)
	}
	_, err = rpcClient.GetContractPackage("wronghash")
	if err == nil {
		t.Errorf("Should have thrown an error")
	}
}

func TestClient_GetContract(t *testing.T) {
	r, err := rpcClient.GetContract("db3a41adea55e5ae65c8cba29d8e8527a16ac5fa998a76dfed553215e3254090")
	log.Println(r)
	if err != nil {
		t.Errorf("Unable to retrieve contract package %s", err)
	}
	_, err = rpcClient.GetContract("wronghash")
	if err == nil {
		t.Errorf("Should have thrown an error")
	}
}

func TestClient_GetEraInfo(t *testing.T) {
	r, err := rpcClient.GetEraInfo("3293b31319a97a6451614f57bdd7f65225d4cb2add24fd78af373b4188413a10")
	if err != nil {
		t.Errorf("Unable to retrieve era info %s", err)
	}
	r, err = rpcClient.GetEraInfo("wronghash")
	if r.EraSummary != nil {
		t.Errorf("Should be nil")
	}
}
