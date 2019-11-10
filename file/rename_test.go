package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRename(t *testing.T) {
	assert.Equal(t, 0, RenamePictures(`/tmp/aaa`))
}
