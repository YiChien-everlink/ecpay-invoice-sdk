package ecpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 綠界電子發票客戶端
type Client struct {
	MerchantID string
	HashKey    string
	HashIV     string
	Env        Environment
	httpClient *http.Client
	debug      bool
	crypto     *CryptoHandler
}

// NewClient 建立新的客戶端
func NewClient(merchantID, hashKey, hashIV string, env Environment) *Client {
	return &Client{
		MerchantID: merchantID,
		HashKey:    hashKey,
		HashIV:     hashIV,
		Env:        env,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		debug:  false,
		crypto: NewCryptoHandler(hashKey, hashIV),
	}
}

// SetDebug 設定除錯模式
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
	c.crypto.SetDebug(debug)
}

// SetTimeout 設定逾時時間
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// sendRequest 發送 API 請求
func (c *Client) sendRequest(apiPath string, data interface{}) ([]byte, error) {
	// 將資料轉換為 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, NewError(ErrCodeRequest, fmt.Sprintf("JSON 編碼失敗: %v", err))
	}
	
	if c.debug {
		fmt.Printf("=== 原始請求資料 ===\n%s\n", string(jsonData))
	}
	
	// AES 加密
	encryptedData, err := c.crypto.Encrypt(string(jsonData))
	if err != nil {
		return nil, NewError(ErrCodeCrypto, fmt.Sprintf("加密失敗: %v", err))
	}
	
	// 建立請求物件
	request := BaseRequest{
		MerchantID: c.MerchantID,
		Data:       encryptedData,
	}
	request.RqHeader.Timestamp = time.Now().Unix()
	request.RqHeader.Revision = "3.0.0"
	
	// 轉換為 JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, NewError(ErrCodeRequest, fmt.Sprintf("建立請求失敗: %v", err))
	}
	
	// 建立 HTTP 請求
	fullURL := fmt.Sprintf("%s%s", c.Env, apiPath)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, NewError(ErrCodeRequest, fmt.Sprintf("建立 HTTP 請求失敗: %v", err))
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	if c.debug {
		fmt.Printf("=== API 請求 ===\n")
		fmt.Printf("URL: %s\n", fullURL)
		fmt.Printf("Request Body: %s\n", string(requestBody))
	}
	
	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewError(ErrCodeNetwork, fmt.Sprintf("發送請求失敗: %v", err))
	}
	defer resp.Body.Close()
	
	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewError(ErrCodeResponse, fmt.Sprintf("讀取回應失敗: %v", err))
	}
	
	if c.debug {
		fmt.Printf("=== API 回應 ===\n")
		fmt.Printf("Response Body: %s\n", string(body))
	}
	
	// 解析基本回應
	var baseResp BaseResponse
	if err := json.Unmarshal(body, &baseResp); err != nil {
		return nil, NewError(ErrCodeParse, fmt.Sprintf("解析回應失敗: %v", err))
	}
	
	// 檢查回應狀態
	if baseResp.TransCode != 1 {
		return nil, NewError(ErrCodeAPI, fmt.Sprintf("%s (Code: %d)", baseResp.TransMsg, baseResp.TransCode))
	}
	
	// 解密回應資料
	decryptedData, err := c.crypto.Decrypt(baseResp.Data)
	if err != nil {
		return nil, NewError(ErrCodeCrypto, fmt.Sprintf("解密回應失敗: %v", err))
	}
	
	if c.debug {
		fmt.Printf("=== 解密後回應 ===\n%s\n", decryptedData)
	}
	
	return []byte(decryptedData), nil
}