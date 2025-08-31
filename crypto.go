package ecpay

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net/url"
)

// CryptoHandler AES 加解密處理器
type CryptoHandler struct {
	key   []byte
	iv    []byte
	debug bool
}

// NewCryptoHandler 建立新的加解密處理器
func NewCryptoHandler(hashKey, hashIV string) *CryptoHandler {
	return &CryptoHandler{
		key: []byte(hashKey),
		iv:  []byte(hashIV),
	}
}

// SetDebug 設定除錯模式
func (ch *CryptoHandler) SetDebug(debug bool) {
	ch.debug = debug
}

// Encrypt AES-128-CBC 加密 (PKCS7 Padding)
func (ch *CryptoHandler) Encrypt(plainText string) (string, error) {
	// Step 1: URL Encode 原始資料
	urlEncoded := url.QueryEscape(plainText)
	
	// Step 2: 建立 AES cipher
	block, err := aes.NewCipher(ch.key)
	if err != nil {
		return "", fmt.Errorf("建立 AES cipher 失敗: %v", err)
	}
	
	// Step 3: PKCS7 Padding
	plainBytes := []byte(urlEncoded)
	plainBytes = ch.pkcs7Padding(plainBytes, block.BlockSize())
	
	// Step 4: CBC 模式加密
	cipherText := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, ch.iv)
	mode.CryptBlocks(cipherText, plainBytes)
	
	// Step 5: Base64 編碼
	result := base64.StdEncoding.EncodeToString(cipherText)
	
	if ch.debug {
		fmt.Printf("=== AES 加密 ===\n")
		fmt.Printf("原始資料長度: %d\n", len(plainText))
		fmt.Printf("URL Encode 後: %s\n", urlEncoded)
		fmt.Printf("加密後 (Base64): %s\n", result)
	}
	
	return result, nil
}

// Decrypt AES-128-CBC 解密
func (ch *CryptoHandler) Decrypt(encryptedText string) (string, error) {
	// Step 1: Base64 解碼
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("Base64 解碼失敗: %v", err)
	}
	
	// Step 2: 建立 AES cipher
	block, err := aes.NewCipher(ch.key)
	if err != nil {
		return "", fmt.Errorf("建立 AES cipher 失敗: %v", err)
	}
	
	// Step 3: CBC 模式解密
	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, ch.iv)
	mode.CryptBlocks(plainText, cipherText)
	
	// Step 4: 移除 PKCS7 Padding
	plainText = ch.pkcs7UnPadding(plainText)
	
	// Step 5: URL Decode
	result, err := url.QueryUnescape(string(plainText))
	if err != nil {
		return "", fmt.Errorf("URL Decode 失敗: %v", err)
	}
	
	if ch.debug {
		fmt.Printf("=== AES 解密 ===\n")
		fmt.Printf("解密後資料長度: %d\n", len(result))
	}
	
	return result, nil
}

// pkcs7Padding PKCS7 填充
func (ch *CryptoHandler) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 移除 PKCS7 填充
func (ch *CryptoHandler) pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	unPadding := int(data[length-1])
	if unPadding > length {
		return data
	}
	return data[:(length - unPadding)]
}