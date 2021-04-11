package prot

import (
	"fmt"
	"testing"
)

var (
	s = `wc,30,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314
	wcc,32,3.13,3.1314`
)

func TestUnMarshal(t *testing.T) {
	var u1 []User
	if err := UnMarshal([]byte(s), &u1); err != nil {
		fmt.Println(err)
	}
	for _, u := range u1 {
		fmt.Println(u)
	}
}

func BenchmarkUnMarshal(b *testing.B) {
	var u1 []User
	for i := 0; i < b.N; i++ {
		UnMarshal([]byte(s), &u1)
	}
}
