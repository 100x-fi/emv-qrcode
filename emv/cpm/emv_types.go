package cpm

import (
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
)

// const ...
const (
	IDPayloadFormatIndicator                 = "85" // (M) Payload Format Indicator
	IDApplicationTemplate                    = "61" // (M) Application Template
	IDCommonDataTemplate                     = "62" // (O) Common Data Template
	IDApplicationSpecificTransparentTemplate = "63" // (O) Application Specific Transparent Template
	IDCommonDataTransparentTemplate          = "64" // (O) Common Data Transparent Template
)

// const Data Object ID
const (
	TagApplicationDefinitionFileName = "4F"
	TagApplicationLabel              = "50"
	TagTrack2EquivalentData          = "57"
	TagApplicationPAN                = "5A"
	TagCardholderName                = "5F20"
	TagLanguagePreference            = "5F2D"
	TagIssuerURL                     = "5F50"
	TagApplicationVersionNumber      = "9F08"
	TagIssuerApplicationData         = "9F10"
	TagTokenRequestorID              = "9F19"
	TagPaymentAccountReference       = "9F24"
	TagLast4DigitsOfPAN              = "9F25"
	TagApplicationCryptogram         = "9F26"
	TagApplicationTransactionCounter = "9F36"
	TagUnpredictableNumber           = "9F37"
)

// DataType ...
type DataType string

// const ...
const (
	DataTypeBinary DataType = "binary"
	DataTypeRaw    DataType = "raw"
)

// ID ...
type ID string

// String ...
func (id ID) String() string {
	return string(id)
}

// ParseInt ...
func (id ID) ParseInt() (int64, error) {
	return strconv.ParseInt(id.String(), 10, 64)
}

// Equal ...
func (id ID) Equal(val ID) bool {
	return id == val
}

// Between ...
func (id ID) Between(start ID, end ID) (bool, error) {
	idNum, err := id.ParseInt()
	if err != nil {
		return false, err
	}
	startNum, err := start.ParseInt()
	if err != nil {
		return false, err
	}
	endNum, err := end.ParseInt()
	if err != nil {
		return false, err
	}
	return idNum >= startNum && idNum <= endNum, nil
}

// TLV ...
type TLV struct {
	Tag    ID
	Length string
	Value  string
}

func (tlv TLV) String() string {
	if tlv.Value == "" {
		return ""
	}
	return tlv.Tag.String() + tlv.Length + tlv.Value
}

// DataWithType ...
func (tlv TLV) DataWithType(dataType DataType, indent string) string {
	if tlv.Value == "" {
		return ""
	}
	if dataType == DataTypeBinary {
		rep := regexp.MustCompile("(.{2})")
		hexStr := hex.EncodeToString([]byte(tlv.Value))
		hexArray := rep.FindAllString(hexStr, -1)
		return indent + tlv.Tag.String() + " " + tlv.Length + " " + strings.Join(hexArray, " ") + "\n"
	}
	if dataType == DataTypeRaw {
		return indent + tlv.Tag.String() + " " + tlv.Length + " " + tlv.Value + "\n"
	}
	return ""
}


