package sql

import (
	"database/sql"
	"errors"

	"github.com/wcnullkill/kwpocore/conf"
)

// New 根据cfg，创建一个sql.db
// 输入参数说明：
//	dbcfg 数据库配置
// 输出参数说明：
//	sql.DB
func New(dbcfg *conf.DBConfig) (*sql.DB, error) {
	if dbcfg.Driver == "mssql" {
		return initMssql(dbcfg)
	} else if dbcfg.Driver == "mysql" {
		return initMysql(dbcfg)
	} else {
		return nil, errors.New("driver error")
	}
}
