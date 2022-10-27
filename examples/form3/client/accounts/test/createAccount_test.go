package test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/cucumber/godog"
)

func (a *AccountClientFeature) iCallTheMethodCreateAccountWithParams(arg1 *godog.DocString) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)
	var output bytes.Buffer

	err = a.executeTemplate(&output, arg1.Content)
	if err != nil {
		return nil
	}

	var reqData types.CreateAccountRequest

	err = json.Unmarshal(output.Bytes(), &reqData)
	if err != nil {
		return err
	}

	got, err := Client.CreateAccount(&reqData)

	if err != nil {
		a.errMessage = err.Error()
		return err
	}

	a.statusCode = got.StatusCode()

	rsp, err := json.Marshal(got.ContextData)
	if err != nil {
		return err
	}

	a.rsp = rsp

	return nil
}

func (a *AccountClientFeature) iCallTheMethodCreateAccountWithContextWithParams(arg1 *godog.DocString) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)
	var output bytes.Buffer

	err = a.executeTemplate(&output, arg1.Content)
	if err != nil {
		return nil
	}

	var reqData types.CreateAccountRequest

	err = json.Unmarshal(output.Bytes(), &reqData)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.timeoutMs)*time.Millisecond)

	defer cancel()
	got, err := Client.CreateAccountWithContext(ctx, &reqData)

	if err != nil {
		a.errMessage = err.Error()
		return nil
	}

	rsp, err := json.Marshal(got.ContextData)
	if err != nil {
		return err
	}

	a.statusCode = got.StatusCode()
	a.rsp = rsp

	return nil
}

func InitializeScenarioCreateAccount(ctx *godog.ScenarioContext) {
	api := &AccountClientFeature{
		generatedInput: &GeneratedInput{
			Id: Id,
		},
	}
	ctx.Step(`^^MockServer has a response delay time for (\d+) milliseconds$$`, api.mockServerHasAResponseDelayTimeForMilliseconds)
	ctx.Step(`^Timeout (\d+) milliseconds$`, api.timeoutMilliseconds)
	ctx.Step(`^ID generated$`, api.iDGenerated)
	ctx.Step(`^I call the method CreateAccount with params$`, api.iCallTheMethodCreateAccountWithParams)
	ctx.Step(`^I call the method CreateAccountWithContext with params$`, api.iCallTheMethodCreateAccountWithContextWithParams)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
	ctx.Step(`^the response should contain error for "([^"]*)"$`, api.theResponseShouldContainErrorFor)
}

func TestFeatures_CreateAccount(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioCreateAccount,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/createAccount"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
