package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/breadysimon/goless/logging"
)

var log *logging.Logger = logging.GetLogger()

// Encrypt encrypt a string and encode it to base64.
func Encrypt(key, text string) (out string, err error) {
	ciphertext, err := encrypt([]byte(key), []byte(text))
	if err == nil {
		out = base64.StdEncoding.EncodeToString(ciphertext)
		return
	}
	log.Error(err)
	return
}

// Decrypt decode base64 string and decrypt it.
func Decrypt(key, text string) (out string, err error) {
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err == nil {
		var data []byte
		data, err = decrypt([]byte(key), ciphertext)
		if err == nil {
			out = string(data)
			return
		}
	}
	log.Error(err)
	return
}

func encrypt(key, text []byte) (out []byte, err error) {
	block, err := aes.NewCipher(key)
	if err == nil {
		// 对IV有随机性要求，但没有保密性要求，所以常见的做法是将IV包含在加密文本当中
		out = make([]byte, aes.BlockSize+len(text))

		// CFB模式，全称Cipher FeedBack模式，译为密文反馈模式
		// 即上一个密文分组作为加密算法的输入，输出与明文异或作为下一个分组的密文。
		// 第一个明文分组，不存在上一个密文分组，因此需要准备与分组等长的初始化向量IV来代替。
		// 常见的做法:将IV包含在加密文本当中
		iv := out[:aes.BlockSize]
		_, err = io.ReadFull(rand.Reader, iv) //ReadFull按buf大小填满由随机生成器产生的数
		if err == nil {
			cfb := cipher.NewCFBEncrypter(block, iv)
			cfb.XORKeyStream(out[aes.BlockSize:], []byte(text))
			return
		}
	}
	return
}

func decrypt(key, text []byte) (out []byte, err error) {
	block, err := aes.NewCipher(key)
	if err == nil {
		if len(text) < aes.BlockSize {
			err = errors.New("ciphertext too short")
		} else {
			iv := text[:aes.BlockSize]
			text = text[aes.BlockSize:]
			cfb := cipher.NewCFBDecrypter(block, iv)
			cfb.XORKeyStream(text, text)
			out = text
			if err == nil {
				return
			}
		}
	}
	return
}

func EncryptUserPassword(key, u, p string) (secret string, err error) {
	s := fmt.Sprintf("%s\n%s", u, p)
	secret, err = Encrypt(key, s)
	return
}
func DecryptUserPassword(key, secret string, u, p *string) (err error) {
	var s string
	s, err = Decrypt(key, secret)
	if err != nil {
		return
	}

	ss := strings.Split(s, "\n")
	if len(ss) != 2 {
		return fmt.Errorf("failed to parse the secret in config.")
	}
	*u = ss[0]
	*p = ss[1]
	return
}
