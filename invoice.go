package ecpay

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// IssueInvoice 開立發票
func (c *Client) IssueInvoice(req *IssueInvoiceRequest) (*IssueInvoiceResponse, error) {
	// 驗證請求
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// 轉換商品明細格式
	req.processItems()
	
	// 發送請求
	respData, err := c.sendRequest("/B2CInvoice/Issue", req)
	if err != nil {
		return nil, err
	}
	
	// 解析回應
	var resp IssueInvoiceResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, NewError(ErrCodeParse, fmt.Sprintf("解析發票回應失敗: %v", err))
	}
	
	// 檢查業務邏輯錯誤
	if resp.RtnCode != 1 {
		return nil, NewError(ErrCodeAPI, resp.RtnMsg)
	}
	
	return &resp, nil
}

// InvalidInvoice 作廢發票
func (c *Client) InvalidInvoice(req *InvalidInvoiceRequest) (*InvalidInvoiceResponse, error) {
	// 驗證請求
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// 發送請求
	respData, err := c.sendRequest("/Invoice/IssueInvalid", req)
	if err != nil {
		return nil, err
	}
	
	// 解析回應
	var resp InvalidInvoiceResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, NewError(ErrCodeParse, fmt.Sprintf("解析作廢回應失敗: %v", err))
	}
	
	// 檢查業務邏輯錯誤
	if resp.RtnCode != 1 {
		return nil, NewError(ErrCodeAPI, resp.RtnMsg)
	}
	
	return &resp, nil
}

// QueryInvoice 查詢發票
func (c *Client) QueryInvoice(req *QueryInvoiceRequest) (*QueryInvoiceResponse, error) {
	// 驗證請求
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// 發送請求
	respData, err := c.sendRequest("/Query/Issue", req)
	if err != nil {
		return nil, err
	}
	
	// 解析回應
	var resp QueryInvoiceResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, NewError(ErrCodeParse, fmt.Sprintf("解析查詢回應失敗: %v", err))
	}
	
	// 檢查業務邏輯錯誤
	if resp.RtnCode != 1 {
		return nil, NewError(ErrCodeAPI, resp.RtnMsg)
	}
	
	return &resp, nil
}

// AllowanceInvoice 開立折讓
func (c *Client) AllowanceInvoice(req *AllowanceInvoiceRequest) (*AllowanceInvoiceResponse, error) {
	// 驗證請求
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// 轉換商品明細格式
	req.processItems()
	
	// 發送請求
	respData, err := c.sendRequest("/Invoice/Allowance", req)
	if err != nil {
		return nil, err
	}
	
	// 解析回應
	var resp AllowanceInvoiceResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, NewError(ErrCodeParse, fmt.Sprintf("解析折讓回應失敗: %v", err))
	}
	
	// 檢查業務邏輯錯誤
	if resp.RtnCode != 1 {
		return nil, NewError(ErrCodeAPI, resp.RtnMsg)
	}
	
	return &resp, nil
}

// AllowanceInvalidInvoice 作廢折讓
func (c *Client) AllowanceInvalidInvoice(req *AllowanceInvalidRequest) (*AllowanceInvalidResponse, error) {
	// 驗證請求
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// 發送請求
	respData, err := c.sendRequest("/Invoice/AllowanceInvalid", req)
	if err != nil {
		return nil, err
	}
	
	// 解析回應
	var resp AllowanceInvalidResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, NewError(ErrCodeParse, fmt.Sprintf("解析作廢折讓回應失敗: %v", err))
	}
	
	// 檢查業務邏輯錯誤
	if resp.RtnCode != 1 {
		return nil, NewError(ErrCodeAPI, resp.RtnMsg)
	}
	
	return &resp, nil
}

// processItems 處理商品明細（內部方法）
func (r *IssueInvoiceRequest) processItems() {
	itemNames := []string{}
	itemCounts := []string{}
	itemWords := []string{}
	itemPrices := []string{}
	itemTaxTypes := []string{}
	itemAmounts := []string{}
	itemRemarks := []string{}
	
	for i, item := range r.Items {
		item.ItemSeq = strconv.Itoa(i + 1)
		itemNames = append(itemNames, item.ItemName)
		itemCounts = append(itemCounts, item.ItemCount)
		itemWords = append(itemWords, item.ItemWord)
		itemPrices = append(itemPrices, item.ItemPrice)
		itemTaxTypes = append(itemTaxTypes, item.ItemTaxType)
		itemAmounts = append(itemAmounts, item.ItemAmount)
		itemRemarks = append(itemRemarks, item.ItemRemark)
	}
	
	// 綠界要求使用 | 分隔商品明細
	r.ItemName = strings.Join(itemNames, "|")
	r.ItemCount = strings.Join(itemCounts, "|")
	r.ItemWord = strings.Join(itemWords, "|")
	r.ItemPrice = strings.Join(itemPrices, "|")
	r.ItemTaxType = strings.Join(itemTaxTypes, "|")
	r.ItemAmount = strings.Join(itemAmounts, "|")
	r.ItemRemark = strings.Join(itemRemarks, "|")
}

// processItems 處理折讓商品明細（內部方法）
func (r *AllowanceInvoiceRequest) processItems() {
	itemNames := []string{}
	itemCounts := []string{}
	itemWords := []string{}
	itemPrices := []string{}
	itemTaxTypes := []string{}
	itemAmounts := []string{}
	
	for i, item := range r.Items {
		item.ItemSeq = strconv.Itoa(i + 1)
		itemNames = append(itemNames, item.ItemName)
		itemCounts = append(itemCounts, item.ItemCount)
		itemWords = append(itemWords, item.ItemWord)
		itemPrices = append(itemPrices, item.ItemPrice)
		itemTaxTypes = append(itemTaxTypes, item.ItemTaxType)
		itemAmounts = append(itemAmounts, item.ItemAmount)
	}
	
	// 綠界要求使用 | 分隔商品明細
	r.ItemName = strings.Join(itemNames, "|")
	r.ItemCount = strings.Join(itemCounts, "|")
	r.ItemWord = strings.Join(itemWords, "|")
	r.ItemPrice = strings.Join(itemPrices, "|")
	r.ItemTaxType = strings.Join(itemTaxTypes, "|")
	r.ItemAmount = strings.Join(itemAmounts, "|")
}