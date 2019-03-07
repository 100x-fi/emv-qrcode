# EMVco QR code

### Documents
https://github.com/dongri/emvco-qrcode-doc

### MPM (Merchant Presented Mode)
```go
package main

import(
	"log"

	"github.com/dongri/emvco-qrcode/emvco/mpm"
)
func main() {
	emvqr := new(mpm.EMVQR)
	emvqr.PayloadFormatIndicator = "01"
	emvqr.PointOfInitiationMethod = "12" // 11 is static qrcode
	merchantAccountInformation := new(mpm.MerchantAccountInformation)
	merchantAccountInformation.Tag = "15"
	merchantAccountInformation.Value = "ABCDEF1234567890"
	emvqr.MerchantAccountInformation = *merchantAccountInformation
	emvqr.MerchantCategoryCode = "5311"
	emvqr.TransactionCurrency = "392"
	emvqr.TransactionAmount = 999
	emvqr.CountryCode = "JP"
	emvqr.MerchantName = "DONGRI"
	emvqr.MerchantCity = "TOKYO"
	additionalTemplate := new(mpm.AdditionalDataFieldTemplate)
	additionalTemplate.BillNumber = "hoge"
	additionalTemplate.ReferenceLabel = "fuga"
	additionalTemplate.TerminalLabel = "piyo"
	emvqr.AdditionalDataFieldTemplate = *additionalTemplate
	mpmQRCode, err := emvqr.GeneratePayload()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(mpmQRCode)
	// 0002010102121516ABCDEF123456789052045311530339254039995802JP5906DONGRI6005TOKYO62240104hoge0504fuga0704piyo63043FE6
}
```

### CPM (Consumer Presented Mode)
```go
package main

import(
	"log"

	"github.com/dongri/emvco-qrcode/emvco/cpm"
)
func main() {
	qr := new(cpm.EMVQR)
	qr.DataPayloadFormatIndicator = "CPV01"

	appTemplate1 := new(cpm.ApplicationTemplate)
	appTemplate1.DataApplicationDefinitionFileName = "A0000000555555"
	appTemplate1.DataApplicationLabel = "Product1"
	qr.ApplicationTemplates = append(qr.ApplicationTemplates, *appTemplate1)

	appTemplate2 := new(cpm.ApplicationTemplate)
	appTemplate2.DataApplicationDefinitionFileName = "A0000000666666"
	appTemplate2.DataApplicationLabel = "Product2"
	qr.ApplicationTemplates = append(qr.ApplicationTemplates, *appTemplate2)

	cdt := new(cpm.CommonDataTemplate)
	cdt.DataApplicationPAN = "1234567890123458"
	cdt.DataCardholderName = "CARDHOLDER/EMV"
	cdt.DataLanguagePreference = "ruesdeen"

	cdtt := new(cpm.CommonDataTransparentTemplate)
	cdtt.DataIssuerApplicationData = "06010A03000000"
	cdtt.DataApplicationCryptogram = "584FD385FA234BCC"
	cdtt.DataApplicationTransactionCounter = "0001"
	cdtt.DataUnpredictableNumber = "6D58EF13"
	cdt.CommonDataTransparentTemplates = append(cdt.CommonDataTransparentTemplates, *cdtt)

	qr.CommonDataTemplates = append(qr.CommonDataTemplates, *cdt)

	comQRCode, err := qr.GeneratePayload()
	if err != nil {
		log.Println(err)
	}
	log.Println(comQRCode)
	// hQVDUFYwMWETTwegAAAAVVVVUAhQcm9kdWN0MWETTwegAAAAZmZmUAhQcm9kdWN0MmJJWggSNFZ4kBI0WF8gDkNBUkRIT0xERVIvRU1WXy0IcnVlc2RlZW5kIZ8QBwYBCgMAAACfJghYT9OF+iNLzJ82AgABnzcEbVjvEw==
}
```
