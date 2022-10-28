## Junyoung CHoi

The project started on October 18, 2022.

I started Go for the first time. It was interesting because it was a little different from other languages that I used before. I think I've implemented an essential interface, at least at the production level. I think I'll get better results if I test it a little more and implement it more and refactor it. First of all, I have created unit tests and BDD scenarios. It's my first time at Go, so I don't know the good directory structure yet, but I organized it as follows.


<details><summary>client api</summary>
```go
got, err = Client.NewCreateAccountRequest(&reqData).WithRetry(client.Retry{
            RetryInterval: r.retryWaitMs,
            RetryMax:      r.retryAttempts,
        }).WithContext(ctx).Do()
```
```go
got, err := client.CreateAccount(tt.args.account)
```
</details>



### Directory Structure
```
|-- examples
|   `-- form3
|       |-- client
|       |   `-- accounts
|       |       |-- features
|       |       |   |-- createAccount
|       |       |   |-- deleteAccount
|       |       |   |-- getAccount
|       |       |   `-- retry
|       |       |-- test
|       |       `-- types
|       |-- commons
|       `-- models
|           `-- account
|-- pkg
|   `-- client
|-- scripts
|   `-- db
- - models/accounts
```
### BDD Example
```

    Scenario: after failing twice, it succeeds at the end
        Given MockServer has 150 ms of latency and 50 ms at the end
        And   MockServer returns the 201 response code
        Given RetryAttempt 3 with RetryWait 300 ms per each request
        When I call the method NewCreateAccountRequest with params
            """
            {
                "data": {
                    "id": "e9af97ac-66bc-42da-8e10-5245d8b216df",
                    "organisation_id": "ba61483c-d5c5-4f50-ae81-6b8c039bea43",
                    "type": "accounts",
                    "version": 0,
                    "attributes": {
                        "name": [],
                        "bank_id": "400300",
                        "bank_id_code": "GBDSC",
                        "bic": "NWBKGB22",
                        "country": "GB",
                        "account_matching_opt_out": true
                    }
                }
            }
            """
        And the request is attempted as many times as a given request as
            """
            first response should should be 500
            second request should should be 500
            last request should get the status code 201 and valid data
            """
        Then the response code should be 201
        Then the request was retried 3 times
        Then the response should match json:
            """
            {
                "data": {
                    "id": "e9af97ac-66bc-42da-8e10-5245d8b216df",
                    "organisation_id": "ba61483c-d5c5-4f50-ae81-6b8c039bea43",
                    "type": "accounts",
                    "version": 0,
                    "attributes": {
                        "name": [],
                        "bank_id": "400300",
                        "bank_id_code": "GBDSC",
                        "bic": "NWBKGB22",
                        "country": "GB",
                        "account_matching_opt_out": true
                    }
                }
            }
            """

```

https://github.com/form3tech-oss/interview-accountapi/

https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts


https://spec.openapis.org/oas/v3.1.0#data-types
https://swagger.io/docs/specification/data-models/data-types/


https://github.com/go-playground/validator



Transport important
Should not use default transport becuase caches connections can have unnecessary connections


Production Requirements

Reliable microservices

Possible to different urls

Retryable
on microservice environments, it has to retry a request when the reuqest has diverse issues.
we have to consider unexpected situations because each service has different url and it can have network issues on microserivce environment.


Simple and light

On form3 tech stack
https://stackshare.io/form3/main


https://filia-aleks.medium.com/aws-lambda-battle-2021-performance-comparison-for-all-languages-c1b441005fd1
Warm start
Cold start

application constiner size

Implement client library

Use :
Opt patterns
Singleton patterns
go generic programming
go has no support method generic, its limit


Should I use singleton for http.client
https://stuartleeks.com/posts/connection-re-use-in-golang-with-http-client/
https://github.com/usbarmory/tamago-go/blob/c117e5d62adf00b99dc5cb9e7e0d3105d87fb09d/src/net/http/transport.go#L63-L66

By default, Transport caches connections for future re-use. This may leave many open connections when accessing many hosts. This behavior can be managed using Transport’s CloseIdleConnections method and the MaxIdleConnsPerHost and DisableKeepAlives fields.

if create new transport per one transaction, may be occured error for maximum socket connections


// Examples:

NewRequestContext(
    &RequestContext[account_types.CreateAccountWithResponse]{
        Method:        http.MethodPost,
        BaseUrl:       baseUrl,
        OperationPath: account_types.OperationPathCreateAccount,
        Header: map[string][]string{
            "Content-Type": {"application/json"},
        },
        Body: reqData1,
    }).Do()



Extra implement 
Validation
Generate Code passed tests by generate tools like openapi, swagger

choose oapi-codegen, but it has no suppport validation, thus I made custom client templates and used validator of google-playground



# Form3 Take Home Exercise

Engineers at Form3 build highly available distributed systems in a microservices environment. Our take home test is designed to evaluate real world activities that are involved with this role. We recognise that this may not be as mentally challenging and may take longer to implement than some algorithmic tests that are often seen in interview exercises. Our approach however helps ensure that you will be working with a team of engineers with the necessary practical skills for the role (as well as a diverse range of technical wizardry). 

## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](https://www.api-docs.form3.tech/api/tutorials/getting-started/create-an-account) for information on how to interact with the API. Please note that the fake account API does not require any authorisation or authentication.

A mapping of account attributes can be found in [models.go](./models.go). Can be used as a starting point, usage of the file is not required.

If you encounter any problems running the fake account API we would encourage you to do some debugging first,
before reaching out for help.

## Submission Guidance

### Shoulds

The finished solution **should:**
- Be written in Go.
- Use the `docker-compose.yaml` of this repository.
- Be a client library suitable for use in another software project.
- Implement the `Create`, `Fetch`, and `Delete` operations on the `accounts` resource.
- Be well tested to the level you would expect in a commercial environment. Note that tests are expected to run against the provided fake account API.
- Be simple and concise.
- Have tests that run from `docker-compose up` - our reviewers will run `docker-compose up` to assess if your tests pass.


## How to submit your exercise

- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, by copying all files you deem necessary for your submission
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) [@form3tech-interviewer-1](https://github.com/form3tech-interviewer-1) to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team

## License

Copyright 2019-2022 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.