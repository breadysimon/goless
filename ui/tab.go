package ui

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

func TablePrint(fields string, lst interface{}) {
	//table.Output(list)
	// fmt.Println(list)
	list := reflect.ValueOf(lst)
	pading := 2
	w := tabwriter.NewWriter(os.Stdout, 0, 0, pading, ' ', 0)
	fmt.Fprintln(w, strings.ToUpper(fields))
	headers := strings.Split(fields, "\t")
	for n := 0; n < list.Len(); n++ {
		vv := list.Index(n)
		for _, f := range headers {
			if strings.TrimSpace(f) != "" {
				fmt.Fprintf(w, "\t%v", vv.FieldByName(f).Interface())
			}
		}
		fmt.Fprintln(w, "\t")
	}
	w.Flush()

}
