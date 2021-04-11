package prot

import (
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	Age      int     `prot:"2"`
	Money1   float32 `prot:"3"`
	Money2   float64 `prot:"4"`
	UserName string  `prot:"1"`
	Man      bool    `prot:"5"`
}

var (
	tab = []User{
		{UserName: "wc", Age: 30, Money1: 3.13},
	}
)

func TestMarshal(t *testing.T) {
	bs, err := Marshal(tab)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bs))
}

func TestJsonMarshal(t *testing.T) {
	bs, err := json.Marshal(tab)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bs))
}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Marshal(tab)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkJsonMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(tab)
		if err != nil {
			fmt.Println(err)
		}
	}
}
