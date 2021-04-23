package prot

import (
	"bytes"
	"testing"
)

type Data struct {
	String string
	User   User
}

func TestUnMarshal(t *testing.T) {
	ss, users := initData()
	bs := write(ss)
	var result []User
	err := UnMarshal(bs[:len(bs)-1], &result)
	if err != nil {
		t.Error(err)
	}
	for i := range users {
		if result[i].UserName != users[i].UserName || result[i].Age != users[i].Age ||
			result[i].Money1 != users[i].Money1 || result[i].Money2 != users[i].Money2 ||
			result[i].UID != users[i].UID || result[i].Man != users[i].Man {
			t.Error(ss[i])
		}
	}
}

func BenchmarkUnMarshal(b *testing.B) {
	ss, _ := initData()
	bs := write(ss)
	b.ResetTimer()
	var u1 []User
	for i := 0; i < b.N; i++ {
		UnMarshal(bs[:len(bs)-1], &u1)
	}
}

func initData() ([]string, []User) {
	return []string{
			"wc,18,3.123,3.123456789,true,987654321",
			"3321,100,3,0.001,false,123456789",
			"0,0,0,0,false,0",
		}, []User{
			{
				UserName: "wc", Age: 18, Money1: 3.123, Money2: 3.123456789, Man: true, UID: 987654321,
			},
			{
				UserName: "3321", Age: 100, Money1: 3, Money2: 0.001, Man: false, UID: 123456789,
			},
			{
				UserName: "0", Age: 0, Money1: 0, Money2: 0, Man: false, UID: 0,
			},
		}
}

func write(ss []string) []byte {
	var bb bytes.Buffer
	for _, s := range ss {
		bb.Write([]byte(s))
		bb.WriteByte('\n')
	}
	bs := bb.Bytes()
	return bs[:len(bs)-1]
}
