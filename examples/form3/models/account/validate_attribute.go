package account

var (
	MapAttributeValid = map[string]map[string]FieldValidation{
		"GB": {
			"iban": FieldValidation{
				RegexStr: IBAN_GB_REGEX_STR,
				Regex:    IBAN_GB_REGEX,
				Error:    "Generated if not provided",
				subtags:  []string{"omitempty"},
			},
		},
	}
)
