package app

import (
	"fmt"
	"github.com/wonderivan/logger"
	"os/exec"
	"strings"
)


func Cmd(baseCmd string, args ...interface{}) (output string, err error) {
	var outputByte []byte
	cmdString := fmt.Sprintf(baseCmd, args...)
	logger.Info(cmdString)
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	if outputByte, err = cmd.CombinedOutput(); err != nil {
		output = string(outputByte)
		return
	}
	return string(outputByte), nil
}


func BindVip() (err error) {
	logger.Info("now binding vip : ip addr add dev ....")
	var cmd string
	ip, err := GetCmdPath("ip")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("ip command absolute path : %v", ip)
	if cmd, err = Cmd("sudo %v addr add %v dev %v ", ip, DbMeta.ClusterVip,  DbMeta.Nic); err != nil {
		logger.Error(cmd)
		logger.Error(err.Error())
		return
	}
	logger.Info(cmd)
	arping, err := GetCmdPath("arping")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("arping command absolute path : %v", arping)
	if cmd, err = Cmd("sudo %v -c 4 -A -I %v %v",  arping, DbMeta.Nic,  DbMeta.ClusterVip); err != nil {
		logger.Error(cmd)
		logger.Error(err.Error())
		return
	}
	logger.Info(cmd)
	if !HasLocalIp(DbMeta.ClusterVip) {
		logger.Error("bind vip failed vip=" +  DbMeta.ClusterVip)
		return
	}
	return nil
}


func UnbindVip() (err error) {
	logger.Info("now unbinding vip : ip addr del dev ....")
	ip, err := GetCmdPath("ip")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("ip command absolute path : %v", ip)
	cmd, err := Cmd("sudo %v addr del dev %v %v", ip,  DbMeta.Nic,  DbMeta.ClusterVip)
	logger.Info(cmd)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func GetCmdPath(cmd string)  (absCmd string, err error) {
	absCmd, err = exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	return absCmd, nil
}


func HasLocalIp(ip string) bool {
	if output, err := Cmd("ip addr"); err != nil {
		logger.Error(output)
		return false
	} else {
		return strings.Contains(string(output), ip+"/")
	}
}
