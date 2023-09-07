package ecb

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

// AES encryption with ECB and PKCS7 padding
func TestAesEcb(t *testing.T)  {
	pt := []byte("Some plain text")
	// Key size for AES is either: 16 bytes (128 bits), 24 bytes (192 bits) or 32 bytes (256 bits)
	key := []byte("secretkey16bytes")

	ct := Encrypt(pt, key)
	fmt.Println("Ciphertext:", fmt.Sprintf("%X", ct))

	recoveredPt := Decrypt(ct, key)
	fmt.Println("Recovered plaintext:", fmt.Sprintf("%s", recoveredPt))

	h := sha1.New()
	h.Write(pt)
	fmt.Println("Sha1:",fmt.Sprintf("%x", h.Sum(nil)))
}
// Output:
// Ciphertext: AF3B0173EAE9DD013A649F4EAABA1376
// Recovered plaintext: Some plain text
// Sha1: d85e382c4a48731d850ec5956a20e5b3ccaa0e7d
//
// SQL
// select hex(aes_encrypt('Some plain text', 'secretkey16bytes'))
// -- AF3B0173EAE9DD013A649F4EAABA1376
// select sha1('Some plain text')
// -- d85e382c4a48731d850ec5956a20e5b3ccaa0e7d
