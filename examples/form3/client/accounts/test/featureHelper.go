package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts"
	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type AccountClientFeature struct {
	baseUrl          string
	timeoutMs        int
	retried          int
	retryAttempts    int
	retryWaitMs      int
	mockResponseCode int
	errMessage       string
	statusCode       int
	rsp              []byte
	generatedInput   *GeneratedInput
}

type GeneratedInput struct {
	Id *uuid.UUID
}

var (
	Id = lo.ToPtr(uuid.New())
)

func getHostNmae() string {
	env := os.Getenv("APP-ENV")
	switch env {
	case "docker":
		return "http://host.docker.internal:8080"
	default:
		return "http://127.0.0.1:8080"
	}
}

func (a *AccountClientFeature) executeTemplate(output *bytes.Buffer, content string) error {
	t := template.Must(template.New("").Parse(content))

	err := t.Execute(output, a.generatedInput)
	if err != nil {
		return err
	}

	return nil
}
func (a *AccountClientFeature) getAccountClientTest(baseUrl string, timeout int) accounts.AccountClientInterface {
	if baseUrl == "" {
		baseUrl = getHostNmae()
	}
	transport := client.NewTransport()
	ccjyclient := client.NewClient(transport, client.ClientConfig{
		BaseUrl: baseUrl,
		Timeout: timeout,
	}, nil)

	return accounts.New(ccjyclient)
}

func (a *AccountClientFeature) iDGenerated() error {
	fmt.Printf("####Given generated Id: %s", a.generatedInput.Id.String())

	return nil
}

func (a *AccountClientFeature) timeoutMilliseconds(arg1 int) error {
	a.timeoutMs = arg1
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
