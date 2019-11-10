package file

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// md5Sig Create MD5 signature of a file
func Md5Sig(path string, sizeLimit int64) (sum string, err error) {

	finfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if size := finfo.Size(); size > sizeLimit {
		return "", fmt.Errorf("file is too big: %s", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	md5hash := md5.New()

	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	sum = base64.StdEncoding.EncodeToString(md5hash.Sum(nil))
	return sum, nil
}
