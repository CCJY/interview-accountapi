package types

import (
	"github.com/ccjy/interview-accountapi/examples/form3/commons"
	models "github.com/ccjy/interview-accountapi/examples/form3/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type CreateAccountRequest = commons.RequestData[models.AccountData]
type CreateAccountResponse = commons.ResponseData[models.AccountData]

type CreateAccountRequestContext = client.RequestContext[CreateAccountResponse]
type CreateAccountResponseContext = client.ResponseContext[CreateAccountResponse]
