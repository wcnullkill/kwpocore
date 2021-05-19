package sql

import (
	"context"
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
func BulkCopy(db *sql.DB, head []string, content [][]interface{}, tbname string, presql, postsql string) (int64, error) {
	return BulkCopyWithCtx(context.Background(), db, head, content, tbname, presql, postsql)
}

func BulkCopyWithCtx(ctx context.Context, db *sql.DB, head []string, content [][]interface{}, tbname string, presql, postsql string) (int64, error) {
	txn, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	if len(content) == 0 {
		err = errors.New("没有数据")
		return 0, err
	}

	_, err = txn.Exec(presql)
	if err != nil {
		return 0, err
	}

	stm, err := txn.Prepare(mssql.CopyIn(tbname, mssql.BulkOptions{}, head...))

	if err != nil {
		return 0, err
	}

	for i := 0; i < len(content); i++ {
		_, err = stm.Exec(content[i]...)
		if err != nil {
			return 0, err
		}

	}
	result, err := stm.Exec()
	if err != nil {
		return 0, err
	}

	err = stm.Close()
	if err != nil {
		return 0, err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	//后处理
	_, err = txn.Exec(postsql)
	if err != nil {
		return 0, err
	}
	err = txn.Commit()
	if err != nil {
		return 0, err
	}
	return rowCount, err
}

func getMSSQLStr(dbcfg *conf.DBConfig) string {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;encrypt=disable;port=%d;database=%s;%s", dbcfg.Host, dbcfg.User, dbcfg.Pwd, dbcfg.Port, dbcfg.DBName, dbcfg.Params)
	return constr
}
