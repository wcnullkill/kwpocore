package conf

import (
	"io/ioutil"
	"testing"
)

const (
	configFile = "../testdata/config.json"
)

func TestReadConfig(t *testing.T) {

	myconfig, err := readConfig()
	if err != nil {
		t.Fatalf("read config error:%d", err)
	}
	db1 := myconfig.DataBases[0]
	if db1.Host != "KWPO-SERVER251\\KWPO251" && db1.Port != 1433 && db1.User != "test_user" && db1.Pwd != "@qq1234" && db1.DBName != "wc_test_db" && db1.Driver != "mssql" && db1.Params != "" {
		t.Log("mssql config content error")
	}
	db2 := myconfig.DataBases[1]
	if db2.Host != "KWPO-SERVER251" && db2.Port != 3306 && db2.User != "test_user" && db2.Pwd != "@qq1234" && db2.DBName != "wc_test_db" && db2.Driver != "mysql" && db2.Params != "timeout=30s" {
		t.Log("mysql config content error")
	}
	mq1 := myconfig.MQs[0]
	if mq1.Host != "10.20.12.80" && mq1.Port != 5672 && mq1.User != "test" && mq1.Pwd != "XSojliSA" && mq1.QueueName != "test" {
		t.Log("mqs config content error")
	}
}

func readConfig() (cfg *MyConfig, err error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	cfg, err = ReadConfig(file)
	if err != nil {
		return
	}
	return
}
