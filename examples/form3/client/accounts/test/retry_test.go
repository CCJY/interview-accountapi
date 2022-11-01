package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/cucumber/godog"
)

func (r *AccountClientFeature) contextOfClientHasTimeLimtForMs(arg1 int) error {
	r.timeoutMs = arg1

	return nil
}

func (r *AccountClientFeature) iCallTheMethodNewCreateAccountRequestWithParams(arg1 *godog.DocString) (err error) {
	Client := r.getAccountClientTest(r.baseUrl, r.timeoutMs)
	var reqData types.CreateAccountRequest

	err = json.Unmarshal([]byte(arg1.Content), &reqData)
	if err != nil {
		return err
	}

	var got *types.CreateAccountResponseContext

	if 0 < r.timeoutMs {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutMs)*time.Millisecond)
		defer cancel()
		got, err = Client.NewCreateAccountRequest(&reqData).WithRetry(
			client.WithRetryPolicyNoBackOff(r.retryWaitMs, r.retryAttempts),
		).WithContext(ctx).Do()
	} else {
		got, err = Client.NewCreateAccountRequest(&reqData).WithRetry(
			client.WithRetryPolicyNoBackOff(r.retryWaitMs, r.retryAttempts),
		).Do()
	}

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

func (retry *AccountClientFeature) mockServerReturnsTheResponseCode(arg1 int) error {
	retry.mockResponseCode = arg1
	return nil
}

func (retry *AccountClientFeature) mockServerHasMsOfLatencyAndMsAtTheEnd(arg1, arg2 int) error {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			retry.retried += 1
			if retry.retried < retry.retryAttempts {
				// fmt.Printf("retry: %d < retry.retryAttempts: %d", retry.retried, retry.retryAttempts)
				time.Sleep(time.Duration(arg1) * time.Millisecond)
				w.WriteHeader(500)
				return
			}
			// fmt.Printf("last retry: %d retryAttempts: %d", retry.retried, retry.retryAttempts)
			time.Sleep(time.Duration(arg2) * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(retry.mockResponseCode)
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()
			fmt.Print(string(bodyBytes))
			fmt.Fprintln(w, string(bodyBytes))
		}),
	)

	retry.baseUrl = s.URL

	return nil
}

func (r *AccountClientFeature) retryAttemptWithRetryWaitMsPerEachRequest(arg1, arg2 int) error {
	r.retryAttempts = arg1
	r.retryWaitMs = arg2

	return nil
}

func (r *AccountClientFeature) theRequestWasRetriedTimes(arg1 int) error {
	if r.retried != arg1 {
		return fmt.Errorf("expected retried: %d, actual retried: %d", arg1, r.retried)
	}

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
	ctx.Step(`^MockServer returns the (\d+) response code$`, api.mockServerReturnsTheResponseCode)
	ctx.Step(`^RetryAttempt (\d+) with RetryWait (\d+) ms per each request$`, api.retryAttemptWithRetryWaitMsPerEachRequest)
	ctx.Step(`^I call the method DeleteAccount with params "([^"]*)" "(\d+)"$`, api.iCallTheMethodDeleteAccountWithParams)
	ctx.Step(`^I call the method NewCreateAccountRequest with params$`, api.iCallTheMethodNewCreateAccountRequestWithParams)
	ctx.Step(`^the request is attempted as many times as a given request as$`, api.theRequestIsAttemptedAsManyTimesAsAGivenRequestAs)
	ctx.Step(`^the request was retried (\d+) times$`, api.theRequestWasRetriedTimes)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
	ctx.Step(`^the response should contain error for "([^"]*)"$`, api.theResponseShouldContainErrorFor)
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
