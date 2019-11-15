package reflection

import (
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

func GetId(structPtr interface{}) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(structPtr))
	return v.FieldByName("ID")
}

func SetIdInt(structPtr interface{}, id int) {
	v := reflect.Indirect(reflect.ValueOf(structPtr))
	v.FieldByName("ID").SetInt(int64(id))
}
func SetValue(ptr interface{}, value interface{}) {
	if value != nil && ptr != nil {
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(value))
	}
}
func MakeInstance(ptr interface{}) interface{} {
	tt := reflect.ValueOf(ptr).Elem().Type()
	x := reflect.New(tt)
	return x.Interface()
}

func GetFieldString(structPtr interface{}, name string) string {
	return reflect.ValueOf(structPtr).Elem().FieldByName(name).String()
}

func GetFields(ptr interface{}) map[string]interface{} {
	st := reflect.Indirect(reflect.ValueOf(ptr)).Type()
	sv := reflect.Indirect(reflect.ValueOf(ptr))
	mm := make(map[string]interface{})
	for i := 0; i < sv.NumField(); i++ {
		switch sv.Field(i).Type().String() {
		case "string":
			v := string(sv.Field(i).String())
			if v != "" {
				k := st.Field(i).Name
				mm[k] = v
			}
		case "int":
			v := int(sv.Field(i).Int())
			if v != 0 {
				k := st.Field(i).Name
				mm[k] = v
			}
		}
	}
	return mm
}

// MakeSlice creates a slice of the arguement's struct type,
// then return its ptr.
// x should be a struct instance ptr, e.g. &T{}
func MakeSlice(ptr interface{}) interface{} {
	// find the instance type from ptr
	t := reflect.ValueOf(ptr).Elem().Type()
	// create slice instance of this type
	sl := reflect.New(reflect.SliceOf(t))
	return sl.Interface()
}

func GetSliceItem(t interface{}, i int) interface{} {
	lst := reflect.Indirect(reflect.ValueOf(t))
	if lst.Len() > 0 {
		return lst.Index(i).Interface()
	}
	return nil
}
func GetSliceLen(ptr interface{}) int {
	return reflect.Indirect(reflect.ValueOf(ptr)).Len()
}

func GetSearchableFieldNames(t interface{}) []string {
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
