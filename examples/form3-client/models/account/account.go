package account

import (
	"github.com/google/uuid"
)

// AccountData defines model for AccountData.
type AccountData struct {
	Attributes *AccountAttributes `json:"attributes,omitempty"`

	// Unique account id of the account
	Id *uuid.UUID `json:"id,omitempty" validate:"required,uuid4"`

	// Unique organisation id of the account
	OrganisationId *uuid.UUID `json:"organisation_id,omitempty" validate:"required,uuid4"`

	// type of the account
	Type *string `json:"type,omitempty"`

	// version of the account
	Version *int64 `json:"version,omitempty"`
}

// AccountAttributes defines model for AccountAttributes.
type AccountAttributes struct {
	// Classification of account, can be either Personal or Business. Defaults to Personal if not provided. Only used for Confirmation of Payee.
	AccountClassification *AccountClassification `json:"account_classification,omitempty"`

	// Flag to indicate if the account has opted out of account matching, only used for Confirmation of Payee. Defaults to false if not provided.
	//
	// This field has been replaced by name_matching_status.
	AccountMatchingOptOut *bool `json:"account_matching_opt_out,omitempty"`

	// A unique account number will be generated automatically if not provided. If provided, the account number is validated to ensure there's no duplicate. It does not undergo an MOD check.
	AccountNumber *string `json:"account_number,omitempty" validate:"omitempty,bank_account_number"`

	// Alternative primary account names, up to 3 alternative account names with one name in each line of the array. Only used for Confirmation of Payee.
	AlternativeNames *[]string `json:"alternative_names,omitempty" validate:"omitempty,max=3,min=0,dive,max=140"`

	// Local country bank identifier, must be a UK sort code.
	BankId string `json:"bank_id" validate:"required,min=4,max=11"`

	// Identifies the type of bank ID being used, must be GBDSC.
	BankIdCode string `json:"bank_id_code" validate:"required,oneof='GBDSC'"`

	// ISO 4217 code  used to identify the base currency of the account. Must be GBP.
	BaseCurrency *BaseCurrency `json:"base_currency,omitempty"`

	// SWIFT BIC in either 8 or 11 character format.
	Bic string `json:"bic" validate:"required,min=8,max=11,bic"`

	// ISO 3166-1 code  used to identify the domicile of the account. Must be GB
	Country string `json:"country" validate:"required,iso3166_1_alpha2"`

	// IBAN of the account. Generated if not provided.
	Iban *string `json:"iban,omitempty" validate:"omitempty,uk_iban"`

	// Flag to indicate if the account is a joint account, set to true if this is a joint account. Defaults to false if not provided. Only used for Confirmation of Payee.
	JointAccount *bool `json:"joint_account,omitempty"`

	// Name of the account holder, up to four lines possible.
	//
	// For Confirmation of Payee, the following rules apply:
	// * Must be the primary account name.
	// * For concatenated personal names, joint account names and organisation names, use the first line.
	// * If first and last names of a personal name are separated, use the first line for first names, the second line for last names.
	// * Titles are ignored and should not be entered."
	Name *[]string `json:"name,omitempty" validate:"omitempty,max=4,min=0,dive,max=140"`

	// Additional information to identify the account and account holder, 140 characters max. Can be any type of additional identification, e.g. a building society roll number. Only used for Confirmation of Payee.
	SecondaryIdentification *string `json:"secondary_identification,omitempty" validate:"omitempty,max=140"`

	// Status of the account. pending and confirmed are set by Form3, closed can be set manually.
	Status *Status `json:"status,omitempty" validate:"omitempty,oneof='closed'"`

	// Flag to indicate if the account has been switched away from this organisation, only used for Confirmation of Payee.
	//
	// This field has been replaced by name_matching_status.
	Switched *bool `json:"switched,omitempty"`
}

// Defines values for AccountAttributesAccountClassification.
const (
	AccountClassificationBusiness AccountClassification = "Business"
	AccountClassificationPersonal AccountClassification = "Personal"
)

// Defines values for AccountAttributesBaseCurrency.
const (
	BaseCurrencyGBP BaseCurrency = "GBP"
)

// Defines values for AccountAttributesStatus.
const (
	StatusClosed    Status = "closed"
	StatusConfirmed Status = "confirmed"
	StatusPending   Status = "pending"
)

// ISO 4217 code  used to identify the base currency of the account. Must be GBP.
type BaseCurrency string

// Classification of account, can be either Personal or Business. Defaults to Personal if not provided. Only used for Confirmation of Payee.
type AccountClassification string

// Status of the account. pending and confirmed are set by Form3, closed can be set manually.
type Status string
