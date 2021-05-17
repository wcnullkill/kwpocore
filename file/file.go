package file

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type decoderFunc func([]byte) interface{}
type decoderCache []decoderFunc

type csvOptions struct {
	comma   rune // 默认","
	comment rune // 默认"/n"
}

type csvOptionFn struct {
	f func(*csvOptions)
}

type CsvOption interface {
	apply(*csvOptions)
}

// ReadCsv 读取csv文件内容，统一返回[][]string
func ReadCsv(file []byte, opts ...CsvOption) ([][]string, error) {
	return readCsv(file, opts...)
}

// trimBom 如果file是utf-8 bom 格式，将会转换成普通utf-8格式
func trimBom(file []byte) []byte {
	if file[0] == 0xef || file[1] == 0xbb || file[2] == 0xbf {
		return file[3:]
	}
	return file
}

func readCsv(file []byte, opts ...CsvOption) ([][]string, error) {
	opt := defaultCsvOptions()
	for _, o := range opts {
		o.apply(opt)
	}
	file = trimBom(file)
	b := bytes.NewReader(file)
	r := csv.NewReader(b)
	r.Comma = opt.comma
	if opt.comment != '\n' {
		r.Comment = opt.comment
	}
	content, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (fo *csvOptionFn) apply(opts *csvOptions) {
	fo.f(opts)
}

func WithCsvComma(comma rune) *csvOptionFn {
	return &csvOptionFn{func(co *csvOptions) {
		co.comma = comma
	}}
}

func WithCsvComment(comment rune) *csvOptionFn {
	return &csvOptionFn{func(co *csvOptions) {
		co.comment = comment
	}}
}

func defaultCsvOptions() *csvOptions {
	return &csvOptions{
		comma:   ',',
		comment: '\n',
	}
}

// ReadCsvWithInterface 读取csv文件内容，返回[][]interface{},interface{}z中包含struct字段的真实类型
// 方便使用mssql bulkcopy
func ReadCsvWithInterface(file []byte, v interface{}, opts ...CsvOption) ([][]interface{}, error) {
	return readCsvWithInterface(file, v, opts...)
}

func readCsvWithInterface(file []byte, v interface{}, opts ...CsvOption) ([][]interface{}, error) {
	opt := defaultCsvOptions()
	for _, o := range opts {
		o.apply(opt)
	}
	file = trimBom(file)
	rows := bytes.Split(file, []byte{byte(opt.comment)})
	if len(rows) == 0 {
		return [][]interface{}{}, nil
	}
	cache, err := scan(v)
	if err != nil {
		return nil, err
	}
	result := make([][]interface{}, len(rows))
	for i := range rows {
		row := make([]interface{}, len(cache))
		cells := bytes.Split(rows[i], []byte{byte(opt.comma)})
		for j := range cells {
			row[j] = cache[j](cells[j])
		}
		result[i] = row
	}
	return result, nil
}

func scan(v interface{}) (decoderCache, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("传入参数必须要struct类型")
	}
	num := t.NumField()
	if num == 0 {
		return nil, errors.New("传入参数必须要包含可访问字段")
	}
	cache := make([]decoderFunc, num)
	for i := 0; i < num; i++ {
		f := t.Field(i)
		// interface{}中，不能只包含值，还需要有类型
		// 从csv中读取出来的数据，都是文本形式，不是二进制形式，不可以使用binary包，结果会错误
		switch f.Type.Kind() {
		case reflect.Int:
			cache[i] = decoderInt
		case reflect.Int8:
			cache[i] = decoderInt8
		case reflect.Int16:
			cache[i] = decoderInt16
		case reflect.Int32:
			cache[i] = decoderInt32
		case reflect.Int64:
			cache[i] = decoderInt64
		case reflect.Uint:
			cache[i] = decoderUint
		case reflect.Uint8:
			cache[i] = decoderUint8
		case reflect.Uint16:
			cache[i] = decoderUint16
		case reflect.Uint32:
			cache[i] = decoderUint32
		case reflect.Uint64:
			cache[i] = decoderUint64
		case reflect.Float32:
			cache[i] = decoderFloat32
		case reflect.Float64:
			cache[i] = decoderFloat64
		case reflect.String:
			cache[i] = decoderString
		case reflect.Bool:
			cache[i] = decoderBool
		}
	}
	return cache, nil
}

func decoderInt(b []byte) interface{} {
	return int(decoderInt64R(b))
}

func decoderInt8(b []byte) interface{} {
	return int8(decoderInt64R(b))
}

func decoderInt16(b []byte) interface{} {
	return int16(decoderInt64R(b))
}

func decoderInt32(b []byte) interface{} {
	return int32(decoderInt64R(b))
}
func decoderInt64(b []byte) interface{} {
	return decoderInt64R(b)
}

func decoderInt64R(b []byte) int64 {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(fmt.Errorf("convert %s to int error:%v", string(b), err))
	}
	return i
}

func decoderUint(b []byte) interface{} {
	return uint(decoderUint64R(b))
}

func decoderUint8(b []byte) interface{} {
	return uint8(decoderUint64R(b))
}

func decoderUint16(b []byte) interface{} {
	return uint16(decoderUint64R(b))
}

func decoderUint32(b []byte) interface{} {
	return uint32(decoderUint64R(b))
}

func decoderUint64(b []byte) interface{} {
	return decoderUint64R(b)
}

func decoderUint64R(b []byte) uint64 {
	i, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Errorf("convert %s to uint error:%v", string(b), err))
	}
	return i
}

func decoderFloat32(b []byte) interface{} {
	f, err := strconv.ParseFloat(string(b), 32)
	if err != nil {
		panic(fmt.Errorf("convert %s to float32 error:%v", string(b), err))
	}
	return float32(f)
}

func decoderFloat64(b []byte) interface{} {
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		panic(fmt.Errorf("convert %s to float64 error:%v", string(b), err))
	}
	return f
}

func decoderString(b []byte) interface{} {
	return string(b[:])
}

func decoderBool(b []byte) interface{} {
	v, _ := strconv.ParseBool(string(b))
	return v
}
