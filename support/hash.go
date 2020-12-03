package support

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)
func Hash(key, s string) string{
	h := md5.New()
	h.Write([]byte(key))
	h.Write([]byte(s))
	v := h.Sum(nil)
	return string(v[:])
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func main() {
	data := []byte("hello")
	fmt.Printf("%x", md5.Sum(data))
}