package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func (a *AccountClientFeature) iCallTheMethodDeleteAccountWithParams(arg1 string, arg2 int) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)
	got, err := Client.DeleteAccount(arg1, fmt.Sprint(arg2))
	if err != nil {
		return err
	}

	rsp, err := json.Marshal(got.ContextData)
	if err != nil {
		a.errMessage = err.Error()
		return err
	}

	a.statusCode = got.StatusCode()
	a.rsp = rsp

	return nil
}

func (a *AccountClientFeature) iCallTheMethodDeleteAccountWithContextWithParams(arg1 string, arg2 int) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.timeoutMs)*time.Millisecond)

	defer cancel()
	got, err := Client.DeleteAccountWithContext(ctx, arg1, fmt.Sprint(arg2))

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

func InitializeScenarioDeleteAccount(ctx *godog.ScenarioContext) {
	api := &AccountClientFeature{
		generatedInput: &GeneratedInput{
			Id: Id,
		},
	}

	ctx.Step(`^^MockServer has a response delay time for (\d+) milliseconds$$`, api.mockServerHasAResponseDelayTimeForMilliseconds)
	ctx.Step(`^Timeout (\d+) milliseconds$`, api.timeoutMilliseconds)
	ctx.Step(`^I call the method CreateAccount with params$`, api.iCallTheMethodCreateAccountWithParams)
	ctx.Step(`^I call the method DeleteAccountWithContext with params "([^"]*)" "(\d+)"$`, api.iCallTheMethodDeleteAccountWithContextWithParams)
	ctx.Step(`^I call the method DeleteAccount with params "([^"]*)" "(\d+)"$`, api.iCallTheMethodDeleteAccountWithParams)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
	ctx.Step(`^the response should contain error for "([^"]*)"$`, api.theResponseShouldContainErrorFor)
}

func TestFeatures_DeleteAccount(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioDeleteAccount,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/deleteAccount"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
