package app

import (
	"fmt"
	"testing"
)


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
}


func TestInitDbMeta(t *testing.T) {
	var filename string
	filename = "/Users/yucf/goproject/mgragent/src/mgr-agent/conf/agent.system"
	if err := InitConfig(filename); err != nil {
		t.Errorf("load config error %v \n", err)
	}
	err := InitDbMeta()
	if err != nil {
		t.Errorf("InitDbMeta")
	}
	t.Logf("admin_user :%v\n", DbMeta.Username)
	t.Logf("admin_password :%v\n", DbMeta.Password)
	t.Logf("admin_address :%v\n", DbMeta.Address)
	t.Logf("admin_port :%v\n", DbMeta.Port)
	t.Logf("cluster_vip :%v\n", DbMeta.ClusterVip)
	t.Logf("polling_interval_Seconds :%v\n", DbMeta.PollingInterval)
	t.Logf("nic :%v\n", DbMeta.Nic)
}


func TestOpenMysql(t *testing.T) {
	sqlDB := OpenMysql(DbMeta)
	t.Log(sqlDB)
	if sqlDB == nil {
		t.Errorf("open mysql error %v", sqlDB)
	}
}


func TestPingDB(t *testing.T) {
	sqlDB := OpenMysql(DbMeta)
	alive, err := PingDB(sqlDB)
	t.Logf("mysql alived or not: %v, desc=%v", alive, err)
}


func TestIsPrimary(t *testing.T) {
	sqlDB := OpenMysql(DbMeta)
	res, err := IsPrimary(sqlDB)
	t.Log(res)
	t.Log(err)
}
