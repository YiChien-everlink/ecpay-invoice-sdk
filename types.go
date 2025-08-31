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
	MerchantID string `json:"MerchantID"`
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
	
	// 商品明細（原始）
	Items []Item `json:"-"`
	
	// 商品明細（處理後，用|分隔）
	ItemName    string `json:"ItemName"`
	ItemCount   string `json:"ItemCount"`
	ItemWord    string `json:"ItemWord"`
	ItemPrice   string `json:"ItemPrice"`
	ItemTaxType string `json:"ItemTaxType"`
	ItemAmount  string `json:"ItemAmount"`
	ItemRemark  string `json:"ItemRemark,omitempty"`
	
	// 其他選填
	DelayDay    int    `json:"DelayDay,omitempty"`
	ECBankID    string `json:"ECBankID,omitempty"`
}

// Item 商品明細
type Item struct {
	ItemSeq     string `json:"ItemSeq,omitempty"`
	ItemName    string `json:"ItemName"`
	ItemCount   string `json:"ItemCount"`
	ItemWord    string `json:"ItemWord"`
	ItemPrice   string `json:"ItemPrice"`
	ItemTaxType string `json:"ItemTaxType"`
	ItemAmount  string `json:"ItemAmount"`
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
		amount, err := strconv.Atoi(item.ItemAmount)
		if err != nil {
			return NewError(ErrCodeValidation, fmt.Sprintf("商品金額格式錯誤: %s", item.ItemAmount))
		}
		totalAmount += amount
	}
	
	salesAmount, err := strconv.Atoi(r.SalesAmount)
	if err != nil {
		return NewError(ErrCodeValidation, "發票金額格式錯誤")
	}
	
	// 檢查金額是否一致（考慮稅額）
	if r.TaxType == "1" { // 應稅
		expectedTotal := int(float64(totalAmount) * 1.05)
		if expectedTotal != salesAmount {
			return NewError(ErrCodeValidation, 
				fmt.Sprintf("發票金額不一致: 預期 %d, 實際 %d", expectedTotal, salesAmount))
		}
	} else {
		if totalAmount != salesAmount {
			return NewError(ErrCodeValidation, 
				fmt.Sprintf("發票金額不一致: 預期 %d, 實際 %d", totalAmount, salesAmount))
		}
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

// QueryInvoiceRequest 查詢發票請求
type QueryInvoiceRequest struct {
	RelateNumber string `json:"RelateNumber"`
}

// Validate 驗證查詢請求
func (r *QueryInvoiceRequest) Validate() error {
	if r.RelateNumber == "" {
		return NewError(ErrCodeValidation, "RelateNumber 不能為空")
	}
	
	return nil
}

// QueryInvoiceResponse 查詢發票回應
type QueryInvoiceResponse struct {
	RtnCode         int    `json:"RtnCode"`
	RtnMsg          string `json:"RtnMsg"`
	InvoiceNo       string `json:"IIS_Number"`
	InvoiceDate     string `json:"IIS_Create_Date"`
	InvoiceStatus   string `json:"IIS_Invoice_Status"`
	SalesAmount     string `json:"IIS_Sales_Amount"`
	TaxAmount       string `json:"IIS_Tax_Amount"`
	TotalAmount     string `json:"IIS_Amount"`
	RemainAmount    string `json:"IIS_Remain_Amount"`
	InvalidStatus   string `json:"IIS_Invalid_Status"`
	UploadStatus    string `json:"IIS_Upload_Status"`
	TurnkeyStatus   string `json:"IIS_Turnkey_Status"`
	PrintFlag       string `json:"IIS_Print_Flag"`
	AwardFlag       string `json:"IIS_Award_Flag"`
	AwardType       string `json:"IIS_Award_Type"`
	RandomNumber    string `json:"InvoiceRandomNumber"`
	CarrierType     string `json:"IIS_Carrier_Type"`
	CarrierNum      string `json:"IIS_Carrier_Num"`
	CustomerName    string `json:"IIS_Customer_Name"`
	CustomerID      string `json:"IIS_Customer_ID"`
}

// AllowanceInvoiceRequest 開立折讓請求
type AllowanceInvoiceRequest struct {
	InvoiceNo       string        `json:"InvoiceNo"`
	InvoiceDate     string        `json:"InvoiceDate"`
	AllowanceNotify string        `json:"AllowanceNotify"`
	CustomerName    string        `json:"CustomerName"`
	NotifyMail      string        `json:"NotifyMail,omitempty"`
	NotifyPhone     string        `json:"NotifyPhone,omitempty"`
	AllowanceAmount string        `json:"AllowanceAmount"`
	
	// 商品明細（原始）
	Items []AllowanceItem `json:"-"`
	
	// 商品明細（處理後）
	ItemName    string `json:"ItemName"`
	ItemCount   string `json:"ItemCount"`
	ItemWord    string `json:"ItemWord"`
	ItemPrice   string `json:"ItemPrice"`
	ItemTaxType string `json:"ItemTaxType"`
	ItemAmount  string `json:"ItemAmount"`
}

// AllowanceItem 折讓商品明細
type AllowanceItem struct {
	ItemSeq     string `json:"ItemSeq,omitempty"`
	ItemName    string `json:"ItemName"`
	ItemCount   string `json:"ItemCount"`
	ItemWord    string `json:"ItemWord"`
	ItemPrice   string `json:"ItemPrice"`
	ItemTaxType string `json:"ItemTaxType"`
	ItemAmount  string `json:"ItemAmount"`
}

// Validate 驗證折讓請求
func (r *AllowanceInvoiceRequest) Validate() error {
	if r.InvoiceNo == "" {
		return NewError(ErrCodeValidation, "InvoiceNo 不能為空")
	}
	
	if r.InvoiceDate == "" {
		return NewError(ErrCodeValidation, "InvoiceDate 不能為空")
	}
	
	if r.AllowanceAmount == "" {
		return NewError(ErrCodeValidation, "AllowanceAmount 不能為空")
	}
	
	if len(r.Items) == 0 {
		return NewError(ErrCodeValidation, "折讓商品明細不能為空")
	}
	
	return nil
}

// AllowanceInvoiceResponse 開立折讓回應
type AllowanceInvoiceResponse struct {
	RtnCode     int    `json:"RtnCode"`
	RtnMsg      string `json:"RtnMsg"`
	AllowanceNo string `json:"IA_Allow_No"`
	InvoiceNo   string `json:"IA_Invoice_No"`
	AllowanceDate string `json:"IA_Date"`
	AllowanceAmount string `json:"IA_Amount"`
	RemainAmount string `json:"IA_Remain_Amount"`
}

// AllowanceInvalidRequest 作廢折讓請求
type AllowanceInvalidRequest struct {
	InvoiceNo   string `json:"InvoiceNo"`
	AllowanceNo string `json:"AllowanceNo"`
	Reason      string `json:"Reason"`
}

// Validate 驗證作廢折讓請求
func (r *AllowanceInvalidRequest) Validate() error {
	if r.InvoiceNo == "" {
		return NewError(ErrCodeValidation, "InvoiceNo 不能為空")
	}
	
	if r.AllowanceNo == "" {
		return NewError(ErrCodeValidation, "AllowanceNo 不能為空")
	}
	
	if r.Reason == "" {
		return NewError(ErrCodeValidation, "Reason 不能為空")
	}
	
	return nil
}

// AllowanceInvalidResponse 作廢折讓回應
type AllowanceInvalidResponse struct {
	RtnCode int    `json:"RtnCode"`
	RtnMsg  string `json:"RtnMsg"`
}