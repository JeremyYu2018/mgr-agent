package app

import "testing"

func TestHasLocalIp(t *testing.T) {
	//
	var hasLocalIp bool
	hasLocalIp = HasLocalIp("192.168.43.35")
	t.Log(hasLocalIp)

	hasLocalIp = HasLocalIp("192.168.43.37")
	t.Log(hasLocalIp)
}


func TestBindVip(t *testing.T) {
	err := BindVip()
	t.Log(err)
}


func TestUnbindVip(t *testing.T) {
	err := UnbindVip()
	t.Log(err)
}