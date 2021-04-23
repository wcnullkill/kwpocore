// 关于默认值
// int int8 int16 int32 int64 默认值0
// uint uint8 uint16 uint32 uint64 默认值0
// string 暂不支持默认值	todo
// float32 float64 默认值0
// bool 默认值 false
package prot

import (
	"errors"
	"reflect"
	"strconv"
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

// UnMarshal 反序列化
//
func UnMarshal(data []byte, v interface{}, opts ...ProtOption) error {
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

	return d.unMarshal(v, fields)
}

func (d *decodeState) unMarshal(v interface{}, fields []field) error {
	return d.simpleUnMarshal(d.data, v, fields)
}

func (d *decodeState) simpleUnMarshal(data []byte, v interface{}, fields []field) error {
	rows := split(data, d.opt.rowSep)
	va := reflect.ValueOf(v).Elem()
	grow(va, len(rows))
	for i := range rows {
		r := rows[i]
		if len(r) == 0 { //防止整行都是空
			continue
		}
		v1 := va.Index(i)
		cells := split(r, d.opt.columnSep)
		fill(cells, fields, v1)

	}
	return nil

}

func split(data []byte, b byte) [][]byte {
	result := make([][]byte, 0, 10)
	j := 0
	for i := 0; i < len(data); i++ {
		if data[i] == b && i > 0 {
			result = append(result, data[j:i])
			i++
			j = i
		}
	}
	result = append(result, data[j:])
	return result
}

func fill(ss [][]byte, fields []field, v reflect.Value) {
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
			subv.FieldByName(fields[i].name).SetString(string(ss[i]))
		case reflect.Bool:
			subv.FieldByName(fields[i].name).SetBool(stringToBool(ss[i]))
		}

		v.Set(subv)
	}
}

func stringToInt(s []byte) int64 {
	i, _ := strconv.Atoi(string(s))
	return int64(i)
}

func stringToUint(s []byte) uint64 {
	i, _ := strconv.Atoi(string(s))
	return uint64(i)
}
func stringToFloat(s []byte) float64 {
	f, _ := strconv.ParseFloat(string(s), 64)
	return f
}
func stringToBool(s []byte) bool {
	i, _ := strconv.ParseBool(string(s))
	return i
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
