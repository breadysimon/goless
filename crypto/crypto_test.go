package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		text   string
		encErr bool
		decErr bool
	}{
		{
			"enc-dec",
			"secret0123456789",
			"username\npassword",
			false, false,
		},
		{
			"enc-dec with 32byte key",
			"secret0123456789secret0123456789",
			"username\npassword",
			false, false,
		},
		{
			"return err if key is not 16byte",
			"secret012345678",
			"",
			true, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := Encrypt(tt.key, tt.text)
			assert.Equal(t, tt.encErr, (err != nil))
			if err == nil {
				t.Log("encoded:", enc)
				dec, err := Decrypt(tt.key, enc)

				assert.Nil(t, err)
				assert.Equal(t, tt.text, dec)
			}

		})
	}
}
