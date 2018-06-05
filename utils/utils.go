package utils

import "os"

func FileIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func AppendByteSlice(s [] byte, t []byte) []byte {
	zlen := len(s) + len(t)
	z := make([]byte, zlen)
	copy(z, s)
	copy(z[len(s):], t)
	return z
}
