package sql

import (
	"database/sql"
	"errors"
	"fmt"

	//mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/wcnullkill/kwpocore/conf"
)

func initMysql(dbcfg *conf.DBConfig) (*sql.DB, error) {
	if dbcfg.Driver == "mysql" {
		conStr := getMYSQLStr(dbcfg)
		sqldb, err := sql.Open(dbcfg.Driver, conStr)
		if err != nil {
			return nil, err
		}

		//判断连接是否有效
		err = sqldb.Ping()
		if err != nil {
			return nil, err
		}

		return sqldb, nil
	}
	return nil, errors.New("driver error")
}

func getMYSQLStr(dbcfg *conf.DBConfig) string {
	var p string
	if dbcfg.Params != "" {
		p = "?" + dbcfg.Params
	}
	constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", dbcfg.User, dbcfg.Pwd, dbcfg.Host, dbcfg.Port, dbcfg.DBName, p)
	return constr
}
