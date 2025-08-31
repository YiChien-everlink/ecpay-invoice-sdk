package ecpay

import (
	"fmt"
	"regexp"
	"strconv"
)

// Environment 環境設定
type Environment string

// BaseRequest 基本請求結構
type BaseRequest struct {
	MerchantID string `json:"MerchantID"`
	RqHeader   struct {
		Timestamp int64  `json:"Timestamp"`
		Revision  string `json:"Revision"`
	} `json:"RqHeader"`
	Data string `json:"Data"` // AES 加密後的資料
}

// BaseResponse 基本回應結構
type BaseResponse struct {
	MerchantID interface{} `json:"MerchantID"` // 可能是 string 或 number
	RpHeader   struct {
		Timestamp int64 `json:"Timestamp"`
	} `json:"RpHeader"`
	TransCode int    `json:"TransCode"`
	TransMsg  string `json:"TransMsg"`
	Data      string `json:"Data"` // AES 加密的回應資料
}

// IssueInvoiceRequest 開立發票請求
type IssueInvoiceRequest struct {
	// 基本資訊
	RelateNumber        string `json:"RelateNumber"`
	CustomerID          string `json:"CustomerID,omitempty"`
	CustomerIdentifier  string `json:"CustomerIdentifier,omitempty"`
	CustomerName        string `json:"CustomerName"`
	CustomerAddr        string `json:"CustomerAddr,omitempty"`
	CustomerPhone       string `json:"CustomerPhone,omitempty"`
	CustomerEmail       string `json:"CustomerEmail"`
	
	// 列印與載具
	Print        string `json:"Print"`
	Donation     string `json:"Donation"`
	LoveCode     string `json:"LoveCode,omitempty"`
	CarrierType  string `json:"CarrierType,omitempty"`
	CarrierNum   string `json:"CarrierNum,omitempty"`
	
	// 稅務資訊
	TaxType      string `json:"TaxType"`
	SalesAmount  string `json:"SalesAmount"`
	InvoiceRemark string `json:"InvoiceRemark,omitempty"`
	InvType      string `json:"InvType"`
	Vat          string `json:"vat,omitempty"`
	
	// 商品明細 - B2C API 使用 Items 陣列
	Items []Item `json:"Items"`
	
	// 其他選填
	DelayDay    int    `json:"DelayDay,omitempty"`
	ECBankID    string `json:"ECBankID,omitempty"`
}

// Item 商品明細
type Item struct {
	ItemSeq     int    `json:"ItemSeq"`
	ItemName    string `json:"ItemName"`
	ItemCount   int    `json:"ItemCount"`
	ItemWord    string `json:"ItemWord"`
	ItemPrice   int    `json:"ItemPrice"`
	ItemTaxType string `json:"ItemTaxType"`
	ItemAmount  int    `json:"ItemAmount"`
	ItemRemark  string `json:"ItemRemark,omitempty"`
}

// Validate 驗證開立發票請求
func (r *IssueInvoiceRequest) Validate() error {
	if r.RelateNumber == "" {
		return NewError(ErrCodeValidation, "RelateNumber 不能為空")
	}
	
	if len(r.RelateNumber) > 30 {
		return NewError(ErrCodeValidation, "RelateNumber 長度不能超過 30")
	}
	
	// 驗證 Email 格式
	if r.CustomerEmail != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(r.CustomerEmail) {
			return NewError(ErrCodeValidation, "Email 格式不正確")
		}
	}
	
	// 驗證手機號碼格式（台灣手機）
	if r.CustomerPhone != "" {
		phoneRegex := regexp.MustCompile(`^09\d{8}$`)
		if !phoneRegex.MatchString(r.CustomerPhone) {
			return NewError(ErrCodeValidation, "手機號碼格式不正確")
		}
	}
	
	// 驗證統一編號（如果有填）
	if r.CustomerIdentifier != "" {
		if !ValidateTaxID(r.CustomerIdentifier) {
			return NewError(ErrCodeValidation, "統一編號格式不正確")
		}
	}
	
	// 驗證捐贈碼
	if r.Donation == "1" && r.LoveCode == "" {
		return NewError(ErrCodeValidation, "選擇捐贈時必須填寫愛心碼")
	}
	
	// 驗證載具
	if r.CarrierType != "" && r.CarrierNum == "" {
		return NewError(ErrCodeValidation, "選擇載具類別時必須填寫載具編號")
	}
	
	// 驗證商品明細
	if len(r.Items) == 0 {
		return NewError(ErrCodeValidation, "商品明細不能為空")
	}
	
	// 驗證金額
	totalAmount := 0
	for _, item := range r.Items {
		totalAmount += item.ItemAmount
	}
	
	salesAmount, err := strconv.Atoi(r.SalesAmount)
	if err != nil {
		return NewError(ErrCodeValidation, "發票金額格式錯誤")
	}
	
	// 檢查金額是否一致
	// B2C API 中 SalesAmount 是未稅金額，稅額另計
	if totalAmount != salesAmount {
		return NewError(ErrCodeValidation, 
			fmt.Sprintf("發票金額不一致: 預期 %d, 實際 %d", totalAmount, salesAmount))
	}
	
	return nil
}

// IssueInvoiceResponse 開立發票回應
type IssueInvoiceResponse struct {
	RtnCode      int    `json:"RtnCode"`
	RtnMsg       string `json:"RtnMsg"`
	InvoiceNo    string `json:"InvoiceNo"`
	InvoiceDate  string `json:"InvoiceDate"`
	RandomNumber string `json:"RandomNumber"`
}

// InvalidInvoiceRequest 作廢發票請求
type InvalidInvoiceRequest struct {
	InvoiceNo   string `json:"InvoiceNo"`
	InvoiceDate string `json:"InvoiceDate"`
	Reason      string `json:"Reason"`
}

// Validate 驗證作廢請求
func (r *InvalidInvoiceRequest) Validate() error {
	if r.InvoiceNo == "" {
		return NewError(ErrCodeValidation, "InvoiceNo 不能為空")
	}
	
	if r.InvoiceDate == "" {
		return NewError(ErrCodeValidation, "InvoiceDate 不能為空")
	}
	
	if r.Reason == "" {
		return NewError(ErrCodeValidation, "Reason 不能為空")
	}
	
	if len(r.Reason) > 20 {
		return NewError(ErrCodeValidation, "Reason 長度不能超過 20")
	}
	
	return nil
}

// InvalidInvoiceResponse 作廢發票回應
type InvalidInvoiceResponse struct {
	RtnCode int    `json:"RtnCode"`
	RtnMsg  string `json:"RtnMsg"`
}