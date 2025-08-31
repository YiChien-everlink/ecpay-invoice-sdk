package ecpay

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// GenerateRelateNumber 產生特店自訂編號
func GenerateRelateNumber(prefix string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s%d", prefix, timestamp)
}

// ParseInvoiceDate 解析發票日期
func ParseInvoiceDate(dateStr string) (time.Time, error) {
	// 綠界回傳格式: 2024-01-01 12:00:00
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, dateStr)
}

// FormatInvoiceDate 格式化發票日期
func FormatInvoiceDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// ValidateTaxID 驗證統一編號
func ValidateTaxID(taxID string) bool {
	if len(taxID) != 8 {
		return false
	}
	
	// 統編檢查邏輯
	weights := []int{1, 2, 1, 2, 1, 2, 4, 1}
	sum := 0
	
	for i := 0; i < 8; i++ {
		num := int(taxID[i] - '0')
		if num < 0 || num > 9 {
			return false
		}
		
		product := num * weights[i]
		sum += product / 10
		sum += product % 10
	}
	
	// 特殊情況：第七碼為 7 時的處理
	if taxID[6] == '7' {
		return sum%10 == 0 || (sum+1)%10 == 0
	}
	
	return sum%10 == 0
}

// PrettyPrint 格式化輸出 JSON
func PrettyPrint(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", data)
	}
	return string(b)
}

// CalculateTax 計算稅額
func CalculateTax(amount int, taxType string) (salesAmount, taxAmount int) {
	switch taxType {
	case TaxTypeRegular: // 應稅 5%
		taxAmount = int(float64(amount) * 0.05)
		salesAmount = amount + taxAmount
	case TaxTypeZero: // 零稅率
		taxAmount = 0
		salesAmount = amount
	case TaxTypeFree: // 免稅
		taxAmount = 0
		salesAmount = amount
	default:
		taxAmount = 0
		salesAmount = amount
	}
	
	return
}

// ConvertToInvoiceItem 轉換為發票商品格式
func ConvertToInvoiceItem(name string, count int, price float64) Item {
	amount := float64(count) * price
	
	return Item{
		ItemName:    name,
		ItemCount:   strconv.Itoa(count),
		ItemWord:    "個",
		ItemPrice:   fmt.Sprintf("%.0f", price),
		ItemTaxType: TaxTypeRegular,
		ItemAmount:  fmt.Sprintf("%.0f", amount),
	}
}

// ValidateLoveCode 驗證愛心碼格式
func ValidateLoveCode(code string) bool {
	// 愛心碼為 3-7 碼數字
	if len(code) < 3 || len(code) > 7 {
		return false
	}
	
	for _, c := range code {
		if c < '0' || c > '9' {
			return false
		}
	}
	
	return true
}

// ValidateCarrierNum 驗證載具編號
func ValidateCarrierNum(carrierType, carrierNum string) bool {
	switch carrierType {
	case CarrierTypeCitizen: // 手機條碼
		// 格式: /開頭，後面7個字元
		if len(carrierNum) != 8 || carrierNum[0] != '/' {
			return false
		}
		return true
		
	case CarrierTypeMember: // 會員載具
		// 各商家自訂，通常是會員編號
		return len(carrierNum) > 0 && len(carrierNum) <= 30
		
	default:
		return true
	}
}