package prot

import (
	"fmt"
	"testing"
)

type User struct {
	UserName string  `prot:"1"`
	Age      uint8   `prot:"2"`
	Money1   float32 `prot:"3"`
	Money2   float64 `prot:"4"`
	Man      bool    `prot:"5"`
	UID      uint64  `prot:"6"`
}

func TestMarshal(t *testing.T) {
	ss, users := initData()
	bs, err := Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	bs1 := write(ss)
	if string(bs1) != string(bs) {
		t.Error(string(bs))
	}
}

func BenchmarkMarshal(b *testing.B) {
	_, users := initData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(users)
		if err != nil {
			fmt.Println(err)
		}
	}
}
