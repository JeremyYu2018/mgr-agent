mysql 8.0 MGR（GroupReplication） vip 切换工具


### 配置说明
```bash
[system]
cluster_vip              =  192.168.43.38
nic                      =  en0
polling_interval_Seconds =  5
admin_user      =  admin_op
admin_password  =  admin_op
admin_address   =  172.16.200.12
#尽量填写MySQL 8.0管理端口, 管理端口不受max_connections的限制。
admin_port      =  5000
#peer_ips同时也会作为group_replication_group_seeds的值。
peer_ips        =  172.16.200.12,172.16.200.13,172.16.200.14
```
- cluster_vip  集群的VIP（三个节点可以不配置一致)；

- polling_interval_Seconds:  检测节点角色的时间频率，单位为秒(S)；

- nic 网卡名，当MGR接到变为Primary时会把ClusterVip绑定在此网卡上；

- peer_ips  MGR三节点的IP地址。其中之一的地址需要在本机上；

admin_user，admin_password，admin_address，admin_port-agent程序用来连接MySQL检测主从角色的用户；

### 运行条件
1. MySQL 8.0版本以上
2. MGR单主模式（Single-Primary-Mode）
3. 原生用户认证模式(mysql_native_password)
在create_user.sql里面有创建用户的语句。


```bash
create user if not exists  admin_op@'%' identified WITH mysql_native_password By 'admin_op';
grant select  on performance_schema.* to admin_op@'%';
```


### 运行方式

帮助：目前只简单提供了两个子命令，启动和关闭。
```bash
# ./bin/mgr-agent -h
mgr-agent 

A Virtual IP failover agent tool for MySQL Group Replication(Single Primary Mode):

Usage:
  mgr-agent [flags]
  mgr-agent [command]

Available Commands:
  help        Help about any command
  start       start mgr-agent
  stop        Stop mgr-agent

Flags:
  -c, --config string   configuration file to use (default "conf/agent.system")
  -h, --help            help for mgr-agent
  -v, --version         version for mgr-agent

Use "mgr-agent [command] --help" for more information about a command.
```

- 前台启动
```bash
# ./bin/mgr-agent start
2020/09/23 15:01:47 mgr-agent start
2020-09-23 15:01:47 [INFO] [mgr-agent/app/config.go:15] load config success
```
- 后台启动
```bash
# ./bin/mgr-agent start -d
2020/09/23 15:02:13 ./bin/mgr-agent start, [PID] %d running...
 28420
```

- 关闭
```$xslt
./bin/mgr-agent stop
2020-09-23 15:01:20 [INFO] [mgr-agent/main.go:94] mgr-agent stopped
```

### 原理
  agent判断逻辑:
   - (1)每次检查如果本机是Primary，且查询另外两个节点都不是Primary(包括是Secondary，offline,....或者连不上或者没有数据返回等情况)，则绑定vip，否则保持不动。
   - (2)每次检查如果本机不是Primary(包括是Secondary，offline,....或者连不上或者没有数据返回等情况)，连另外两个节点，如果任意一个等于Primary, 则卸载vip，否则保持不动。
 > 注意事项：
   1. 默认情况下，节点退出MGR集群，会变为超级只读（super_read_only=1）,从退出集群到vip正式被卸载存在检测时间，为防止数据混乱，业务用户切不可有super权限。
   2. 实际上，MySQL 5.7版本也可使用此工具，但是考虑到5.7MGR的不稳定性，推荐使用8.0版本。
   3. MGR的部署请自行遵循最佳实践。
