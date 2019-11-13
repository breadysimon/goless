package rest

import (
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

func SetId(m interface{}, id int) {
	v := reflect.Indirect(reflect.ValueOf(m))
	v.FieldByName("ID").SetInt(int64(id))
}
func makeItem(t interface{}) interface{} {
	tt := reflect.ValueOf(t).Elem().Type()
	x := reflect.New(tt)
	return x.Interface()
}
func readFieldString(t interface{}, name string) string {
	return reflect.ValueOf(t).Elem().FieldByName(name).String()
}

// makeSlice creates a slice of the arguement's struct type,
// then return its ptr.
// x should be a struct instance ptr, e.g. &T{}
func makeSlice(x interface{}) interface{} {
	// find the instance type from ptr
	t := reflect.ValueOf(x).Elem().Type()
	// create slice instance of this type
	sl := reflect.New(reflect.SliceOf(t))
	return sl.Interface()
}

func getSearchableFields(t interface{}) []string {
	fields := []string{}
	st := reflect.ValueOf(t).Elem().Type()
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		restTag := f.Tag.Get("rest")
		if strings.Contains(restTag, "search") {
			fields = append(fields, gorm.ToColumnName(f.Name))
		}
	}
	return fields
}
