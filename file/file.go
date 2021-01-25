package file

import (
	"bytes"
	"encoding/csv"
)

// ReadCsv 读取csv文件内容，统一返回[][]string，如果需要特殊格式，自行转换
//
func ReadCsv(file []byte, comma rune) ([][]string, error) {

	//解决部分csv是utf-8 bom编码问题
	if file[0] == 0xef || file[1] == 0xbb || file[2] == 0xbf {
		file = file[3:]
	}
	b := bytes.NewReader(file)
	r := csv.NewReader(b)
	r.Comma = comma

	content, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return content, nil
}
