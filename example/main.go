package main

import (
	"fmt"
	"log"
	"time"
	
	ecpay "github.com/YiChien-everlink/ecpay-invoice-sdk"
)

func main() {
	// 建立客戶端 (使用測試環境)
	client := ecpay.NewClient(
		"2000132",           // 測試商店代號
		"ejCk326UnaZWKisg",  // 測試 HashKey (16碼)
		"q9jcZX8Ib9LM8wYk",  // 測試 HashIV (16碼)
		ecpay.Stage,         // 測試環境
	)
	
	// 啟用除錯模式
	client.SetDebug(true)
	
	// 模擬從前端收到的訂單編號（隨機生成）
	orderID := ecpay.GenerateRelateNumber("ORD")  // 實際應用中，這會從前端或資料庫取得
	
	// 開立發票
	fmt.Println("========== 開立發票 ==========")
	invoice, err := issueInvoice(client, orderID)
	if err != nil {
		log.Fatal("開立發票失敗:", err)
	}
	
	fmt.Printf("發票開立成功！\n")
	fmt.Printf("訂單編號: %s\n", orderID)
	fmt.Printf("發票號碼: %s\n", invoice.InvoiceNo)
	fmt.Printf("發票日期: %s\n", invoice.InvoiceDate)
	fmt.Printf("隨機碼: %s\n", invoice.RandomNumber)
	
	// 等待一下再查詢
	time.Sleep(2 * time.Second)
	
	// 作廢發票（選擇性）
	fmt.Println("\n========== 作廢發票 ==========")
	if err := invalidInvoice(client, invoice); err != nil {
		log.Printf("作廢發票失敗: %v\n", err)
	} else {
		fmt.Println("發票作廢成功！")
	}
}

// issueInvoice 開立發票
// orderID: 前端傳來的訂單編號（作為 RelateNumber）
func issueInvoice(client *ecpay.Client, orderID string) (*ecpay.IssueInvoiceResponse, error) {
	// 實際應用中，這些資料應該從資料庫或前端取得
	req := &ecpay.IssueInvoiceRequest{
		RelateNumber:  orderID,  // 使用前端提供的訂單編號
		CustomerID:    "",        // 客戶代號（選填）
		CustomerName:  "測試客戶",
		CustomerAddr:  "台北市信義區信義路五段7號",
		CustomerPhone: "0912345678",
		CustomerEmail: "test@example.com",
		Print:         ecpay.PrintNo,
		Donation:      ecpay.DonationNo,
		LoveCode:      "",
		CarrierType:   "",
		CarrierNum:    "",
		TaxType:       ecpay.TaxTypeRegular,
		SalesAmount:   "100",
		InvoiceRemark: "測試發票備註",
		InvType:       ecpay.InvTypeGeneral,
		Vat:           ecpay.VatYes,
		Items: []ecpay.Item{
			{
				ItemName:    "測試商品A",
				ItemCount:   1,
				ItemWord:    "個",
				ItemPrice:   50,
				ItemTaxType: ecpay.TaxTypeRegular,
				ItemAmount:  50,
			},
			{
				ItemName:    "測試商品B",
				ItemCount:   1,
				ItemWord:    "個",
				ItemPrice:   50,
				ItemTaxType: ecpay.TaxTypeRegular,
				ItemAmount:  50,
			},
		},
	}
	
	return client.IssueInvoice(req)
}

func invalidInvoice(client *ecpay.Client, invoice *ecpay.IssueInvoiceResponse) error {
	req := &ecpay.InvalidInvoiceRequest{
		InvoiceNo:   invoice.InvoiceNo,
		InvoiceDate: invoice.InvoiceDate,
		Reason:      "測試作廢",
	}
	
	_, err := client.InvalidInvoice(req)
	return err
}