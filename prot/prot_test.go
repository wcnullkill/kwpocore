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
}

var (
	tab = []User{
		{UserName: "wc", Age: 30, Money1: 28977988998789.093409203410234, Money2: 342394808.0980349582098540980349850239845023984},
		{UserName: "wcc", Age: 32, Money1: 28977988998789.093409203410234, Money2: 342394808.0980349582098540980349850239845023984},
		{UserName: "wccc", Age: 33, Money1: 28977988998789.093409203410234, Money2: 342394808.0980349582098540980349850239845023984},
		{UserName: "wcccc", Age: 34, Money1: 28977988998789.093409203410234, Money2: 342394808.0980349582098540980349850239845023984},
		{UserName: "wccccc", Age: 35, Money1: 28977988998789.093409203410234, Money2: 342394808.0980349582098540980349850239845023984},
	}
)

func TestMarshal(t *testing.T) {
	bs, err := Marshal(tab)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(bs))
}

func TestJsonMarshal(t *testing.T) {
	bs, err := json.Marshal(tab)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(bs))
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
