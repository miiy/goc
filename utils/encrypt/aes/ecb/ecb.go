package ecb

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func Encrypt(pt, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

func Decrypt(ct, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		fmt.Println(key)
		panic(err.Error())
	}
	return pt
}

func EncryptToHex(pt, key []byte) string {
	return fmt.Sprintf("%X", Encrypt(pt, key))
}

func DecryptFromHex(pt string, key []byte) string {
	b, err := hex.DecodeString(pt)
	if err != nil {
		return ""
	}
	return string(Decrypt(b, key))
}

func Sha1ToHex(pt []byte) string {
	h := sha1.New()
	h.Write(pt)
	return fmt.Sprintf("%x", h.Sum(nil))
}
