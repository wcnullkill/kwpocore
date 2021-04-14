package prot

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	scanContinue = iota
	scanObject
	scanValue
	scanEnd
	scanError
)

type decodeState struct {
	data   []byte
	off    int
	fstart int
	opt    protOptions
	fields []field
	opcode int
}

// getDecodeFields 获取有序字段
// 目前要求：struct内每个字段，都设置prot:{num}，{num}为int型，从1开始，且连续
// v类型为*slice类型
func getDecodeFields(v interface{}) ([]field, error) {
	//fmt.Println(reflect.TypeOf(v).Kind())
	if err := decodeValid(v); err != nil {
		return nil, err
	}
	return getFields(reflect.TypeOf(v).Elem())
}

func UnMarshal(data []byte, v interface{}, opts ...ProtOption) error {
	rv := reflect.ValueOf(v)

	if err := decodeValid(reflect.TypeOf(v)); err != nil {
		return err
	}

	opt := defaultOpt()
	for _, o := range opts {
		o.apply(&opt)
	}
	fields, err := getDecodeFields(v)
	if err != nil {
		return err
	}
	d := decodeState{data: data, opt: opt, fields: fields}
	fmt.Println(rv.Kind())
	return d.unMarshal(v, fields)
}

func (d *decodeState) unMarshal(v interface{}, fields []field) error {
	return d.simpleUnMarshal(d.data, v, fields)
}

func (d *decodeState) readRow() []byte {
	i := d.off
	for i < len(d.data) {
		if d.data[i] == d.opt.rowSep {
			i++
			j := d.fstart
			d.fstart = i
			d.off = i
			return d.data[j : i-1]
		}
	}
	return nil
}
func (d *decodeState) scanWhile(op int) {
	i := d.off
	for i < len(d.data) {
		switch d.data[i] {
		case d.opt.rowSep:

		}
	}
}

func (d *decodeState) simpleUnMarshal(data []byte, v interface{}, fields []field) error {
	rows := strings.Split(string(data), string(d.opt.rowSep))
	va := reflect.ValueOf(v).Elem()
	grow(va, len(rows))
	for i, row := range rows {
		v1 := va.Index(i)
		cells := strings.Split(row, string(d.opt.columnSep))
		fill(cells, fields, v1)
	}
	return nil

}

func fill(ss []string, fields []field, v reflect.Value) {
	subv := v
	for i, field := range fields {
		switch field.t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			subv.FieldByName(fields[i].name).SetInt(stringToInt(ss[i]))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			subv.FieldByName(fields[i].name).SetUint(stringToUint(ss[i]))
		case reflect.Float32, reflect.Float64:
			subv.FieldByName(fields[i].name).SetFloat(stringToFloat(ss[i]))
		case reflect.String:
			subv.FieldByName(fields[i].name).SetString(ss[i])
		}
		v.Set(subv)
	}
}

func stringToInt(s string) int64 {
	i, _ := strconv.Atoi(s)
	return int64(i)
}

func stringToUint(s string) uint64 {
	i, _ := strconv.Atoi(s)
	return uint64(i)
}
func stringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func decodeValid(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("必须为pointer类型")
	}
	return nil
}

// grow 初始化时使用
func grow(v reflect.Value, len int) {
	if v.Len() == v.Cap() {
		newCap := len + len/2
		if newCap < 4 {
			newCap = 4
		}
		newv := reflect.MakeSlice(v.Type(), v.Len(), newCap)
		reflect.Copy(newv, v)
		v.Set(newv)
	}
	v.SetLen(len)
}
