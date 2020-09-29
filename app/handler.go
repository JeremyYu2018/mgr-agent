package app

import (
	"github.com/wonderivan/logger"
	"strconv"
	"strings"
	"time"
)

type MysqlCheckHandler struct {
}

func NewMysqlCheckHandler() *MysqlCheckHandler {
	return &MysqlCheckHandler{}
}

//agent判断逻辑:
//(1)每次检查如果本机是Primary，且查询另外两个节点都不是Primary(包括是Secondary，offline,....或者连不上或者没有数据返回等情况)，则绑定vip，否则保持不动。
//(2)每次检查如果本机不是Primary(包括是Secondary，offline,....或者连不上或者没有数据返回等情况)，连另外两个节点，如果任意一个等于Primary, 则卸载vip，否则保持不动。

func (this *MysqlCheckHandler) CheckIsPrimaryLoop() {
	var idx uint64
	ticker := time.NewTicker(time.Duration(DbMeta.PollingInterval) * time.Second)
	for {
		logger.Info("#############entering CheckIsPrimaryLoop Num idx: %v for loop#############", idx)
		var (
			isPrimry bool
			err      error
		)
		sqlDb := OpenMysql(DbMeta)
		isPrimry, err = IsPrimary(sqlDb)
		if isPrimry {
			logger.Info("current mgr node is Primary, check again...")
			//连续判断三次, 如果任意一次判断当前节点不是Primary节点。跳出内部for loop
			for i := 1; i <= 3; i++ {
				logger.Info("%v(st/nd/rd) time local mysql checking member %v:%v state", i, DbMeta.Address, DbMeta.Port)
				isPrimry, err = IsPrimary(sqlDb)
				if err != nil {
					logger.Info("member local mysql %v:%v state, desc=%v", DbMeta.Address, DbMeta.Port, err.Error())
				}
				if isPrimry == false {
					break
				}
				<-ticker.C
			}
			//判断另外两个节点, 如果有一个是Primary节点则跳出此次循环。
			if peerIpIsPrimry, _ := this.CheckPeerMysqlIsPrimary(); peerIpIsPrimry {
				logger.Info("skip this loop")
				//及时关闭，不然可能有连接不释放。在for loop的开头使用defer语句来释放连接好像不行。因为一直都在for loop里面。函数永远退不出去。
				sqlDb.Close()
				continue
			}

		} else {
			//同理如上。
			logger.Info("current mgr node is not Primary, check again...")
			for i := 1; i <= 3; i++ {
				logger.Info("%v(st/nd/rd) time local mysql checking member %v:%v state", i, DbMeta.Address, DbMeta.Port)
				isPrimry, err = IsPrimary(sqlDb)
				if err != nil {
					logger.Info("member local mysql %v:%v state, desc=%v", DbMeta.Address, DbMeta.Port, err.Error())
				}
				if isPrimry == true {
					break
				}
				<-ticker.C
			}
			//判断另外两个节点, 如果都不是Primary则跳出此次循环。
			if peerIpIsPrimry, _ := this.CheckPeerMysqlIsPrimary(); !peerIpIsPrimry {
				logger.Info("skip this loop")
				sqlDb.Close()
				continue
			}
		}

		if isPrimry {
			if HasLocalIp(DbMeta.ClusterVip) {
				logger.Info("cluster vip already binded, skip this loop")
			} else {
				BindVip()
			}
		} else {
			if HasLocalIp(DbMeta.ClusterVip) {
				UnbindVip()
			} else {
				logger.Info("cluster vip already unbinded, skip this loop")
			}
		}
		sqlDb.Close()
		idx++
		<-ticker.C
	}
}

func (this *MysqlCheckHandler) CheckPeerMysqlIsPrimary() (isPrimary bool, err error) {
	logger.Info("#############entering CheckPeerMysqlIsPrimary stage#############")
	peer_ips := GetCfgString("system", "peer_ips")
	ips := strings.Split(peer_ips, ",")
	port, _ := strconv.Atoi(GetCfgString("system", "admin_port"))
	logger.Info("all node ips: %v:%v", ips, port)
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if !HasLocalIp(ip) {
			logger.Info("check peer ip  %v:%v mysql state", ip, port)
			var opDbMeta *OpDbMeta
			if opDbMeta, err = FetchDBMeta(ip); err != nil {
				return false, err
			}
			sqlDb := OpenMysql(opDbMeta)
			isPrimary, _ := IsPrimary(sqlDb)
			sqlDb.Close()
			if isPrimary {
				return true, nil
			}
		}
	}
	return false, nil
}

func FetchDBMeta(address string) (opDbMeta *OpDbMeta, err error) {
	logger.Info("#####entering fetch peers DBMeta stage#######")
	var port int
	if port, err = strconv.Atoi(GetCfgString("system", "admin_port")); err != nil {
		return nil, err
	}
	return &OpDbMeta{
		Username: GetCfgString("system", "admin_user"),
		Password: GetCfgString("system", "admin_password"),
		Address:  address,
		Port:     port,
	}, nil
}