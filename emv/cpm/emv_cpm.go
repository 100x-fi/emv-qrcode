package cpm

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// EMVQR ...
type EMVQR struct {
	DataPayloadFormatIndicator string                // 85
	ApplicationTemplates       []ApplicationTemplate // 61
	CommonDataTemplates        []CommonDataTemplate  // 62
	customBERTLVMapID4         map[string]*CustomBERTLVField
	customBERTLVMapID2         map[string]*CustomBERTLVField
}

// CustomBERTLVField ...
type CustomBERTLVField struct {
	Description string
	IsHex       bool
}

// ApplicationTemplate ...
type ApplicationTemplate struct {
	BERTLV
	ApplicationSpecificTransparentTemplates []ApplicationSpecificTransparentTemplate // 63
}

// CommonDataTemplate ...
type CommonDataTemplate struct {
	BERTLV
	CommonDataTransparentTemplates []CommonDataTransparentTemplate // 64
}

// CommonDataTransparentTemplate ...
type CommonDataTransparentTemplate struct {
	BERTLV
}

// ApplicationSpecificTransparentTemplate ...
type ApplicationSpecificTransparentTemplate struct {
	BERTLV
}

// BERTLV ...
type BERTLV struct {
	DataApplicationDefinitionFileName string // "4F"
	DataApplicationLabel              string // "50"
	DataTrack2EquivalentData          string // "57"
	DataApplicationPAN                string // "5A"
	DataCardholderName                string // "5F20"
	DataLanguagePreference            string // "5F2D"
	DataIssuerURL                     string // "5F50"
	DataApplicationVersionNumber      string // "9F08"
	DataIssuerApplicationData         string // "9F10"
	DataTokenRequestorID              string // "9F19"
	DataPaymentAccountReference       string // "9F24"
	DataLast4DigitsOfPAN              string // "9F25"
	DataApplicationCryptogram         string // "9F26"
	DataApplicationTransactionCounter string // "9F36"
	DataUnpredictableNumber           string // "9F37"
	AdditionalDataMap                 map[string]string
}

// AddAdditionalData ...
func (b *BERTLV) AddAdditionalData(id, val string) {
	if b.AdditionalDataMap == nil {
		b.AdditionalDataMap = make(map[string]string)
	}
	b.AdditionalDataMap[id] = val
}

// AddCustomBERTLVID4 ...
func (c *EMVQR) AddCustomBERTLVID4(id, description string, isHex bool) {
	if c.customBERTLVMapID4 == nil {
		c.customBERTLVMapID4 = make(map[string]*CustomBERTLVField)
	}
	c.customBERTLVMapID4[id] = &CustomBERTLVField{
		Description: description,
		IsHex:       isHex,
	}
}

// AddCustomBERTLVID2 ...
func (c *EMVQR) AddCustomBERTLVID2(id, description string, isHex bool) {
	if c.customBERTLVMapID2 == nil {
		c.customBERTLVMapID2 = make(map[string]*CustomBERTLVField)
	}
	c.customBERTLVMapID2[id] = &CustomBERTLVField{
		Description: description,
		IsHex:       isHex,
	}
}

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() (string, error) {
	s := ""
	if c.DataPayloadFormatIndicator != "" {
		s += format(IDPayloadFormatIndicator, toHex(c.DataPayloadFormatIndicator))
	} else {
		return "", fmt.Errorf("DataPayloadFormatIndicator is mandatory")
	}
	if len(c.ApplicationTemplates) > 0 {
		for _, t := range c.ApplicationTemplates {
			template := formattingTemplate((t.BERTLV))
			if len(t.ApplicationSpecificTransparentTemplates) > 0 {
				for _, tt := range t.ApplicationSpecificTransparentTemplates {
					ttemplate := formattingTemplate((tt.BERTLV))
					template += format(IDApplicationSpecificTransparentTemplate, ttemplate)
				}
			}
			s += format(IDApplicationTemplate, template)
		}
	}
	if len(c.CommonDataTemplates) > 0 {
		for _, t := range c.CommonDataTemplates {
			template := formattingTemplate(t.BERTLV)
			if len(t.CommonDataTransparentTemplates) > 0 {
				for _, tt := range t.CommonDataTransparentTemplates {
					ttemplate := formattingTemplate(tt.BERTLV)
					template += format(IDCommonDataTransparentTemplate, ttemplate)
				}
			}
			s += format(IDCommonDataTemplate, template)
		}
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString([]byte(string(decoded)))
	return s, nil
}

// Decode ...
func (c *EMVQR) Decode(payload string) (*EMVQR, error) {
	s, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	encoded := hex.EncodeToString(s)
	emvqr := new(EMVQR)
	p := NewParser(encoded)
	idWordCount := IDWordCount
	for p.Next(idWordCount) {
		id := strings.ToUpper(string(p.ID(idWordCount)))
		// length := p.ValueLength(idWordCount)
		hexValue := p.Value(idWordCount)
		switch id {
		case IDPayloadFormatIndicator:
			value, err := fromHex(hexValue)
			if err != nil {
				return nil, err
			}
			emvqr.DataPayloadFormatIndicator = value
		case IDApplicationTemplate:
			applicationTemplate, err := c.ParseApplication(hexValue)
			if err != nil {
				return nil, err
			}
			emvqr.ApplicationTemplates = append(emvqr.ApplicationTemplates, *applicationTemplate)
		case IDCommonDataTemplate:
			commonDataTemplate, err := c.ParseCommonDataTemplate(hexValue)
			if err != nil {
				return nil, err
			}
			emvqr.CommonDataTemplates = append(emvqr.CommonDataTemplates, *commonDataTemplate)
		default:
			// nothing
		}
	}
	return emvqr, nil
}

// ParseApplication ...
func (c *EMVQR) ParseApplication(hexString string) (*ApplicationTemplate, error) {
	applicationTemplate := new(ApplicationTemplate)
	p := NewParser(hexString)
	idWordCount := IDWordCount
	for p.Next(idWordCount) {
		idWordCount = IDWordCount
		id := strings.ToUpper(string(p.ID(idWordCount)))
		// length := p.ValueLength(idWordCount)
		hexVal := p.Value(idWordCount)
		value, err := fromHex(hexVal)
		if err != nil {
			return nil, err
		}
		switch id {
		case TagApplicationDefinitionFileName:
			applicationTemplate.DataApplicationDefinitionFileName = value
		case TagApplicationLabel:
			applicationTemplate.DataApplicationLabel = value
		case TagTrack2EquivalentData:
			applicationTemplate.DataTrack2EquivalentData = hexVal
		case TagApplicationPAN:
			applicationTemplate.DataApplicationPAN = hexVal
		case IDApplicationSpecificTransparentTemplate:
			bertlv, err := c.ParseBERTLV(hexVal)
			if err != nil {
				return nil, err
			}
			applicationTemplate.ApplicationSpecificTransparentTemplates = append(
				applicationTemplate.ApplicationSpecificTransparentTemplates,
				ApplicationSpecificTransparentTemplate{
					*bertlv,
				},
			)
		default:
			if c.customBERTLVMapID2[id] != nil {
				if c.customBERTLVMapID2[id].IsHex {
					applicationTemplate.AddAdditionalData(id, hexVal)
				} else {
					applicationTemplate.AddAdditionalData(id, value)
				}
			} else {
				idWordCount = 4
				// id = strings.ToUpper(string(p.ID(idWordCount)))
				//length = p.ValueLength(idWordCount)
				hexVal = p.Value(idWordCount)
				value, err = fromHex(hexVal)
				if err != nil {
					return nil, err
				}
				switch id {
				case TagCardholderName:
					applicationTemplate.DataCardholderName = value
				case TagLanguagePreference:
					applicationTemplate.DataLanguagePreference = value
				case TagIssuerURL:
					applicationTemplate.DataIssuerURL = value
				case TagApplicationVersionNumber:
					applicationTemplate.DataApplicationVersionNumber = hexVal
				case TagIssuerApplicationData:
					applicationTemplate.DataIssuerApplicationData = hexVal
				case TagTokenRequestorID:
					applicationTemplate.DataTokenRequestorID = hexVal
				case TagPaymentAccountReference:
					applicationTemplate.DataPaymentAccountReference = hexVal
				case TagLast4DigitsOfPAN:
					applicationTemplate.DataLast4DigitsOfPAN = hexVal
				case TagApplicationCryptogram:
					applicationTemplate.DataApplicationCryptogram = hexVal
				case TagApplicationTransactionCounter:
					applicationTemplate.DataApplicationTransactionCounter = hexVal
				case TagUnpredictableNumber:
					applicationTemplate.DataUnpredictableNumber = hexVal
				default:
					if c.customBERTLVMapID4[id] != nil {
						if c.customBERTLVMapID4[id].IsHex {
							applicationTemplate.AddAdditionalData(id, hexVal)
						} else {
							applicationTemplate.AddAdditionalData(id, value)
						}
					}
					// skip unknown
				}
			}
		}
	}
	return applicationTemplate, nil
}

// ParseCommonDataTemplate ...
func (c *EMVQR) ParseCommonDataTemplate(hexString string) (*CommonDataTemplate, error) {
	commonDataTemplate := new(CommonDataTemplate)
	p := NewParser(hexString)
	idWordCount := IDWordCount
	for p.Next(idWordCount) {
		idWordCount = IDWordCount
		id := strings.ToUpper(string(p.ID(idWordCount)))
		// length := p.ValueLength(idWordCount)
		hexVal := p.Value(idWordCount)
		value, err := fromHex(hexVal)
		if err != nil {
			return nil, err
		}
		switch id {
		case TagApplicationDefinitionFileName:
			commonDataTemplate.DataApplicationDefinitionFileName = value
		case TagApplicationLabel:
			commonDataTemplate.DataApplicationLabel = value
		case TagTrack2EquivalentData:
			commonDataTemplate.DataTrack2EquivalentData = hexVal
		case TagApplicationPAN:
			commonDataTemplate.DataApplicationPAN = hexVal
		case IDCommonDataTransparentTemplate:
			bertlv, err := c.ParseBERTLV(hexVal)
			if err != nil {
				return nil, err
			}
			commonDataTemplate.CommonDataTransparentTemplates = append(
				commonDataTemplate.CommonDataTransparentTemplates,
				CommonDataTransparentTemplate{
					*bertlv,
				},
			)
		default:
			if c.customBERTLVMapID2[id] != nil {
				if c.customBERTLVMapID2[id].IsHex {
					commonDataTemplate.AddAdditionalData(id, hexVal)
				} else {
					commonDataTemplate.AddAdditionalData(id, value)
				}
			} else {
				idWordCount = 4
				id = strings.ToUpper(string(p.ID(idWordCount)))
				//length = p.ValueLength(idWordCount)
				hexVal = p.Value(idWordCount)
				value, err = fromHex(hexVal)
				if err != nil {
					return nil, err
				}
				switch id {
				case TagCardholderName:
					commonDataTemplate.DataCardholderName = value
				case TagLanguagePreference:
					commonDataTemplate.DataLanguagePreference = value
				case TagIssuerURL:
					commonDataTemplate.DataIssuerURL = value
				case TagApplicationVersionNumber:
					commonDataTemplate.DataApplicationVersionNumber = hexVal
				case TagIssuerApplicationData:
					commonDataTemplate.DataIssuerApplicationData = hexVal
				case TagTokenRequestorID:
					commonDataTemplate.DataTokenRequestorID = hexVal
				case TagPaymentAccountReference:
					commonDataTemplate.DataPaymentAccountReference = hexVal
				case TagLast4DigitsOfPAN:
					commonDataTemplate.DataLast4DigitsOfPAN = hexVal
				case TagApplicationCryptogram:
					commonDataTemplate.DataApplicationCryptogram = hexVal
				case TagApplicationTransactionCounter:
					commonDataTemplate.DataApplicationTransactionCounter = hexVal
				case TagUnpredictableNumber:
					commonDataTemplate.DataUnpredictableNumber = hexVal
				default:
					if c.customBERTLVMapID4[id] != nil {
						if c.customBERTLVMapID4[id].IsHex {
							commonDataTemplate.AddAdditionalData(id, hexVal)
						} else {
							commonDataTemplate.AddAdditionalData(id, value)
						}
					}
					// skip unknown
				}
			}
		}
	}
	return commonDataTemplate, nil
}

// ParseBERTLV ...
func (c *EMVQR) ParseBERTLV(hexString string) (*BERTLV, error) {
	bertlv := new(BERTLV)
	p := NewParser(hexString)
	idWordCount := IDWordCount
	for p.Next(idWordCount) {
		idWordCount = IDWordCount
		id := strings.ToUpper(string(p.ID(idWordCount)))
		// length := p.ValueLength(idWordCount)
		hexVal := p.Value(idWordCount)
		value, err := fromHex(hexVal)
		if err != nil {
			return nil, err
		}
		switch id {
		case TagApplicationDefinitionFileName:
			bertlv.DataApplicationDefinitionFileName = value
		case TagApplicationLabel:
			bertlv.DataApplicationLabel = value
		case TagTrack2EquivalentData:
			bertlv.DataTrack2EquivalentData = hexVal
		case TagApplicationPAN:
			bertlv.DataApplicationPAN = hexVal
		default:
			if c.customBERTLVMapID2[id] != nil {
				if c.customBERTLVMapID2[id].IsHex {
					bertlv.AddAdditionalData(id, hexVal)
				} else {
					bertlv.AddAdditionalData(id, value)
				}
			} else {
				idWordCount = 4
				id = strings.ToUpper(string(p.ID(idWordCount)))
				//length = p.ValueLength(idWordCount)
				hexVal = p.Value(idWordCount)
				value, err = fromHex(hexVal)
				if err != nil {
					return nil, err
				}
				switch id {
				case TagCardholderName:
					bertlv.DataCardholderName = value
				case TagLanguagePreference:
					bertlv.DataLanguagePreference = value
				case TagIssuerURL:
					bertlv.DataIssuerURL = value
				case TagApplicationVersionNumber:
					bertlv.DataApplicationVersionNumber = hexVal
				case TagIssuerApplicationData:
					bertlv.DataIssuerApplicationData = hexVal
				case TagTokenRequestorID:
					bertlv.DataTokenRequestorID = hexVal
				case TagPaymentAccountReference:
					bertlv.DataPaymentAccountReference = hexVal
				case TagLast4DigitsOfPAN:
					bertlv.DataLast4DigitsOfPAN = hexVal
				case TagApplicationCryptogram:
					bertlv.DataApplicationCryptogram = hexVal
				case TagApplicationTransactionCounter:
					bertlv.DataApplicationTransactionCounter = hexVal
				case TagUnpredictableNumber:
					bertlv.DataUnpredictableNumber = hexVal
				default:
					if c.customBERTLVMapID4[id] != nil {
						if c.customBERTLVMapID4[id].IsHex {
							bertlv.AddAdditionalData(id, hexVal)
						} else {
							bertlv.AddAdditionalData(id, value)
						}
					}
					// skip unknown
				}
			}
		}
	}
	return bertlv, nil
}

func format(id, value string) string {
	length := utf8.RuneCountInString(value) / 2
	lengthStr := strconv.Itoa(length)
	lengthStr = "00" + fmt.Sprintf("%X", length)
	return id + lengthStr[len(lengthStr)-2:] + value
}

func toHex(s string) string {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return string(dst)
}

func fromHex(h string) (string, error) {
	src := []byte(h)
	dst := make([]byte, utf8.RuneCountInString(h)/2)
	_, err := hex.Decode(dst, src)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}

func formattingTemplate(t BERTLV) string {
	template := ""
	if t.DataApplicationDefinitionFileName != "" {
		template += format(TagApplicationDefinitionFileName, t.DataApplicationDefinitionFileName)
	}
	if t.DataApplicationLabel != "" {
		template += format(TagApplicationLabel, toHex(t.DataApplicationLabel))
	}
	if t.DataTrack2EquivalentData != "" {
		template += format(TagTrack2EquivalentData, t.DataTrack2EquivalentData)
	}
	if t.DataApplicationPAN != "" {
		template += format(TagApplicationPAN, t.DataApplicationPAN)
	}
	if t.DataCardholderName != "" {
		template += format(TagCardholderName, toHex(t.DataCardholderName))
	}
	if t.DataLanguagePreference != "" {
		template += format(TagLanguagePreference, toHex(t.DataLanguagePreference))
	}
	if t.DataIssuerURL != "" {
		template += format(TagIssuerURL, toHex(t.DataIssuerURL))
	}
	if t.DataApplicationVersionNumber != "" {
		template += format(TagApplicationVersionNumber, t.DataApplicationVersionNumber)
	}
	if t.DataIssuerApplicationData != "" {
		template += format(TagIssuerApplicationData, t.DataIssuerApplicationData)
	}
	if t.DataTokenRequestorID != "" {
		template += format(TagTokenRequestorID, t.DataTokenRequestorID)
	}
	if t.DataPaymentAccountReference != "" {
		template += format(TagPaymentAccountReference, t.DataPaymentAccountReference)
	}
	if t.DataLast4DigitsOfPAN != "" {
		template += format(TagLast4DigitsOfPAN, t.DataLast4DigitsOfPAN)
	}
	if t.DataApplicationCryptogram != "" {
		template += format(TagApplicationCryptogram, t.DataApplicationCryptogram)
	}
	if t.DataApplicationTransactionCounter != "" {
		template += format(TagApplicationTransactionCounter, t.DataApplicationTransactionCounter)
	}
	if t.DataUnpredictableNumber != "" {
		template += format(TagUnpredictableNumber, t.DataUnpredictableNumber)
	}
	return template
}
