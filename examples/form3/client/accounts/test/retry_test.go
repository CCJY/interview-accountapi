package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/cucumber/godog"
)

func (r *AccountClientFeature) contextOfClientHasTimeLimtForMs(arg1 int) error {
	r.timeoutMs = arg1

	return nil
}

func (r *AccountClientFeature) iCallTheMethodNewCreateAccountRequestWithParams(arg1 *godog.DocString) error {
	Client := r.getAccountClientTest(r.baseUrl, r.timeoutMs)
	var reqData types.CreateAccountRequest

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutMs)*time.Millisecond)

	defer cancel()
	got, err := Client.NewCreateAccountRequest(&reqData).WithContext(ctx).Do()

	if err != nil {
		r.errMessage = err.Error()
		return nil
	}

	rsp, err := json.Marshal(got.ContextData)
	if err != nil {
		return err
	}

	r.statusCode = got.StatusCode()
	r.rsp = rsp

	return nil
}

func (retry *AccountClientFeature) mockServerHasMsOfLatencyAndMsAtTheEnd(arg1, arg2 int) error {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if retry.retried < retry.retryAttempts {
				time.Sleep(time.Duration(arg1) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(arg2) * time.Millisecond)
			}
			retry.retried += 1
		}),
	)

	retry.baseUrl = s.URL

	return nil
}

func (r *AccountClientFeature) retryAttemptWithRetryWaitSecondsPerEachRequest(arg1, arg2 int) error {
	r.retryAttempts = arg1
	r.retryWaitMs = arg2

	return nil
}

func (r *AccountClientFeature) theRequestIsAttemptedAsManyTimesAsAGivenRequestAs(arg1 *godog.DocString) error {
	return nil
}

func InitializeScenario_Retry(ctx *godog.ScenarioContext) {
	api := &AccountClientFeature{
		generatedInput: &GeneratedInput{
			Id: Id,
		},
	}
	ctx.Step(`^Context of client has time limt for (\d+) ms$`, api.contextOfClientHasTimeLimtForMs)
	ctx.Step(`^MockServer has (\d+) ms of latency and (\d+) ms at the end$`, api.mockServerHasMsOfLatencyAndMsAtTheEnd)
	ctx.Step(`^RetryAttempt (\d+) with RetryWait (\d+) seconds per each request$`, api.retryAttemptWithRetryWaitSecondsPerEachRequest)
	ctx.Step(`^I call the method NewCreateAccountRequest with params$`, api.iCallTheMethodNewCreateAccountRequestWithParams)
	ctx.Step(`^the request is attempted as many times as a given request as$`, api.theRequestIsAttemptedAsManyTimesAsAGivenRequestAs)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
}

func TestFeatures_Retry(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario_Retry,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/retry"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
