package file

import (
	"io/ioutil"
	"testing"
)

func TestReadCsvUTF8(t *testing.T) {
	file, err := ioutil.ReadFile("../testdata/utf8.csv")
	if err != nil {
		t.FailNow()
	}

	testResult1, err := ReadCsv(file, ';')
	if err != nil {
		t.FailNow()
	}
	if testResult1[0][0] != "1" || testResult1[0][1] != "asd" {
		t.Fail()
	}
	if testResult1[1][0] != "2" || testResult1[1][1] != "1111111111" {
		t.Fail()
	}
	if testResult1[2][0] != "3" || testResult1[2][1] != "asdf1234" {
		t.Fail()
	}
	if testResult1[3][0] != "4" || testResult1[3][1] != "哈哈和完全" {
		t.Fail()
	}
	if testResult1[4][0] != "五" || testResult1[4][1] != "哈哈我" {
		t.Fail()
	}
}

func TestReadCsvUTF8BOM(t *testing.T) {
	file, err := ioutil.ReadFile("../testdata/utf8bom.csv")
	if err != nil {
		t.FailNow()
	}

	testResult1, err := ReadCsv(file, ';')
	if err != nil {
		t.FailNow()
	}

	if file[0] != 0xef || file[1] != 0xbb || file[2] != 0xbf {
		t.FailNow()
	}

	if testResult1[0][0] != "1" || testResult1[0][1] != "asd" {
		t.Fail()
	}
	if testResult1[1][0] != "2" || testResult1[1][1] != "1111111111" {
		t.Fail()
	}
	if testResult1[2][0] != "3" || testResult1[2][1] != "asdf1234" {
		t.Fail()
	}
	if testResult1[3][0] != "4" || testResult1[3][1] != "哈哈和完全" {
		t.Fail()
	}
	if testResult1[4][0] != "五" || testResult1[4][1] != "哈哈我" {
		t.Fail()
	}
}
