package file

import (
	"bytes"
	"encoding/csv"
)

type csvOptions struct {
	comma   rune
	comment rune
}

type csvOptionFn struct {
	f func(*csvOptions)
}

type CsvOption interface {
	apply(*csvOptions)
}

// ReadCsv 读取csv文件内容，统一返回[][]string，如果需要特殊格式，自行转换
func ReadCsv(file []byte, opts ...CsvOption) ([][]string, error) {
	return readCsv(file, opts...)
}

// trimBom 如果file是utf-8 bom 格式，将会转换成普通utf-8格式
func trimBom(file []byte) {
	if file[0] == 0xef || file[1] == 0xbb || file[2] == 0xbf {
		file = file[3:]
	}
}

func readCsv(file []byte, opts ...CsvOption) ([][]string, error) {
	opt := defaultCsvOptions()
	for _, o := range opts {
		o.apply(opt)
	}
	trimBom(file)
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
