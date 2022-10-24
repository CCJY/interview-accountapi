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
