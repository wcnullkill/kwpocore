package sql

import (
	// mssql
	"database/sql"
	"errors"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/wcnullkill/kwpocore/conf"
)

func initMssql(dbcfg *conf.DBConfig) (*sql.DB, error) {
	if dbcfg.Driver == "mssql" {
		conStr := getMSSQLStr(dbcfg)
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

// BulkCopy 批量添加，tb第一行就是列名，返回处理行数
// 输入字段说明：
// 	db 数据库
// 	head 表列名
// 	content  要导入的各数据
// 	tbname 表名
// 	presql 执行bulk insert 前需要执行的语句
// 	postsql 执行bulk insert 后将要执行的语句
// 输出字段说明：
//	insert 行数,error
func BulkCopy(db *sql.DB, head []string, content [][]interface{}, tbname string, presql, postsql string) (rowCount int64, err error) {
	txn, err := db.Begin()
	if err != nil {
		return
	}

	if len(content) == 0 {
		err = errors.New("没有数据")
	}

	_, err = db.Exec(presql)
	if err != nil {
		return
	}

	stm, err := txn.Prepare(mssql.CopyIn(tbname, mssql.BulkOptions{}, head...))

	if err != nil {
		return
	}

	for i := 0; i < len(content); i++ {
		_, err = stm.Exec(content[i]...)
		if err != nil {
			return
		}

	}
	result, err := stm.Exec()

	if err != nil {
		return
	}

	err = stm.Close()
	if err != nil {
		return
	}

	err = txn.Commit()
	rowCount, err = result.RowsAffected()

	//后处理
	_, err = db.Exec(postsql)
	if err != nil {
		return
	}
	return
}

func getMSSQLStr(dbcfg *conf.DBConfig) string {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;encrypt=disable;port=%d;database=%s;%s", dbcfg.Host, dbcfg.User, dbcfg.Pwd, dbcfg.Port, dbcfg.DBName, dbcfg.Params)
	return constr
}
