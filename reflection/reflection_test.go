package reflection

import (
	"fmt"
	"reflect"
	"testing"
)

type aaa struct {
	F1 string
	F2 int
	f3 int
	f4 string
}

func Test_guessResourceFields(t *testing.T) {
	a := aaa{F1: "asdf", F2: 12, f4: "xxx"}
	m := GetFields(&a)
	fmt.Println(m)
}

func Test(t *testing.T) {
	lst := []aaa{{F1: "fff1"}}
	x := GetSliceItem(lst, 0)
	fmt.Println(x)
}

func TestSetValue(t *testing.T) {
	a := aaa{F2: 1}
	b := aaa{F2: 2}
	var x, y interface{}
	x = &a
	y = &b
	reflect.ValueOf(x).Elem().Set(reflect.ValueOf(y).Elem())
	// fmt.Println(p)
	fmt.Println(reflect.Indirect(reflect.ValueOf(x)), y)
}
