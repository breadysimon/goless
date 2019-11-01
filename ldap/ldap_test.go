package ldap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCN(t *testing.T) {
	tests := []struct {
		name string
		dn   string
		want string
	}{
		{
			"chenji",
			"CN=Simon,OU=信息技术部,OU=XXX,OU=CCCC,OU=ZZZ,DC=aaa,DC=com",
			"Simon",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetCN(tt.dn))
		})
	}
}
