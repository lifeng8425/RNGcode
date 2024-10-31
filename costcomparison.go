package main

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"time"

	drbg "github.com/canonical/go-sp800.90a-drbg"
)

var (
	Byt  = int(16)
	echo = 1e6
)

func main() {
	UK := RS(16)
	D := RS(4)
	Cnt := RS(4)
	ot := time.Now()
	for i := 0; i < 1e5; i++ {
		h := sha256.New()
		h.Write(getAnd3(UK, D, Cnt))
		h.Sum(nil)
	}
	Ot := time.Since(ot)
	fmt.Println(Ot)

	ot = time.Now()
	for i := 0; i < 1e5; i++ {
		h := sha256.New()
		h.Write(getOxr(UK, getAnd(D, Cnt)))
		h.Sum(nil)
	}
	Ot = time.Since(ot)
	fmt.Println(Ot)
	hashtime(int(echo), 2)
	hashtime(int(echo), 4)
	hashtime(int(echo), 6)
	hashtime(int(echo), 8)
	hashtime(int(echo), 16)
	hashtime(int(echo), 24)
	hashtime(int(echo), 32)
	hashtime(int(echo), 64)
	hashtime(int(echo), 64*2)
	hashtime(int(echo), 64*4)
	hashAndtime(int(echo), 4)
	hashOxrtime(int(echo), 4)
	hashAndtime(int(echo), 4*2)
	hashOxrtime(int(echo), 4*2)
	hashAndtime(int(echo), 4*3)
	hashOxrtime(int(echo), 4*3)
	hashAndtime(int(echo), 4*4)
	hashOxrtime(int(echo), 4*4)

	HASH_DRBGTime(int(echo))
	HASH_DRBGWithExternalEntropyTime(int(echo))
	CTR_DRBGTime(int(echo))
	CTR_DRBGWithExternalEntropyTime(int(echo))
	HMAC_DRBGTime(int(echo))
	HMAC_DRBGWithExternalEntropyTime(int(echo))
}
func HASH_DRBGTime(echo int) {
	entropySource := rand.Reader
	hashFunc := crypto.SHA1
	personalization := []byte("personalization string")
	drbg, _ := drbg.NewHash(hashFunc, personalization, entropySource)
	data := make([]byte, 16)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		drbg.Generate(nil, data)
	}
	Ot := time.Since(ot)
	fmt.Println("HASH:", Ot)
}
func HASH_DRBGWithExternalEntropyTime(echo int) {

	entropyInput := []byte("initial entropy input")
	nonce := []byte("nonce value")
	hashFunc := crypto.SHA1

	personalization := []byte("personalization string")

	entropySource := rand.Reader

	drbg, _ := drbg.NewHashWithExternalEntropy(hashFunc, entropyInput, nonce, personalization, entropySource)

	data := make([]byte, 16)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		drbg.Generate(nil, data)
	}
	Ot := time.Since(ot)
	fmt.Println("HASHWith:", Ot)

}
func HMAC_DRBGWithExternalEntropyTime(echo int) {
	entropyInput := []byte("initial entropy input")
	nonce := []byte("nonce value")
	hashFunc := crypto.SHA1
	personalization := []byte("personalization string")

	entropySource := rand.Reader

	drbg, _ := drbg.NewHMACWithExternalEntropy(hashFunc, entropyInput, nonce, personalization, entropySource)
	data := make([]byte, 16)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		drbg.Generate(nil, data)
	}
	Ot := time.Since(ot)
	fmt.Println("HMACWith:", Ot)
}
func HMAC_DRBGTime(echo int) {
	entropySource := rand.Reader
	hashFunc := crypto.SHA1
	personalization := []byte("personalization string")
	drbg, _ := drbg.NewHMAC(hashFunc, personalization, entropySource)

	data := make([]byte, 16)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		drbg.Generate(nil, data)
	}
	Ot := time.Since(ot)
	fmt.Println("HMAC:", Ot)
}

func CTR_DRBGTime(echo int) {

	entropySource := rand.Reader
	personalization := []byte("personalization string")

	drbg, _ := drbg.NewCTR(aes.BlockSize, personalization, entropySource)

	data := make([]byte, 16)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		drbg.Generate(nil, data)
	}
	Ot := time.Since(ot)
	fmt.Println("CTR:", Ot)

}
func CTR_DRBGWithExternalEntropyTime(echo int) {

	entropyInput := []byte("initial entropy input")
	nonce := []byte("nonce value")

	personalization := []byte("personalization string")

	entropySource := rand.Reader

	drbg, _ := drbg.NewCTRWithExternalEntropy(aes.BlockSize, entropyInput, nonce, personalization, entropySource)
	data := make([]byte, 16)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		drbg.Generate(nil, data)
	}
	Ot := time.Since(ot)
	fmt.Println("CTRWith:", Ot)

}
func hashtime(echo int, Byt int) {
	b := make([]byte, Byt)
	rand.Read(b)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		h := sha1.New()
		h.Write(b)
		h.Sum(nil)
	}
	Ot := time.Since(ot)
	fmt.Println(Byt, "bytes,hash:", Ot)
}
func hashOxrtime(echo int, Byt int) {
	a := make([]byte, Byt)
	aa := make([]byte, Byt)
	aaa := make([]byte, Byt)

	ot := time.Now()
	for i := 0; i < echo; i++ {
		h := sha1.New()
		h.Write(getOxr3(a, aa, aaa))
		h.Sum(nil)
	}
	Ot := time.Since(ot)
	fmt.Println(Byt, "bytes,hashoxr:", Ot)
}
func hashAndtime(echo int, Byt int) {
	a := make([]byte, Byt)
	aa := make([]byte, Byt)
	aaa := make([]byte, Byt)
	ot := time.Now()
	for i := 0; i < echo; i++ {
		h := sha1.New()
		h.Write(getAnd3(a, aa, aaa))
		h.Sum(nil)
	}
	Ot := time.Since(ot)
	fmt.Println(Byt, "bytes,hashand:", Ot)
}

func RS(Byt int) []byte {
	r := make([]byte, Byt)
	rand.Read(r)
	return r
}

func get16(a []byte) []byte {
	c := make([]byte, 16)
	for i := 0; i < 16; i++ {
		c[i] = a[i]
	}
	return c
}
func getAnd(a []byte, b []byte) []byte {
	la := len(a)
	lb := len(b)
	c := make([]byte, la+lb)
	for i := 0; i < la; i++ {
		c[i] = a[i]
	}
	for i := la; i < la+lb; i++ {
		c[i] = b[i-la]
	}
	return c
}
func getAnd3(a []byte, b []byte, c []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	cc := make([]byte, la+lb+lc)
	i := 0
	for ; i < la; i++ {
		cc[i] = a[i]
	}
	for ; i < la+lb; i++ {
		cc[i] = b[i-la]
	}
	for ; i < la+lb+lc; i++ {
		cc[i] = c[i-la-lb]
	}
	return cc
}
func getAnd4(a []byte, b []byte, c []byte, d []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	ld := len(d)
	cc := make([]byte, la+lb+lc+ld)
	i := 0
	for ; i < la; i++ {
		cc[i] = a[i]
	}
	for ; i < la+lb; i++ {
		cc[i] = b[i-la]
	}
	for ; i < la+lb+lc; i++ {
		cc[i] = c[i-la-lb]
	}
	for ; i < la+lb+lc+ld; i++ {
		cc[i] = c[i-la-lb-lc]
	}
	return cc
}
func getAnd5(a []byte, b []byte, c []byte, d []byte, e []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	ld := len(d)
	le := len(e)
	cc := make([]byte, la+lb+lc+ld+le)
	i := 0
	for ; i < la; i++ {
		cc[i] = a[i]
	}
	for ; i < la+lb; i++ {
		cc[i] = b[i-la]
	}
	for ; i < la+lb+lc; i++ {
		cc[i] = c[i-la-lb]
	}
	for ; i < la+lb+lc+ld; i++ {
		cc[i] = c[i-la-lb-lc]
	}
	for ; i < la+lb+lc+ld+le; i++ {
		cc[i] = c[i-la-lb-lc-ld]
	}
	return cc
}
func getOxr(a []byte, b []byte) []byte {
	la := len(a)
	lb := len(b)
	if la > lb {
		c := make([]byte, la)
		i := 0
		t := 0
		for ; i < la; i++ {
			if i%lb == 0 {
				t = 0
			}
			c[i] = a[i] ^ b[t]
			t++
		}
		return c
	}
	c := make([]byte, lb)
	i := 0
	t := 0
	for ; i < lb; i++ {
		if i%la == 0 {
			t = 0
		}
		c[i] = b[i] ^ a[t]
		t++
	}
	return c
}
func getOxr3(a []byte, b []byte, c []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	if la >= lb && la >= lc {
		cc := make([]byte, la)
		i := 0
		t := 0
		tt := 0
		for ; i < la; i++ {
			if i%lb == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			cc[i] = a[i] ^ b[t] ^ c[tt]
			t++
			tt++
		}
		return cc
	}
	if lb >= la && lb >= lc {
		cc := make([]byte, lb)
		i := 0
		t := 0
		tt := 0
		for ; i < lb; i++ {
			if i%la == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			cc[i] = b[i] ^ a[t] ^ c[tt]
			t++
			tt++
		}
		return cc
	}
	cc := make([]byte, lc)
	i := 0
	t := 0
	tt := 0
	for ; i < lc; i++ {
		if i%la == 0 {
			t = 0
		}
		if i%lb == 0 {
			tt = 0
		}
		cc[i] = c[i] ^ a[t] ^ b[tt]
		t++
		tt++
	}
	return cc
}
func getOxr4(a []byte, b []byte, c []byte, d []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	ld := len(d)
	if la >= lb && la >= lc && la >= ld {
		cc := make([]byte, la)
		i := 0
		t := 0
		tt := 0
		ttt := 0
		for ; i < la; i++ {
			if i%lb == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			if i%ld == 0 {
				ttt = 0
			}
			cc[i] = a[i] ^ b[t] ^ c[tt] ^ d[ttt]
			t++
			tt++
			ttt++
		}
		return cc
	}
	if lb >= la && lb >= lc && lb >= ld {
		cc := make([]byte, lb)
		i := 0
		t := 0
		tt := 0
		ttt := 0
		for ; i < lb; i++ {
			if i%la == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			if i%ld == 0 {
				ttt = 0
			}
			cc[i] = b[i] ^ a[t] ^ c[tt] ^ d[ttt]
			t++
			tt++
			ttt++
		}
		return cc
	}
	if ld >= la && ld >= lc && ld >= lb {
		cc := make([]byte, ld)
		i := 0
		t := 0
		tt := 0
		ttt := 0
		for ; i < ld; i++ {
			if i%la == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			if i%lb == 0 {
				ttt = 0
			}
			cc[i] = d[i] ^ a[t] ^ c[tt] ^ b[ttt]
			t++
			tt++
			ttt++
		}
		return cc
	}
	cc := make([]byte, lc)
	i := 0
	t := 0
	tt := 0
	ttt := 0
	for ; i < lc; i++ {
		if i%la == 0 {
			t = 0
		}
		if i%lb == 0 {
			tt = 0
		}
		if i%ld == 0 {
			ttt = 0
		}
		cc[i] = c[i] ^ a[t] ^ b[tt] ^ d[ttt]
		t++
		tt++
		ttt++
	}
	return cc
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func Hash(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	return get16(h.Sum(nil))
}

// AES CTR
func AESEncryptCTR(PlainText, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("err")
	}
	iv := []byte("12345678asdfghjk")
	stream := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(PlainText))
	stream.XORKeyStream(cipherText, PlainText)
	return cipherText
}
func AESDecryptCTR(cipherText, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("err")
	}
	iv := []byte("12345678asdfghjk")
	stream := cipher.NewCTR(block, iv)

	stream.XORKeyStream(cipherText, cipherText)
	return cipherText
}
