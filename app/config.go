package app

import (
	"github.com/wonderivan/logger"
	"gopkg.in/ini.v1"
)

var Cfg *ini.File


func InitConfig(filename string) (err error) {
	if Cfg, err = ini.Load(filename); err != nil {
		logger.Error("load config errors desc=%v", err)
		return
	}
	logger.Info("load config success, config file=%v", filename)
	return
}


func GetCfgString(section, key string) string {
	 return  Cfg.Section(section).Key(key).String()
}


