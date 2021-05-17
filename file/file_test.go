package file

import (
	"fmt"
	"io/ioutil"
	"math"
	"testing"
)

const (
	csvUTF8File    = "../testdata/utf8.csv"
	csvUTF8BomFile = "../testdata/utf8bom.csv"
	csvCRLF        = "../testdata/crlf.csv"
)

type data struct {
	C1 string
	C2 string
}

func TestReadCsvUTF8(t *testing.T) {
	file, err := ioutil.ReadFile(csvUTF8File)
	if err != nil {
		t.Fatalf("io read file [%s] error:%d", csvUTF8File, err)
	}

	opt := []CsvOption{WithCsvComment("\n")}
	testResult1, err := readCsvWithInterface(file, data{}, opt...)
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

	opt := []CsvOption{WithCsvComment("\n")}
	testResult1, err := readCsvWithInterface(file, data{}, opt...)
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

func TestReadCsvCRLF(t *testing.T) {
	file, err := ioutil.ReadFile(csvCRLF)
	if err != nil {
		t.Fatalf("io read file [%s] error:%d", csvUTF8File, err)
	}

	opt := []CsvOption{WithCsvComment("\r\n")}
	testResult1, err := readCsvWithInterface(file, data{}, opt...)
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

func TestReadCsvWithInterface(t *testing.T) {
	type data struct {
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		String  string
		Bool    bool
	}
	table := []data{
		{
			Int:     math.MaxInt64,
			Int8:    math.MaxInt8,
			Int16:   math.MaxInt16,
			Int32:   math.MaxInt32,
			Int64:   math.MaxInt64,
			Uint:    math.MaxUint64,
			Uint8:   math.MaxUint8,
			Uint16:  math.MaxUint16,
			Uint32:  math.MaxUint32,
			Uint64:  math.MaxUint64,
			Float32: float32(3.1415),
			Float64: float64(3.1415926),
			String:  "asdf",
			Bool:    false,
		},
		{
			Int:     math.MinInt64,
			Int8:    math.MinInt8,
			Int16:   math.MinInt16,
			Int32:   math.MinInt32,
			Int64:   math.MinInt64,
			Uint:    0,
			Uint8:   0,
			Uint16:  0,
			Uint32:  0,
			Uint64:  0,
			Float32: float32(-3.1415),
			Float64: float64(-3.1415926),
			String:  "",
			Bool:    true,
		},
	}
	for i := range table {
		d := table[i]
		s := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%v,%v,%s,%t",
			d.Int, d.Int8, d.Int16, d.Int32, d.Int64, d.Uint, d.Uint8, d.Uint16, d.Uint32, d.Uint64, d.Float32, d.Float64, d.String, d.Bool)
		d1 := data{}
		d2, err := readCsvWithInterface([]byte(s), d1)
		if err != nil || len(d2) != 1 {
			t.FailNow()
		}
		d3 := d2[0]
		if d3[0] != d.Int {
			t.Logf("%v:%v", d3[0], d.Int)
		}
		if d3[1] != d.Int8 {
			t.Logf("%v:%v", d3[1], d.Int)
		}
		if d3[2] != d.Int16 {
			t.Logf("%v:%v", d3[2], d.Int)
		}
		if d3[3] != d.Int32 {
			t.Logf("%v:%v", d3[3], d.Int)
		}
		if d3[4] != d.Int64 {
			t.Logf("%v:%v", d3[4], d.Int64)
		}
		if d3[5] != d.Uint {
			t.Logf("%v:%v", d3[5], d.Uint)
		}
		if d3[6] != d.Uint8 {
			t.Logf("%v:%v", d3[6], d.Uint8)
		}
		if d3[7] != d.Uint16 {
			t.Logf("%v:%v", d3[7], d.Uint16)
		}
		if d3[8] != d.Uint32 {
			t.Logf("%v:%v", d3[8], d.Uint32)
		}
		if d3[9] != d.Uint64 {
			t.Logf("%v:%v", d3[9], d.Uint64)
		}
		if d3[10] != d.Float32 {
			t.Logf("%v:%v", d3[10], d.Float32)
		}
		if d3[11] != d.Float64 {
			t.Logf("%v:%v", d3[11], d.Float64)
		}

		if d3[12] != d.String {
			t.Logf("%v:%v", d3[12], d.String)
		}
		if d3[13] != d.Bool {
			t.Logf("%v:%v", d3[13], d.Bool)
		}
	}

}
