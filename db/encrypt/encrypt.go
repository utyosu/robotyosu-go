package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"github.com/utyosu/robotyosu-go/env"
)

var (
	c cipher.Block
)

func init() {
	var err error
	c, err = aes.NewCipher([]byte(env.EncryptKey))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(env.EncryptKey), err)
		os.Exit(-1)
	}
}

func EncryptString(s string) []byte {
	cfb := cipher.NewCFBEncrypter(c, []byte(env.EncryptCommonIV))
	ciphertext := make([]byte, len(s))
	cfb.XORKeyStream(ciphertext, []byte(s))
	return ciphertext
}

func DecryptString(s []byte) string {
	cfbdec := cipher.NewCFBDecrypter(c, []byte(env.EncryptCommonIV))
	plaintextCopy := make([]byte, len(s))
	cfbdec.XORKeyStream(plaintextCopy, s)
	return string(plaintextCopy)
}
