package conf

import (
	"io/ioutil"
	"testing"
)

func TestReadConfig(t *testing.T) {
	bytes, err := ioutil.ReadFile("../testdata/config.json")

	if err != nil {
		t.FailNow()
	}

	myconfig, err := ReadConfig(bytes)
	if err != nil {
		t.FailNow()
	}
	db1 := myconfig.DataBases[0]
	if db1.Host != "KWPO-SERVER251\\KWPO251" && db1.Port != 1433 && db1.User != "test_user" && db1.Pwd != "@qq1234" && db1.DBName != "wc_test_db" && db1.Driver != "mssql" && db1.Params != "" {
		t.Fail()
	}
	db2 := myconfig.DataBases[1]
	if db2.Host != "KWPO-SERVER251" && db2.Port != 3306 && db2.User != "test_user" && db2.Pwd != "@qq1234" && db2.DBName != "wc_test_db" && db2.Driver != "mysql" && db2.Params != "timeout=30s" {
		t.Fail()
	}
	mq1 := myconfig.MQs[0]
	if mq1.Host != "10.20.12.80" && mq1.Port != 5672 && mq1.User != "test" && mq1.Pwd != "XSojliSA" && mq1.QueueName != "test" {
		t.Fail()
	}
}
