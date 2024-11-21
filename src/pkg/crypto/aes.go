// https://www.cnblogs.com/Nihility/p/14695647.html
// https://www.cnblogs.com/remixnameless/p/15894694.html
// https://wisp888.github.io/Golang%E5%8A%A0%E5%AF%86%E6%96%B9%E6%B3%95%E5%A4%A7%E5%85%A8.html

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	// keyStr is a constant string used as the encryption key.
	keyStr = "Npf4zWUvqDp6LmQtNxkorgn1qSAgSMGW"
)

var key = []byte(keyStr)

// =================== CFB ======================

// aesEncryptCFB encrypts the original data using AES in Cipher Feedback (CFB) mode.
func aesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}

// aesDecryptCFB decrypts the encrypted data using AES in CFB mode.
// TODO 携程中解密失败（panic）会导致主进程奔溃
func aesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

// Encryption encrypts a plain text string.
func Encryption(plain string) string {
	plaByt := []byte(plain)
	encrypted := aesEncryptCFB(plaByt, key)
	_cipher := hex.EncodeToString(encrypted)
	return _cipher
}

// Decryption decrypts a cipher text string.
func Decryption(cipherText string) (string, error) {
	encryptedHex, err := hex.DecodeString(cipherText)
	if err != nil {
		hlog.Error("Unable to decode hexadecimal string: ", err)
		return "", err
	}
	plain := aesDecryptCFB(encryptedHex, key)
	return string(plain), nil
}

// func Test() {
// 	data := []byte("hello word") // 待加密的数据
// 	fmt.Println("------------------ CFB模式 --------------------")
// 	encrypted := AesEncryptCFB(data, key)
// 	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
// 	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
// 	decrypted := AesDecryptCFB(encrypted, key)
// 	fmt.Println("解密结果：", string(decrypted))

// 	fmt.Println("解密结果：", Decryption("4e5ed44794c42720a8656da96f5f275b5cd7c34b6a72"))
// 	fmt.Println("加密结果：", Encryption("123123"))
// }
