package rpc

import (
	"math"
	"testing"
)

var rpcClient = NewRpcClient("http://node.testnet.casperholders.com:7777/rpc")

func TestClient_GetBlock(t *testing.T) {
	_, _, err := rpcClient.GetBlock(1894552)
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
	_, _, err := rpcClient.GetDeploy("dc9db1422a0c9b212b7dd052033c4ecd54d66423e286ee6790c3af36f2174250")
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
	_, err := rpcClient.GetContract("db3a41adea55e5ae65c8cba29d8e8527a16ac5fa998a76dfed553215e3254090")
	if err != nil {
		t.Errorf("Unable to retrieve contract package %s", err)
	}
	_, err = rpcClient.GetContract("wronghash")
	if err == nil {
		t.Errorf("Should have thrown an error")
	}
}

func TestClient_GetEraInfo(t *testing.T) {
	r, err := rpcClient.GetEraInfo("8d6a98a977482af4eb308cfb4ebdf6981643afdc06f56d6589792808992f56fe")
	if err != nil {
		t.Errorf("Unable to retrieve era info %s", err)
	}
	r, err = rpcClient.GetEraInfo("wronghash")
	if r.EraSummary != nil {
		t.Errorf("Should be nil")
	}
}

func TestClient_GetAuction(t *testing.T) {
	_, err := rpcClient.GetAuction()
	if err != nil {
		t.Errorf("Unable to retrieve auction info %s", err)
	}
}

func TestClient_GetMainPurse(t *testing.T) {
	_, err := rpcClient.GetMainPurse("account-hash-fa12d2dd5547714f8c2754d418aa8c9d59dc88780350cb4254d622e2d4ef7e69")
	if err != nil {
		t.Errorf("Unable to retrieve main purse info %s", err)
	}
}

func TestClient_GetPurseBalance(t *testing.T) {
	_, err := rpcClient.GetPurseBalance("uref-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1-007")
	if err != nil {
		t.Errorf("Unable to retrieve balance %s", err)
	}
}

func TestClient_GetUrefValue(t *testing.T) {
	_, _, err := rpcClient.GetUrefValue("uref-bb9f47c30ddbe192438fad10b7db8200247529d6592af7159d92c5f3aa7716a1-007")
	if err != nil {
		t.Errorf("Unable to retrieve uref %s", err)
	}
	_, _, err = rpcClient.GetUrefValue("hash-d204aaea638a26d580fc0b40af97c468469f3c11c7aa60f2866adc46f03b5033")
	if err != nil {
		t.Errorf("Unable to retrieve uref %s", err)
	}
	initValue, _, err := rpcClient.GetUrefValue("hash-0000000000000000000000000000000000000000000000000000000000000000")
	if err == nil || initValue != "null" {
		t.Errorf("Init value should be null. Received : %s", initValue)
	}
	_, _, err = rpcClient.GetUrefValue("uref-d4a9e949503f14a524ee5a163386aec4ff231b87e4e856f68d8840432ecd693e-007")
	if err != nil {
		t.Errorf("Unable to retrieve uref %s", err)
	}
}
