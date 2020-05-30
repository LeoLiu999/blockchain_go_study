package main

import (
	"bytes"
	"math/big"
)

//base58 比base64少了 0(零)O(大写O)I（大写I）l(小写L) + /
var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

//Base58Encode base58加密
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(base58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}
	//翻转bytes 获取base58的取模结果
	ReverseBytes(result)

	//将版本放入base58加密的结果中 input的第一个字节 = 0x00
	for b := range input {
		if b == 0x00 {
			result = append([]byte{base58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	return result
}

//Base58Decode  decodes Base58-encoded data
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(base58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}
