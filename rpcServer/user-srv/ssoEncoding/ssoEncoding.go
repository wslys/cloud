package ssoEncoding

import (
	"encoding/base64"
)

var (
	encodeX64  = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")
	decodeX64  = make([]byte, 256)
	base64Code = base64.URLEncoding
)

//将base64字符串转换为uint32
//base64	是输入base64字符串
//i32	是输出无符号32位int数
//err	是输出错误信息
func DecodeBase64ToUint32(base64 string) (i32 uint32, err error) {
	var codeBytes []byte
	codeBytes, err = base64Code.DecodeString(base64)
	if err != nil {
		return
	}
	L := len(codeBytes)
	switch L {
	case 0:
		return 0, nil
	case 1:
		return uint32(codeBytes[0]), nil
	case 2:
		return uint32(codeBytes[0]) | uint32(codeBytes[1])<<8, nil
	case 3:
		return uint32(codeBytes[0]) | uint32(codeBytes[1])<<8 |
			uint32(codeBytes[2])<<16, nil
	case 4:
		return uint32(codeBytes[0]) | uint32(codeBytes[1])<<8 |
			uint32(codeBytes[2])<<16 | uint32(codeBytes[3])<<24, nil
	default:
		return uint32(codeBytes[L-4]) | uint32(codeBytes[L-3])<<8 |
			uint32(codeBytes[L-2])<<16 | uint32(codeBytes[L-1])<<24, nil
	}
}

//将base64字符串转换为uin64
//base64	是输入64进制字符串
//i32	是输出无符号32位int数
//err	是输出错误信息
func DecodeBase64ToUint64(base64 string) (i64 uint64, err error) {
	var codeBytes []byte
	codeBytes, err = base64Code.DecodeString(base64)
	if err != nil {
		return
	}
	L := len(codeBytes)
	switch L {
	case 0:
		return 0, nil
	case 1:
		return uint64(codeBytes[0]), nil
	case 2:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8, nil
	case 3:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8 |
			uint64(codeBytes[2])<<16, nil
	case 4:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8 |
			uint64(codeBytes[2])<<16 | uint64(codeBytes[3])<<24, nil
	case 5:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8 |
			uint64(codeBytes[2])<<16 | uint64(codeBytes[3])<<24 |
			uint64(codeBytes[4])<<32, nil
	case 6:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8 |
			uint64(codeBytes[2])<<16 | uint64(codeBytes[3])<<24 |
			uint64(codeBytes[4])<<32 | uint64(codeBytes[5])<<40, nil
	case 7:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8 |
			uint64(codeBytes[2])<<16 | uint64(codeBytes[3])<<24 |
			uint64(codeBytes[4])<<32 | uint64(codeBytes[5])<<40 |
			uint64(codeBytes[6])<<48, nil
	case 8:
		return uint64(codeBytes[0]) | uint64(codeBytes[1])<<8 |
			uint64(codeBytes[2])<<16 | uint64(codeBytes[3])<<24 |
			uint64(codeBytes[4])<<32 | uint64(codeBytes[5])<<40 |
			uint64(codeBytes[6])<<48 | uint64(codeBytes[7])<<56, nil
	default:
		return uint64(codeBytes[L-8]) | uint64(codeBytes[L-7])<<8 |
			uint64(codeBytes[L-6])<<16 | uint64(codeBytes[L-5])<<24 |
			uint64(codeBytes[L-4])<<32 | uint64(codeBytes[L-3])<<40 |
			uint64(codeBytes[L-2])<<48 | uint64(codeBytes[L-1])<<56, nil
	}
	return
}

//将uint32转换为base64进制字符串
//i32	是输入无符号32位int数
//fill	是输入随机填充数
func EncodeUint32ToBase64(i32 uint32, fill []byte) (s string) {
	if fill == nil {
		codeBytes := make([]byte, 4)
		codeBytes[0], codeBytes[1] = byte(i32), byte(i32>>8)
		codeBytes[2], codeBytes[3] = byte(i32>>16), byte(i32>>24)
		return base64Code.EncodeToString(codeBytes)
	} else {
		codeBytes := make([]byte, len(fill), len(fill)+4)
		copy(codeBytes, fill)
		codeBytes = append(codeBytes, byte(i32), byte(i32>>8),
			byte(i32>>16), byte(i32>>24))
		return base64Code.EncodeToString(codeBytes)
	}

}

//将uin642转换为base64进制字符串
//i64	是输入无符号64位int数
//fill	是输入随机填充数
func EncodeUint64ToBase64(i64 uint64, fill []byte) (s string) {
	if fill == nil {
		codeBytes := []byte{byte(i64), byte(i64 >> 8),
			byte(i64 >> 16), byte(i64 >> 24), byte(i64 >> 32), byte(i64 >> 40),
			byte(i64 >> 48), byte(i64 >> 56)}
		return base64Code.EncodeToString(codeBytes)
	} else {
		codeBytes := make([]byte, len(fill), len(fill)+8)
		copy(codeBytes, fill)
		codeBytes = append(codeBytes, byte(i64), byte(i64>>8),
			byte(i64>>16), byte(i64>>24), byte(i64>>32), byte(i64>>40),
			byte(i64>>48), byte(i64>>56))
		return base64Code.EncodeToString(codeBytes)
	}

}

func initDecodeX64() {
	for i, _ := range decodeX64 {
		decodeX64[i] = 0
	}
	for i, b := range encodeX64 {
		decodeX64[b] = byte(i)
	}
}

func init() {
	initDecodeX64()
}
