package ecpay

import (
	"encoding/json"
	"fmt"
)

// IssueInvoice 開立發票
func (c *Client) IssueInvoice(req *IssueInvoiceRequest) (*IssueInvoiceResponse, error) {
	// 驗證請求
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// 設定商品序號
	for i := range req.Items {
		req.Items[i].ItemSeq = i + 1
	}
	
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
	respData, err := c.sendRequest("/B2CInvoice/Invalid", req)
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