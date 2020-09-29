package app

import "testing"

func TestInitConfig(t *testing.T) {
	var filename string
	filename = "/Users/yucf/goproject/mgragent/src/mgr-agent/conf/agent.system"
	if err := InitConfig(filename); err != nil {
		t.Errorf("load config error %v \n", err)
	}

	filename = "/do/not/exists/agent.system"
	if err := InitConfig(filename); err == nil {
		t.Errorf("load config error %v \n", err)
	}
}

func TestGetCfgString(t *testing.T) {
	var filename string
	filename = "/Users/yucf/goproject/mgragent/src/mgr-agent/conf/agent.system"
	if err := InitConfig(filename); err != nil {
		t.Errorf("load config error %v \n", err)
	}

	pollingIntervals  := GetCfgString("system", "polling_interval_Seconds")
	if pollingIntervals != "5" {
		t.Errorf("get value error PollingIntervalSeconds=%v\n", pollingIntervals)
	}

	admin_user  := GetCfgString("system", "admin_user")
	t.Log(admin_user)
	if admin_user != "admin_op" {
		t.Errorf("get value error admin_user=%v\n", admin_user)
	}

	peer_ips := GetCfgString("system", "peer_ips")
	if peer_ips != "172.16.200.12,172.16.200.13,172.16.200.14" {
		t.Errorf("%v\n", peer_ips)
	}

}