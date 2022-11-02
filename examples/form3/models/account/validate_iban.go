package account

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	IBAN_GB_REGEX_STR = "^GB\\d{2}[A-Z]{4}\\d{14}$"
	IBAN_BE_REGEX_STR = "^BE\\d{14}$"

	// When Country code: DK (Denmark), FO (Faroe Islands), GL (Greenland)
	IBAN_DK_REGEX_STR = "^DK\\d{16}$|^FO\\d{16}$|^GL\\d{16}$"

	IBAN_EE_REGEX_STR = "^EE\\d{18}$"
	IBAN_FI_REGEX_STR = "^FI\\d{16}$"
	IBAN_FR_REGEX_STR = "^FR\\d{12}[0-9A-Z]{11}\\d{2}$"
	IBAN_DE_REGEX_STR = "^DE\\d{20}$"
	IBAN_GR_REGEX_STR = "^GR\\d{9}[0-9A-Z]{16}$"
	IBAN_IE_REGEX_STR = "^IE\\d{2}[A-Z]{4}\\d{14}$"
	IBAN_IT_REGEX_STR = "^IT\\d{2}[A-Z]\\d{10}[0-9A-Z]{12}$"
	IBAN_LU_REGEX_STR = "^LU\\d{5}[0-9A-Z]{13}$"
	IBAN_NL_REGEX_STR = "^NL\\d{2}[A-Z]{4}\\d{10}$"
	IBAN_PT_REGEX_STR = "^PT\\d{23}$"
	IBAN_PL_REGEX_STR = "^PL\\d{10}[0-9A-Z]{,16}n$"
	IBAN_ES_REGEX_STR = "^ES\\d{22}$"
	IBAN_SE_REGEX_STR = "^SE\\d{22}$"
	IBAN_CH_REGEX_STR = "^CH\\d{7}[0-9A-Z]{12}$"
)

var (
	IBAN_GB_REGEX = regexp.MustCompile(IBAN_GB_REGEX_STR)
	IBAN_BE_REGEX = regexp.MustCompile(IBAN_BE_REGEX_STR)
	IBAN_DK_REGEX = regexp.MustCompile(IBAN_DK_REGEX_STR)
	IBAN_EE_REGEX = regexp.MustCompile(IBAN_EE_REGEX_STR)
	IBAN_FI_REGEX = regexp.MustCompile(IBAN_FI_REGEX_STR)
	IBAN_FR_REGEX = regexp.MustCompile(IBAN_FR_REGEX_STR)
	IBAN_DE_REGEX = regexp.MustCompile(IBAN_DE_REGEX_STR)
	IBAN_GR_REGEX = regexp.MustCompile(IBAN_GR_REGEX_STR)
	IBAN_IE_REGEX = regexp.MustCompile(IBAN_IE_REGEX_STR)
	IBAN_IT_REGEX = regexp.MustCompile(IBAN_IT_REGEX_STR)
	IBAN_LU_REGEX = regexp.MustCompile(IBAN_LU_REGEX_STR)
	IBAN_NL_REGEX = regexp.MustCompile(IBAN_NL_REGEX_STR)
	IBAN_PT_REGEX = regexp.MustCompile(IBAN_PT_REGEX_STR)
	IBAN_PL_REGEX = regexp.MustCompile(IBAN_PL_REGEX_STR)
	IBAN_ES_REGEX = regexp.MustCompile(IBAN_ES_REGEX_STR)
	IBAN_SE_REGEX = regexp.MustCompile(IBAN_SE_REGEX_STR)
	IBAN_CH_REGEX = regexp.MustCompile(IBAN_CH_REGEX_STR)
)

type IbanValidation struct {
	Empty    bool
	RegexStr string
	Regex    *regexp.Regexp
	Error    string
	subtags  []string
}

var (
	IBAN_VALID_MAP = map[string]IbanValidation{
		"GB": {
			RegexStr: IBAN_GB_REGEX_STR,
			Regex:    IBAN_GB_REGEX,
			Error:    "Generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"AU": {
			Error:   "has to be empty",
			subtags: []string{"gt=0,omitempty"},
		},
		"BE": {
			RegexStr: IBAN_BE_REGEX_STR,
			Regex:    IBAN_BE_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"CA": {
			Error:   "not supported, has to be empty",
			subtags: []string{"gt=0,omitempty"},
		},
		"DE": {
			RegexStr: IBAN_DE_REGEX_STR,
			Regex:    IBAN_DE_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"DK": {
			RegexStr: IBAN_DK_REGEX_STR,
			Regex:    IBAN_DK_REGEX,
			Error:    "required, 18 characters",
			subtags:  []string{"required"},
		},
		"FO": {
			RegexStr: IBAN_DK_REGEX_STR,
			Regex:    IBAN_DK_REGEX,
			Error:    "required, 18 characters",
			subtags:  []string{"required"},
		},
		"GL": {
			RegexStr: IBAN_DK_REGEX_STR,
			Regex:    IBAN_DK_REGEX,
			Error:    "required, 18 characters",
			subtags:  []string{"required"},
		},
		"EE": {
			RegexStr: IBAN_EE_REGEX_STR,
			Regex:    IBAN_EE_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"FI": {
			RegexStr: IBAN_FI_REGEX_STR,
			Regex:    IBAN_FI_REGEX,
			Error:    "optional, 18 characters, generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"FR": {
			RegexStr: IBAN_FR_REGEX_STR,
			Regex:    IBAN_FR_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"GR": {
			RegexStr: IBAN_GR_REGEX_STR,
			Regex:    IBAN_GR_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"HK": {
			Error:   "not supported, has to be empty",
			subtags: []string{"gt=0,omitempty"},
		},
		"IE": {
			RegexStr: IBAN_IE_REGEX_STR,
			Regex:    IBAN_IE_REGEX,
			Error:    "not supported, has to be empty",
			subtags:  []string{"omitempty"},
		},
		"IT": {
			RegexStr: IBAN_IT_REGEX_STR,
			Regex:    IBAN_IT_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"LU": {
			RegexStr: IBAN_LU_REGEX_STR,
			Regex:    IBAN_LU_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"NL": {
			RegexStr: IBAN_NL_REGEX_STR,
			Regex:    IBAN_NL_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"PL": {
			RegexStr: IBAN_PL_REGEX_STR,
			Regex:    IBAN_PL_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"PT": {
			RegexStr: IBAN_PT_REGEX_STR,
			Regex:    IBAN_PT_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"ES": {
			RegexStr: IBAN_ES_REGEX_STR,
			Regex:    IBAN_ES_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"SE": {
			RegexStr: IBAN_SE_REGEX_STR,
			Regex:    IBAN_SE_REGEX,
			Error:    "required, 24 characters, never generated",
			subtags:  []string{"required"},
		},
		"CH": {
			RegexStr: IBAN_CH_REGEX_STR,
			Regex:    IBAN_CH_REGEX,
			Error:    "generated if not provided",
			subtags:  []string{"omitempty"},
		},
		"US": {
			Error:   "not supported, has to be empty",
			subtags: []string{"gt=0,omitempty"},
		},
	}
)

// https://www.api-docs.form3.tech/api/tutorials/getting-started/create-an-account
// https://www.iban.com/structure
// http://ht5ifv.serprest.pt/extensions/tools/IBAN/index.html
func IsIBAN(fl validator.FieldLevel) bool {
	account, ok := fl.Parent().Interface().(AccountAttributes)
	if ok {
		switch account.Country {
		case "AU":
			return fl.Field().Interface() == nil
		case "BE":
			return IBAN_BE_REGEX.MatchString(fl.Field().String())
		case "CA":
			return fl.Field().Interface() == nil
		case "EE":
			return IBAN_EE_REGEX.MatchString(fl.Field().String())
		case "FI":
			return IBAN_FI_REGEX.MatchString(fl.Field().String())
		case "FR":
			return IBAN_FR_REGEX.MatchString(fl.Field().String())
		case "DE":
			return IBAN_DE_REGEX.MatchString(fl.Field().String())
		case "GR":
			return IBAN_GR_REGEX.MatchString(fl.Field().String())
		case "HK":
			return fl.Field().Interface() == nil
		case "IE":
			return IBAN_IE_REGEX.MatchString(fl.Field().String())
		case "IT":
			return IBAN_IT_REGEX.MatchString(fl.Field().String())
		case "LU":
			return IBAN_LU_REGEX.MatchString(fl.Field().String())
		case "NL":
			return IBAN_NL_REGEX.MatchString(fl.Field().String())
		case "PT":
			return IBAN_PT_REGEX.MatchString(fl.Field().String())
		case "ES":
			return IBAN_ES_REGEX.MatchString(fl.Field().String())
		case "SE":
			return IBAN_SE_REGEX.MatchString(fl.Field().String())
		case "CH":
			return IBAN_CH_REGEX.MatchString(fl.Field().String())
		case "US":
			return fl.Field().Interface() == nil
		case "GB":
			return IBAN_GB_REGEX.MatchString(fl.Field().String())
		}
	}

	return false
}
