package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"hash"
	"hash/adler32"

	"golang.org/x/crypto/ripemd160"
)

func HashHmac(algo string, data string, key string) string {
	var m hash.Hash
	dataByte := []byte(data)
	keyByte := []byte(key)
	switch algo {
	case "sha256":
		m = hmac.New(sha256.New, keyByte)
	case "ripemd160":
		m = hmac.New(ripemd160.New, keyByte)
	default:
		panic("未定义加密类型")
	}
	_, _ = m.Write(dataByte)
	result := m.Sum(nil)
	return hex.EncodeToString(result)
}

func Hash(algo string, data string) string {
	var m hash.Hash
	dataByte := []byte(data)
	switch algo {
	case "adler32":
		m = adler32.New()
	}
	_, _ = m.Write(dataByte)
	result := m.Sum(nil)
	return hex.EncodeToString(result)
}

// Md5 字符串md5加密
func Md5(str string) string {
	hashs := md5.New()
	hashs.Write([]byte(str))
	return hex.EncodeToString(hashs.Sum(nil))
}

// Base64Encode base64编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// JsonEncode 用于对变量进行 JSON 编码，该函数如果执行成功返回 JSON 数据
func JsonEncode(arg any) string {
	jsonData, err := json.Marshal(arg)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

// JsonDecode JSON 格式的字符串进行解码，并转换为 map[string]any
func JsonDecode(jsonStr string) map[string]any {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil
	}
	return result
}
