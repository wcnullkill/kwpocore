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
//	insert 行数
func BulkCopy(db *sql.DB, head []string, content [][]interface{}, tbname string, presql, postsql string) int64 {
	txn, err := db.Begin()
	if err != nil {
		panic(err)
	}

	//预处理
	db.Exec(presql)

	if len(content) == 0 {
		panic("没有数据")
	}

	stm, err := txn.Prepare(mssql.CopyIn(tbname, mssql.BulkOptions{}, head...))

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(content); i++ {
		_, err = stm.Exec(content[i]...)
		if err != nil {
			panic(err)
		}

	}
	result, err := stm.Exec()

	if err != nil {
		panic(err)
	}

	err = stm.Close()
	if err != nil {
		panic(err)
	}

	err = txn.Commit()
	rowCount, _ := result.RowsAffected()

	fmt.Printf("批量copy数据%d条", rowCount)

	//后处理
	db.Exec(postsql)

	return rowCount
}

func getMSSQLStr(dbcfg *conf.DBConfig) string {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;encrypt=disable;port=%d;database=%s;%s", dbcfg.Host, dbcfg.User, dbcfg.Pwd, dbcfg.Port, dbcfg.DBName, dbcfg.Params)
	return constr
}
