package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type AccountClientFeature struct {
	baseUrl        string
	time           int
	errMessage     string
	statusCode     int
	rsp            []byte
	generatedInput *GeneratedInput
}

type GeneratedInput struct {
	Id *uuid.UUID
}

var (
	Id = lo.ToPtr(uuid.New())
)

func (a *AccountClientFeature) executeTemplate(output *bytes.Buffer, content string) error {
	t := template.Must(template.New("").Parse(content))

	err := t.Execute(output, a.generatedInput)
	if err != nil {
		return err
	}

	return nil
}
func (a *AccountClientFeature) getAccountClientTest(baseUrl string, timeout int) AccountClientInterface {
	if baseUrl == "" {
		baseUrl = getHostNmae()
	}
	transport := client.NewTransport()
	ccjyclient := client.NewClient(transport, client.ClientConfig{
		BaseUrl: baseUrl,
		Timeout: timeout,
	}, nil)

	return New(ccjyclient)
}

func (a *AccountClientFeature) iDGenerated() error {
	fmt.Printf("####Given generated Id: %s", a.generatedInput.Id.String())

	return nil
}
func (a *AccountClientFeature) timeoutMilliseconds(arg1 int) error {
	a.time = arg1
	return nil
}

func (a *AccountClientFeature) iCallTheMethodCreateAccountWithParams(arg1 *godog.DocString) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.time)
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

func (a *AccountClientFeature) iCallTheMethodCreateAccountContextWithParams(arg1 *godog.DocString) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.time)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.time)*time.Millisecond)

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

func (a *AccountClientFeature) iCallTheMethodDeleteAccountWithParams(arg1 string, arg2 int) (err error) {
	Client := a.getAccountClientTest(a.baseUrl, a.time)
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

func (a *AccountClientFeature) iCallTheMethodGetAccountWithParams(arg1 string) error {
	Client := a.getAccountClientTest(a.baseUrl, a.time)

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

func (a *AccountClientFeature) theResponseCodeShouldBe(arg1 int) error {
	if a.statusCode != arg1 {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", arg1, a.statusCode)
	}

	return nil
}

func (a *AccountClientFeature) theResponseShouldMatchJson(body *godog.DocString) (err error) {
	var expected, actual interface{}

	var output bytes.Buffer

	err = a.executeTemplate(&output, body.Content)
	if err != nil {
		return nil
	}

	// re-encode expected response
	if err = json.Unmarshal(output.Bytes(), &expected); err != nil {
		return
	}
	// re-encode expected response
	if err = json.Unmarshal(a.rsp, &actual); err != nil {
		return
	}

	// the matching may be adapted per different requirements.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}

	return nil
}

func (a *AccountClientFeature) theResponseShouldContainErrorFor(arg1 string) (err error) {
	expected := arg1

	if !strings.Contains(a.errMessage, arg1) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, err.Error())
	}
	return nil
}
func (a *AccountClientFeature) mockServerHasAResponseDelayTimeForMilliseconds(arg1 int) error {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(arg1) * time.Millisecond)
		}),
	)

	a.baseUrl = s.URL

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &AccountClientFeature{
		generatedInput: &GeneratedInput{
			Id: Id,
		},
	}
	ctx.Step(`^ID generated$`, api.iDGenerated)
	ctx.Step(`^^MockServer has a response delay time for (\d+) milliseconds$$`, api.mockServerHasAResponseDelayTimeForMilliseconds)
	ctx.Step(`^Timeout (\d+) milliseconds$`, api.timeoutMilliseconds)
	ctx.Step(`^I call the method CreateAccountContext with params$`, api.iCallTheMethodCreateAccountContextWithParams)
	ctx.Step(`^I call the method CreateAccount with params$`, api.iCallTheMethodCreateAccountWithParams)
	ctx.Step(`^I call the method DeleteAccount with params "([^"]*)" "(\d+)"$`, api.iCallTheMethodDeleteAccountWithParams)
	ctx.Step(`^I call the method GetAccount with params "([^"]*)"$`, api.iCallTheMethodGetAccountWithParams)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
	ctx.Step(`^the response should contain error for "([^"]*)"$`, api.theResponseShouldContainErrorFor)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
