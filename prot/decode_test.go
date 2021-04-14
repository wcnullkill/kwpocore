package prot

import (
	"fmt"
	"testing"
)

func TestUnMarshal(t *testing.T) {
	s := `wc,30,3.13,3.1314
	wcc,32,3.13,3.1314
	wccc,33,3.13,3.1314
	wcccc,34,3.13,3.1314
	wccccc,35,3.13,3.1314`
	var u1 []User
	if err := UnMarshal([]byte(s), &u1); err != nil {
		fmt.Println(err)
	}
	fmt.Println(u1)
}
