package conf_center

import (
	"crypto/aes"
	"crypto/cipher"
	"bytes"
	"strconv"
	"strings"
)

func MapMinus(newConfProperties map[string]map[string]string, oldConfProperties map[string]map[string]string) []string {
	oldKeys := Keys(oldConfProperties)
	newKeys := Keys(newConfProperties)
	minus := []string{}
	for _, key := range newKeys {
		if !Contains(oldKeys, key) {
			minus = append(minus, key)
		}
	}
	return minus
}

func Keys(aMap map[string]map[string]string) []string {
	keys := []string{}
	for key := range aMap {
		keys = append(keys, key)
	}
	return keys
}

func Contains(keys []string, k string) bool {
	for _, key := range keys {
		if key == k {
			return true
		}
	}
	return false
}

func AesEncrypt(origData string, key string) (string,error) {
	secretKey := []byte(Padding(key))
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "",err
	}
	originBytes := []byte(Padding(origData))
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, secretKey)
	crypted := make([]byte, len(originBytes))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, originBytes)
	return ByteToHex(crypted),nil
}

func AesDecrypt(ciphertext string, key string) (string,error) {
	secretKey := []byte(Padding(key))
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "",err
	}
	blockMode := cipher.NewCBCDecrypter(block, secretKey)
	cipherBytes := HexToBye(ciphertext)
	origData := make([]byte, len(cipherBytes))
	// origData := crypted
	blockMode.CryptBlocks(origData, cipherBytes)
	origData = UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData),nil
}

//保持一直，自己做padding
func Padding(ciphertext string) string {
	padding := 16 - len(ciphertext) % 16
	padtext := strings.Repeat("\x00", padding)
	return ciphertext + padtext
}

func UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}

func ByteToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

//16进制字符串转[]byte
func HexToBye(hex string) []byte {
	length := len(hex) / 2
	slice := make([]byte, length)
	rs := []rune(hex)

	for i := 0; i < length; i++ {
		s := string(rs[i*2 : i*2+2])
		value, _ := strconv.ParseInt(s, 16, 10)
		slice[i] = byte(value & 0xFF)
	}
	return slice
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}




