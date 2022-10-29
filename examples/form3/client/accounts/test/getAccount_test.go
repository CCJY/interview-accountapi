package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/cucumber/godog"
)

func (a *AccountClientFeature) iCallTheMethodGetAccountWithParams(arg1 string) error {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)

	got, err := Client.GetAccount(arg1)
	if err != nil {
		a.errMessage = err.Error()
		return err
	}

	rsp, err := json.Marshal(got.ContextData)
	if err != nil {
		return err
	}

	a.statusCode = got.StatusCode()
	a.rsp = rsp

	return nil
}

func (a *AccountClientFeature) iCallTheMethodGetAccountWithContextWithParams(arg1 string) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.timeoutMs)*time.Millisecond)

	defer cancel()
	got, err := Client.GetAccountWithContext(ctx, arg1)

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

func (a *AccountClientFeature) iCallTheMethodGetAllAccount() error {
	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)

	got, err := Client.GetAllAccount()

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

func (a *AccountClientFeature) iCallTheMethodGetAllAccountWithParams(arg1 *godog.DocString) error {
	params := types.GetAllAccountParams{}

	err := json.Unmarshal([]byte(arg1.Content), &params)
	if err != nil {
		return err
	}

	Client := a.getAccountClientTest(a.baseUrl, a.timeoutMs)
	opts := []types.GetAllAccountOpt{}

	if params.Page != nil {
		opts = append(opts, types.WithPage(params.Page.Number, params.Page.Size))
	}

	for _, filter := range params.Filters {
		opts = append(opts, types.WithFilter(string(filter.Key), filter.Value))
	}

	got, err := Client.GetAllAccount(opts...)

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

func InitializeScenarioGetAccount(ctx *godog.ScenarioContext) {
	api := &AccountClientFeature{
		generatedInput: &GeneratedInput{
			Id: Id,
		},
	}

	ctx.Step(`^ID generated$`, api.iDGenerated)
	ctx.Step(`^Timeout (\d+) milliseconds$`, api.timeoutMilliseconds)
	ctx.Step(`^^MockServer has a response delay time for (\d+) milliseconds$$`, api.mockServerHasAResponseDelayTimeForMilliseconds)
	ctx.Step(`^I call the method CreateAccount with params$`, api.iCallTheMethodCreateAccountWithParams)
	ctx.Step(`^I call the method RandomCreateAccount (\d+)$`, api.iCallTheMethodRandomCreateAccount)
	ctx.Step(`^I call the method GetAccount with params "([^"]*)"$`, api.iCallTheMethodGetAccountWithParams)
	ctx.Step(`^I call the method GetAccountWithContext with params "([^"]*)"$`, api.iCallTheMethodGetAccountWithContextWithParams)
	ctx.Step(`^I call the method GetAllAccount$`, api.iCallTheMethodGetAllAccount)
	ctx.Step(`^I call the method GetAllAccount with params$`, api.iCallTheMethodGetAllAccountWithParams)
	ctx.Step(`^I call the method DeleteAccount with params "([^"]*)" "(\d+)"$`, api.iCallTheMethodDeleteAccountWithParams)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
	ctx.Step(`^the response should contain error for "([^"]*)"$`, api.theResponseShouldContainErrorFor)
}

func TestFeatures_GetAccount(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioGetAccount,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/getAccount"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
