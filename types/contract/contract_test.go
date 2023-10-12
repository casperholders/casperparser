package contract_test

import (
	"casperParser/rpc"
	"casperParser/types/config"
	"casperParser/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"testing"
)

var rpcClient = rpc.NewRpcClient("https://node.testnet.casperholders.com/rpc")

func TestResult_GetContractType(t *testing.T) {
	err := utils.InitViper()
	if err != nil {
		t.Errorf("Unable to init viper : %s", err)
	}
	dt := viper.Get("config")
	log.Println(dt)
	err = mapstructure.Decode(dt, &config.ConfigParsed)
	log.Println(config.ConfigParsed)
	r, err := rpcClient.GetContract("31bfdc9591902bda8f921d6c31f3e974bda18ec5614222f3bec55390decd05a0")
	println(r.GetContractTypeAndScore())
}
