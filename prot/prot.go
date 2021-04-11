
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

// getFields 获取有序字段名
func getFields(v interface{}) ([]string, error) {
	//fmt.Println(reflect.TypeOf(v).Kind())
	va := reflect.ValueOf(v)
	if va.Kind() != reflect.Slice {
		return nil, errors.New("必须是slice类型")
	}

	if va.IsNil() {
		return nil, errors.New("不能为nil")
	}
	a := va.Index(0).Interface()
	t := reflect.TypeOf(a)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("slice元素必须为struct类型")
	}
	fieldNum := t.NumField()
	if fieldNum == 0 {
		return nil, errors.New("struct必须有字段")
	}
	fields := make([]string, fieldNum)

	// 将prot的值填充tags
	for i := 0; i < fieldNum; i++ {
		s := t.Field(i).Tag.Get(tagName)
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%s不能转换为int", s)
		}
		if n > fieldNum {
			return nil, errors.New("prot标签必须从1开始，且连续，不重复")
		}
		if len(fields[n-1]) > 0 {
			return nil, errors.New("prot标签必须从1开始，且连续，不重复")
		}
		fields[n-1] = t.Field(i).Name
	}

	//如果tags中有空字符串，则认为prot值有重复
	for _, tag := range fields {
		if len(tag) == 0 {
			return nil, errors.New("prot标签必须从1开始，且连续，不重复")
		}
	}
	return fields, nil
}

func Marshal(v interface{}, opts ...ProtOption) ([]byte, error) {
	opt := defaultOpt()
	for _, o := range opts {
		o.apply(&opt)
	}

	fields, err := getFields(v)
	if err != nil {
		return nil, err
	}

	arr := convertSlice(v)

	return marshal(arr, fields, opt)
}

func marshal(arr []interface{}, fields []string, opt protOptions) ([]byte, error) {
	// rs, cs := opt.rowSep, opt.columnSep
	// for row := range arr {

	// }
	return nil, nil
}

func encodeRow(v reflect.Value, sep byte) []byte {
	buf := new(bytes.Buffer)
	l := v.Len()
	for i := 0; i < l; i++ {
		//buf.Write(v.Bytes())
		buf.WriteByte(sep)
	}
	return buf.Bytes()[:buf.Len()-1]
}

func convertSlice(v interface{}) []interface{} {
	a := reflect.ValueOf(v)
	l := a.Len()
	s := make([]interface{}, l)
	for i := 0; i < l; i++ {
		s[i] = a.Index(i).Interface()
	}
	return s
}

func feildTypeValid(v reflect.Value) bool {
	switch v.Kind() {
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
