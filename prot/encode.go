package prot

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagName string = "prot"
)

type encodeState struct {
	out bytes.Buffer
}

type field struct {
	name string
	t    reflect.Type
}

// getEncodeFields 获取有序字段
// 目前要求：struct内每个字段，都设置prot:{num}，{num}为int型，从1开始，且连续
// v类型为slice类型
func getEncodeFields(v interface{}) ([]field, error) {
	//fmt.Println(reflect.TypeOf(v).Kind())
	if err := encodeValid(v); err != nil {
		return nil, err
	}
	return getFields(reflect.TypeOf(v))
}

// getFields 传入Type slice
func getFields(v reflect.Type) ([]field, error) {

	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("getFields方法要求传入参数必须为slice类型，当前类型为%s", v.Kind())
	}

	va := v.Elem()

	fieldNum := va.NumField()

	fields := make([]field, fieldNum)

	// 将prot的值填充tags
	for i := 0; i < fieldNum; i++ {
		s := va.Field(i).Tag.Get(tagName)
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%s不能转换为int", s)
		}
		if n > fieldNum {
			return nil, errors.New("prot值大于字段数量")
		}
		if len(fields[n-1].name) > 0 {
			return nil, errors.New("prot有重复")
		}
		fields[n-1] = field{va.Field(i).Name, va.Field(i).Type}
	}

	//如果tags中有空字符串，则认为prot值有重复
	for _, tag := range fields {
		if len(tag.name) == 0 {
			return nil, errors.New("prot值不连续")
		}
	}
	return fields, nil
}

func Marshal(v interface{}, opts ...ProtOption) ([]byte, error) {
	opt := defaultOpt()
	for _, o := range opts {
		o.apply(&opt)
	}

	fields, err := getEncodeFields(v)
	if err != nil {
		return nil, err
	}

	arr := convertSlice(v)

	return marshal(arr, fields, opt), nil
}

func marshal(arr []reflect.Value, fields []field, opt protOptions) []byte {
	e := new(encodeState)

	rs, cs := opt.rowSep, opt.columnSep
	for i, row := range arr {
		encodeRow(e, row, fields, cs)
		if i < len(arr)-1 {
			e.out.WriteByte(rs)
		}
	}
	return e.out.Bytes()
}

func encodeRow(e *encodeState, v reflect.Value, fs []field, sep byte) {
	for i := range fs {
		fieldEncode(e, v.FieldByName(fs[i].name))
		if i < len(fs)-1 {
			e.out.WriteByte(sep)
		}
	}
}

func convertSlice(v interface{}) []reflect.Value {
	a := reflect.ValueOf(v)
	l := a.Len()
	s := make([]reflect.Value, l)
	for i := 0; i < l; i++ {
		s[i] = a.Index(i)
	}
	return s
}

// encodeValid
func encodeValid(v interface{}) error {
	va := reflect.ValueOf(v)
	if va.Kind() != reflect.Slice || va.IsNil() {
		return errors.New("必须是slice类型")
	}

	if va.Type().Elem().Kind() != reflect.Struct {
		return errors.New("slice元素必须为struct类型")
	}
	return nil
}

func fieldTypeValid(v reflect.Value) bool {
	kind := v.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	default:
		return false
	}
}

func fieldEncode(e *encodeState, v reflect.Value) {
	kind := v.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intEncode(e, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uIntEncode(e, v)
	case reflect.Float32, reflect.Float64:
		floatEncode(e, v)
	case reflect.String:
		stringEncode(e, v)
	}
}

func intEncode(e *encodeState, v reflect.Value) {
	var buf []byte
	e.out.Write(strconv.AppendInt(buf, v.Int(), 10))
}
func uIntEncode(e *encodeState, v reflect.Value) {
	var buf []byte
	e.out.Write(strconv.AppendUint(buf, v.Uint(), 10))
}
func floatEncode(e *encodeState, v reflect.Value) {
	var buf []byte
	e.out.Write(strconv.AppendFloat(buf, v.Float(), 'f', 10, 64))
}
func stringEncode(e *encodeState, v reflect.Value) {
	e.out.Write([]byte(v.String()))
}
