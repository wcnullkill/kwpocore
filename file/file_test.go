package file

import (
	"io/ioutil"
	"testing"
)

const (
	csvUTF8File    = "../testdata/utf8.csv"
	csvUTF8BomFile = "../testdata/utf8bom.csv"
)

func TestReadCsvUTF8(t *testing.T) {
	file, err := ioutil.ReadFile(csvUTF8File)
	if err != nil {
		t.Fatalf("io read file [%s] error:%d", csvUTF8File, err)
	}

	testResult1, err := ReadCsv(file)
	if err != nil {
		t.Fatalf("read csv error:%d", err)
	}
	if testResult1[0][0] != "1" || testResult1[0][1] != "asd" {
		t.Log("csv content error")
	}
	if testResult1[1][0] != "2" || testResult1[1][1] != "1111111111" {
		t.Log("csv content error")
	}
	if testResult1[2][0] != "3" || testResult1[2][1] != "asdf1234" {
		t.Log("csv content error")
	}
	if testResult1[3][0] != "4" || testResult1[3][1] != "哈哈和完全" {
		t.Log("csv content error")
	}
	if testResult1[4][0] != "五" || testResult1[4][1] != "哈哈我" {
		t.Log("csv content error")
	}
}

func TestReadCsvUTF8BOM(t *testing.T) {
	file, err := ioutil.ReadFile(csvUTF8BomFile)
	if err != nil {
		t.Fatalf("io read file [%s] error:%d", csvUTF8BomFile, err)
	}

	testResult1, err := ReadCsv(file)
	if err != nil {
		t.Fatalf("read csv error:%d", err)
	}

	if file[0] != 0xef || file[1] != 0xbb || file[2] != 0xbf {
		t.Fatalf("csv [%s] isn't utf8 bom", csvUTF8BomFile)
	}

	if testResult1[0][0] != "1" || testResult1[0][1] != "asd" {
		t.Log("csv content error")
	}
	if testResult1[1][0] != "2" || testResult1[1][1] != "1111111111" {
		t.Log("csv content error")
	}
	if testResult1[2][0] != "3" || testResult1[2][1] != "asdf1234" {
		t.Log("csv content error")
	}
	if testResult1[3][0] != "4" || testResult1[3][1] != "哈哈和完全" {
		t.Log("csv content error")
	}
	if testResult1[4][0] != "五" || testResult1[4][1] != "哈哈我" {
		t.Log("csv content error")
	}
}
