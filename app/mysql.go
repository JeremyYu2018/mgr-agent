package app

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wonderivan/logger"
	"strconv"
)

type OpDbMeta struct {
	//操作数据库的用户名相关信息
	Username        string
	Password        string
	Address         string
	Port            int
	PollingInterval int
	ClusterVip      string
	Nic             string
}

var DbMeta *OpDbMeta

func InitDbMeta() (err error) {
	var (
		port                   int
		pollingIntervalSeconds int
	)
	pollingIntervalSeconds, err = strconv.Atoi(GetCfgString("system", "polling_interval_Seconds"))
	if err != nil {
		logger.Error(err)
		return
	}

	if port, err = strconv.Atoi(GetCfgString("system", "admin_port")); err != nil {
		logger.Error(err)
		return
	}
	DbMeta = &OpDbMeta{
		Username:        GetCfgString("system", "admin_user"),
		Password:        GetCfgString("system", "admin_password"),
		Address:         GetCfgString("system", "admin_address"),
		Port:            port,
		PollingInterval: pollingIntervalSeconds,
		ClusterVip:      GetCfgString("system", "cluster_vip"),
		Nic:             GetCfgString("system", "nic"),
	}
	return
}

func OpenMysql(db *OpDbMeta) (sqlDB *sqlx.DB) {
	logger.Info("try to connect to mysql %v:%v", db.Address, db.Port)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", db.Username, db.Password, db.Address, db.Port)
	sqlDB, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.Error("open mysql connection %v:%v failed", db.Address, db.Port)
		return nil
	}
	logger.Info("open mysql connection %v:%v succeed", db.Address, db.Port)
	return
}

func PingDB(SqlDB *sqlx.DB) (alive bool, err error) {
	if err := SqlDB.Ping(); err != nil {
		logger.Error(err.Error())
		return false, err
	}
	logger.Info("ping success")
	return true, nil
}

func IsPrimary(SqlDB *sqlx.DB) (bool, error) {
	if err := SqlDB.Ping(); err != nil {
		logger.Error(err.Error())
		return false, err
	}
	logger.Info("ping success")
	sqlStr := `SELECT CASE 
		WHEN a.member_id = b.primary_host THEN 1 
		ELSE 0
	END AS ifPrimary
FROM performance_schema.replication_group_members a, (
		SELECT variable_value AS primary_host
		FROM performance_schema.global_status
		WHERE variable_name = 'group_replication_primary_member'
	) b
WHERE a.member_id = @@server_uuid;`
	res, err := SqlQuery1FieldAnd1Row(SqlDB, "ifPrimary", sqlStr)
	if err != nil {
		return false, err
	}
	logger.Info("ifPrimary=" + res["ifPrimary"])
	role, err := strconv.Atoi(res["ifPrimary"])
	if err != nil {
		return false, err
	}
	switch role {
	case 1:
		return true, nil
	case 0:
		return false, errors.New("not primary node")
	default:
		return false, errors.New("unkown reason, please check mysql performance_schema.replication_group_members")
	}
}

func SqlQuery1FieldAnd1Row(sqldb *sqlx.DB, field string, sqlStr string) (result map[string]string, err error) {
	result = make(map[string]string)
	logger.Info("exec sql: %v", sqlStr)
	var rows *sql.Rows
	if rows, err = sqldb.Query(sqlStr); err != nil {
		return nil, err
	}
	if rows == nil {
		logger.Error(field + ": Empty set (0.00 sec)")
	}
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
		result[field] = value
	}
	rows.Close()
	return result, nil
}


