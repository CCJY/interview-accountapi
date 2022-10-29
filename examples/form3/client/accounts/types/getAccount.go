package types

import (
	"github.com/ccjy/interview-accountapi/examples/form3/commons"
	models "github.com/ccjy/interview-accountapi/examples/form3/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type GetAccountResponse = commons.ResponseData[models.AccountData]

type GetAccountRequestContext = client.RequestContext[GetAccountResponse]
type GetAccountResponseContext = client.ResponseContext[GetAccountResponse]
