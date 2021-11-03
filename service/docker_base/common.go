package docker_base

import "github.com/IceWhaleTech/CasaOS/model"

//过滤mysql关键字
func MysqlFilter(c MysqlConfig, envs model.EnvArray) model.EnvArray {
	for i := 0; i < len(envs); i++ {
		switch envs[i].Value {
		case "$MYSQL_HOST":
			envs[i].Value = c.DataBaseHost
		case "$MYSQL_PORT":
			envs[i].Value = c.DataBasePort
		case "$MYSQL_USERNAME":
			envs[i].Value = c.DataBaseUser
		case "$MYSQL_PASSWORD":
			envs[i].Value = c.DataBasePassword
		case "$MYSQL_DBNAME":
			envs[i].Value = c.DataBaseDB
		case "$MYSQL_HOST_AND_PORT":
			envs[i].Value = c.DataBaseHost + ":" + c.DataBasePort
		}
	}
	return envs
}
