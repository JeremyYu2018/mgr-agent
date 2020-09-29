package app

import (
	"fmt"
	"testing"
)

var mysqlChekHandler *MysqlCheckHandler
func init() {
	var filename string
	filename = "/Users/yucf/goproject/mgragent/src/mgr-agent/conf/agent.system"
	if err := InitConfig(filename); err != nil {
		fmt.Errorf("load config error %v \n", err)
	}
	err := InitDbMeta()
	if err != nil {
		fmt.Errorf("InitDbMeta")
	}

	mysqlChekHandler = NewMysqlCheckHandler()
}


func TestCheckPeerMysqlState(t *testing.T) {
	state, err := mysqlChekHandler.CheckPeerMysqlIsPrimary()
	t.Log(state)
	t.Log(err)
}