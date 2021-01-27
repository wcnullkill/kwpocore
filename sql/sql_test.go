package sql

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/wcnullkill/kwpocore/conf"
)

const (
	configFile = "../testdata/config.json"
)

func TestNewMSSQL(t *testing.T) {
	cfg, err := readConfig()
	if err != nil {
		t.Fatalf("read config error:%d", err)
	}

	db, err := New(&cfg.DataBases[0])
	if err != nil || db == nil {
		t.Errorf("new mssql db error:%d", err)
	}
	defer db.Close()
}

func TestNewMYSQL(t *testing.T) {
	cfg, err := readConfig()
	if err != nil {
		t.Fatalf("read config error:%d", err)
	}
	db, err := New(&cfg.DataBases[1])
	if err != nil || db == nil {
		t.Errorf("new mysql db error:%d", err)
	}
	defer db.Close()
}

func TestMSSQLBulkCopy(t *testing.T) {
	tableName := "test_bulk_table"

	cfg, err := readConfig()
	if err != nil {
		t.Fatalf("read config error:%d", err)
	}

	db, err := New(&cfg.DataBases[0])
	if err != nil || db == nil {
		t.Fatal(err)
	}
	// 需要sql账号role为db_ddladmin
	createSQL := creatTableStr(tableName)
	// 需要sql账号role为db_ddladmin
	dropSQL := fmt.Sprintf("drop table %s", tableName)

	head := []string{"c1", "c2"}
	content := [][]interface{}{
		{"1", "asd"},
		{"2", "1111111111"},
		{"3", "asdf1234"},
		{"4", "哈哈和完全"},
		{"五", "哈哈我"},
	}
	// 需要sql账号role为db_datawriter
	ct, err := BulkCopy(db, head, content, tableName, createSQL, dropSQL)

	if err != nil || int(ct) != len(content) {
		t.Error(err)
	}

}

func creatTableStr(tableName string) (str string) {
	str = `create table ` + tableName + ` (
		c1 nvarchar(100) not null,
		c2 nvarchar(100) not null
	)`
	return
}

func readConfig() (cfg *conf.MyConfig, err error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	cfg, err = conf.ReadConfig(file)
	if err != nil {
		return
	}
	return
}
