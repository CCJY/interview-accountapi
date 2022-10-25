package account

const (
	ContentType string = "application/json"
)

// DeleteAccountByIdAndVersionParams defines parameters for DeleteAccountByIdAndVersion.
type DeleteAccountByIdAndVersionParams struct {
	// Current version number of the Account resource.
	Version int64 `form:"version" json:"version"`
}

type DeleteAccountByWithResponse struct {
	ErrorMessage string `json:"error_message,omitempty"`
}

// PageFilter defines model for PageFilter.
type FilterPage struct {
	// Page number being requested, defaults to 0.
	Number int

	// Size of the page being requested, defaults to 100.
	Size int
}

type FilterKey string

const (
	FilterAccountNumber FilterKey = "account_number"
	FilterBankId        FilterKey = "bank_id"
	FilterBankIdCode    FilterKey = "bank_id_code"
	FilterCountry       FilterKey = "country"
	FilterCustomerId    FilterKey = "customer_id"
	FilterIban          FilterKey = "iban"
)

type Filter struct {
	Key   FilterKey
	Value string
}
