package sql

import (
	"io/ioutil"
	"testing"

	"github.com/wcnullkill/kwpocore/conf"
)

func TestNewMSSQL(t *testing.T) {
	file, err := ioutil.ReadFile("../testdata/config.json")
	if err != nil {
		t.FailNow()
	}

	cfg, err := conf.ReadConfig(file)
	if err != nil {
		t.FailNow()
	}
	db, err := New(&cfg.DataBases[0])
	if err != nil || db == nil {
		t.FailNow()
	}
}

func TestNewMYSQL(t *testing.T) {
	file, err := ioutil.ReadFile("../testdata/config.json")
	if err != nil {
		t.FailNow()
	}

	cfg, err := conf.ReadConfig(file)
	if err != nil {
		t.FailNow()
	}
	db, err := New(&cfg.DataBases[1])
	if err != nil || db == nil {
		t.FailNow()
	}
}
