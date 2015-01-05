package lib

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

const (
	PBKDF2_ITERATION = 1000
	SALT_BYTE_SIZE   = 20
	HASH_BYTE_SIZE   = 40
)

func GenerateSalt(size ...int) string {
	c := 10
	if len(size) > 0 {
		c = size[0]
	}
	b := make([]byte, c)
	_, err := rand.Read(b)
	_ = err
	return hex.EncodeToString(b)
}

// use pbkdf2 encode password
func EncodePbkdf2(password string, salt string) string {
	pwd := PBKDF2([]byte(password), []byte(salt), PBKDF2_ITERATION, HASH_BYTE_SIZE, sha256.New)
	return hex.EncodeToString(pwd)
}

// Encode string to md5 hex value
func EncodeMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func EncodeHmac(secret, value string, params ...func() hash.Hash) string {
	var h func() hash.Hash
	if len(params) > 0 {
		h = params[0]
	} else {
		h = sha1.New
	}

	hm := hmac.New(h, []byte(secret))
	hm.Write([]byte(value))

	return hex.EncodeToString(hm.Sum(nil))
}

// http://code.google.com/p/go/source/browse/pbkdf2/pbkdf2.go?repo=crypto
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

/***************************************************************/
// hashPassword=salt+hash编码过的连合字符串
func ValidatePassword(password string, hashPassword string) bool {
	// var salt string
	salt := hashPassword[:SALT_BYTE_SIZE]
	// hash = hashPassword[SALT_BYTE_SIZE:]

	return salt+EncodePbkdf2(password, salt) == hashPassword
}