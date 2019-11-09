package edit

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditFile(t *testing.T) {
	tests := []struct {
		name    string
		src     [][]string
		yaml    string
		in, out string
	}{
		{
			"generic",
			[][]string{
				{"S", "aaa", "bbb"},
				{"D", "ccc", ""},
			},
			`
- 
 - S # replace 
 - aaa
 - bbb
- 
 - S # replace with group 
 - (vv)
 - ${1}zz 
- 
  - D # delete multiple lines
  - ccc
-
  - D # delete first line
  - ff`,
			"fff\nxaaavv\nbbbb\nccc1234\ndddd\n1234cccc45674",
			"xbbbvvzz\nbbbb\ndddd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := "/tmp/aaa"
			writeToFile(f, tt.in)
			yaml := "/tmp/sss"
			writeToFile(yaml, tt.yaml)

			err := RunScript(yaml, f, f)
			assert.Nil(t, err)

			assert.Equal(t, tt.out, readFromFile(f))
		})
	}
}

func readFromFile(f string) string {
	if out, err := ioutil.ReadFile(f); err != nil {
		panic(err)
	} else {
		return string(out)
	}
}

func writeToFile(f, data string) {
	if err := ioutil.WriteFile(f, []byte(data), os.ModePerm); err != nil {
		panic(err)
	}
}
