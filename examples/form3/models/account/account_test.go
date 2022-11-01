package account

// const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// var seededRand *rand.Rand = rand.New(
// 	rand.NewSource(time.Now().UnixNano()))

// func StringWithCharset(length int, charset string) string {
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[seededRand.Intn(len(charset))]
// 	}
// 	return string(b)
// }

// func String(length int) string {
// 	return StringWithCharset(length, charset)
// }

// func DefaultAccountData() *AccountData {
// 	attributes := AccountAttributes{
// 		BankId:     "400300",
// 		BankIdCode: "GBDSC",
// 		Bic:        "NWBKGB22",
// 		Country:    "GB",
// 		Name: &[]string{
// 			"it", "is", "example",
// 		},
// 	}

// 	return &AccountData{
// 		Attributes:     &attributes,
// 		Id:             uuid.New().String(),
// 		OrganisationId: uuid.New().String(),
// 		Version:        lo.ToPtr(int64(0)),
// 		Type:           lo.ToPtr("accounts"),
// 	}
// }

// func TestValidateDataId(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Id = ""
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Id = uuid.New().String()
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // organization id required and must be uuid
// func TestValidateDataOrganisationId(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.OrganisationId = ""
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.OrganisationId = uuid.New().String()
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // bank_id requied and max length is 11
// func TestValidateDataAttributesBankId(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.BankId = "400300400300400300400300"
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.BankId = "400300"
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // bank_id_code must be GBDSC
// func TestValidateDataAttributesBankIdCode(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.BankIdCode = "EOIJFEFE"
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.BankIdCode = "GBDSC"
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// func TestValidateDataAttributesBIC(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Bic = "EOIJFEFEA"
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Bic = "NWBKGB22"
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Bic = "EOIJFEFEA12"
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// func TestValidateDataAttributesCountry(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Country = "GBE"
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Country = "GB"
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// func TestValidateDataAttributesAccountClassification(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				a := AccountClassification(AccountClassificationBusiness)
// 				accountData.Attributes.AccountClassification = &a
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)
// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// func TestValidateDataAttributesAccountNumber(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.AccountNumber = lo.ToPtr("EOIJFEEFE")
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.AccountNumber = lo.ToPtr("08464524")
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // alternative names [3][140]
// func TestValidateDataAttributesAlternativeNames(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.AlternativeNames = &[]string{"a", "b", "c", "d"}
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.AlternativeNames = &[]string{"add", "bdd", "1aaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaef"}
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.AlternativeNames = &[]string{"add", "bdd", "aa"}
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // currency ISO 4217 must be GBP
// func TestValidateDataAttributesCurrency(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				baseCurrency := BaseCurrency(BaseCurrencyGBP)
// 				accountData.Attributes.BaseCurrency = &baseCurrency
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // iban ISO 13616
// func TestValidateDataAttributesIban(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Iban = lo.ToPtr("aefaefaefaefaefaefaefaefaefaefaefaeaa")
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Iban = lo.ToPtr("GB11NWBK40030041426811")
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// func TestValidateDataAttributesName(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Name = &[]string{"a", "b", "c", "d", "e"}
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Name = &[]string{"add", "bdd", "1aaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaef"}
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Name = &[]string{"add", "bdd", "aa"}
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.Name = &[]string{"add", "bdd", "aa", "a"}
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)
// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// func TestValidateDataAttributesSecondaryIdentification(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.SecondaryIdentification = lo.ToPtr(String(141))
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				accountData.Attributes.SecondaryIdentification = lo.ToPtr(String(1))
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }

// // Status of the account. pending and confirmed are set by Form3, closed can be set manually
// func TestValidateDataAttributesStatus(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    *AccountData
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "No Error",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				e := AccountStatus(StatusClosed)
// 				accountData.Attributes.Status = &e

// 				return accountData
// 			}(),
// 			wantErr: false,
// 		},
// 		{
// 			name: "Error, confirmed is set by Form3",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				e := AccountStatus(StatusConfirmed)
// 				accountData.Attributes.Status = &e

// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 		{
// 			name: "Error, confirmed is set by Form3",
// 			data: func() *AccountData {
// 				accountData := DefaultAccountData()
// 				e := AccountStatus(StatusPending)
// 				accountData.Attributes.Status = &e
// 				return accountData
// 			}(),
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			validate := commons.NewValidator()

// 			err := validate.Struct(tt.data)

// 			if (err != nil) != tt.wantErr {
// 				t.Error(t, err)
// 				return
// 			}
// 		})
// 	}
// }
