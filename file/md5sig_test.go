package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SIZE_LIMIT = 10 << 30 //10G
)

func TestMd5(t *testing.T) {

	path := filepath.Join(os.TempDir(), "md5testxxxmd5testxxx")
	ioutil.WriteFile(path, []byte("12345"), os.ModePerm)

	sum, err := Md5Sig(path, SIZE_LIMIT)
	assert.Nil(t, err)
	assert.Equal(t, "gnzLDuqKcGxMNKFokfhOew==", sum)
}
