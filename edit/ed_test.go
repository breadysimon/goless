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
		src     []EditStep
		in, out string
	}{
		{
			"generic",
			[]EditStep{
				{"S", "aaa", "bbb"},
				{"D", "ccc", ""},
			},
			"xaaavv\nbbbb\nccc1234\ndddd\n1234cccc45674",
			"xaaavv\nbbbb\ndddd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := "/tmp/aaa"
			ioutil.WriteFile(f, []byte(tt.in), os.ModePerm)
			out, err := EditFile(f, tt.src)
			assert.Nil(t, err)
			assert.Equal(t, tt.out, out)

			fout, _ := ioutil.ReadFile(f)
			assert.Equal(t, tt.out, string(fout))
		})
	}
}
