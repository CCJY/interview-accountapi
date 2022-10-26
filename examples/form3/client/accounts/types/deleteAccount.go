package types

import (
	"github.com/ccjy/interview-accountapi/examples/form3/commons"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type DeleteAccountResponseData struct{}

type DeleteAccountResponse = commons.ResponseData[DeleteAccountResponseData]

type DeleteAccountRequestContext = client.RequestContext[DeleteAccountResponse]
type DeleteAccountResponseContext = client.ResponseContext[DeleteAccountResponse]
