package ecpay

const (
	// 環境設定
	Production Environment = "https://einvoice.ecpay.com.tw"
	Stage      Environment = "https://einvoice-stage.ecpay.com.tw"
	
	// 載具類別
	CarrierTypeNone    = ""  // 無載具
	CarrierTypeMember  = "2" // 會員載具  
	CarrierTypeCitizen = "3" // 手機條碼
	
	// 列印旗標
	PrintNo  = "0" // 不列印
	PrintYes = "1" // 列印
	
	// 捐贈旗標
	DonationNo  = "0" // 不捐贈
	DonationYes = "1" // 捐贈
	
	// 課稅類別
	TaxTypeRegular = "1" // 應稅
	TaxTypeZero    = "2" // 零稅率
	TaxTypeFree    = "3" // 免稅
	TaxTypeSpecial = "4" // 應稅(特種稅率)
	TaxTypeMixed   = "9" // 混合應稅與免稅
	
	// 字軌類別
	InvTypeGeneral = "07" // 一般稅額
	InvTypeSpecial = "08" // 特種稅額
	
	// VAT 設定
	VatYes = "1" // 商品單價含稅
	VatNo  = "0" // 商品單價未稅
	
	// 發票狀態
	InvoiceStatusNormal  = "1" // 正常
	InvoiceStatusInvalid = "0" // 作廢
	
	// 上傳狀態
	UploadStatusYes = "1" // 已上傳
	UploadStatusNo  = "0" // 未上傳
)