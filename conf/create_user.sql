set @@sql_log_bin=0;

#状态检测用户
create user if not exists  admin_op@'%' identified WITH mysql_native_password By 'admin_op';
grant select  on performance_schema.* to admin_op@'%';

set @@sql_log_bin=1;