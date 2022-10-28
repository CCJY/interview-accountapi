package types

import (
	"github.com/ccjy/interview-accountapi/examples/form3/commons"
	models "github.com/ccjy/interview-accountapi/examples/form3/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type GetAllAccountResponse = commons.ResponseDataArray[models.AccountData]

type GetAllAccountRequestContext = client.RequestContext[GetAllAccountResponse]
type GetAllAccountResponseContext = client.ResponseContext[GetAllAccountResponse]
