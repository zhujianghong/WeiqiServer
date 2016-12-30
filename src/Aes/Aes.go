// Aes
package Aes

import (
	"crypto/aes"
	"errors"
	"log"
)

var AESKEY []byte = []byte("F0ECDA106091DBD598EF4F941F94DDD2")

func GetEncryptAfterLen(length int32) int32 {
	l := (length/aes.BlockSize + 1) * aes.BlockSize
	return l
}

func Encrypt(origData, key []byte) ([]byte, error) {
	Aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := Aes.BlockSize()

	length := (len(origData)/blockSize + 1) * blockSize
	pIn := make([]byte, length)
	pOut := make([]byte, length)
	copy(pIn[0:len(origData)], origData)
	l := length - len(origData)

	for i := 0; i < l; i++ {
		pIn[len(origData)+i] = byte(l)
	}
	log.Println(pIn)
	en_len := 0
	for {
		if en_len < length {
			begin := en_len
			end := en_len + blockSize
			Aes.Encrypt(pOut[begin:end], pIn[begin:end])
			en_len += blockSize
		} else {
			break
		}
	}
	return pOut, nil
}
func Decrypt(crypted, key []byte) ([]byte, error) {
	Aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypted)%Aes.BlockSize() != 0 {
		return nil, errors.New("data aes len error")
	}
	blockSize := Aes.BlockSize()

	length := len(crypted)
	pIn := make([]byte, length)
	pOut := make([]byte, length)
	copy(pIn[0:length], crypted)
	en_len := 0
	for {
		if en_len < length {
			begin := en_len
			end := en_len + blockSize
			Aes.Decrypt(pOut[begin:end], pIn[begin:end])
			en_len += blockSize
		} else {
			break
		}
	}
	l := int(pOut[length-1])
	log.Println(l)
	if l > 16 || l < 1 {
		return nil, errors.New("AES Decrypt ERROR")
	}
	return pOut[0 : length-l], nil
}
