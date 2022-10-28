# Junyoung CHoi

The project started on October 18, 2022.

I started Go for the first time. It was interesting because it was a little different from other languages that I used before. 


In terms of microservices environments, data may or may not be sent or received for unexpected reasons. It should require implement a retry function on the client side to increase the transmission probability. In addition, in order to reduce unnecessary traffic, it is necessary to check in advance whether data is valid on the client side. Additionally, in a microservices environment, distributed message queues should be used to guarantee the requested data; it can use things like Kafka, Nats and Redis. Also, it can be used common event delivery methods like CloudEvents with Knative to communicate between services and use pub/sub. Container cold start and warm start in Cloud with Kubernetes environment will also have an impact in microservices environment. 



## Form3 Take Home Exercise

<details><summary>Directory Structure</summary>

<p>

> I organized it as follows.

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

</p>
</details>
&nbsp;

## Example Use

```go
got, err := client.CreateAccount(&account)
if err != nil {
    ...
}
fmt.Println(got.ContextData.Data.Id)
```

```go
got, err = Client.NewCreateAccountRequest(&reqData).WithRetry(client.Retry{
            RetryInterval: r.retryWaitMs,
            RetryMax:      r.retryAttempts,
        }).WithContext(ctx).Do()
```

```go
got, err = client.CreateAccountWithContext(ctx, &account)
```


```go
got, err = client.NewCreateAccountRequest(&account).
    	WhenBeforeDo(func(rc *types.CreateAccountRequestContext) error {
            ... // validate
			return nil
		}).WhenAfterDo(func(rc *types.CreateAccountResponseContext) error {
            ...
		return nil
	}).Do()
```


## Tests
The client used fake account API, which is provided as a Docker container in the file `docker-compose.yaml` for operations HTTP Methods `CREATE`, `DELETE`, and `GET`, and Mockserver used it for context, timeout, and retry.
### Requirement Packages
- testing package for TDD
- [godog](https://github.com/cucumber/godog) for BDD


<details><summary>Scenario Example</summary>
<p>


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

</p>
</details>

<br>

# Form3 Take Home Exercise

### Shoulds

The finished solution **should:**
- Be written in Go.
>> I have used the latest version of Golang (1.19)
- Use the `docker-compose.yaml` of this repository.
- Be a client library suitable for use in another software project.
>> The client library is made to be general. Account API related contents are in examples/form3.
- Implement the `Create`, `Fetch`, and `Delete` operations on the `accounts` resource.
- Be well tested to the level you would expect in a commercial environment. Note that tests are expected to run against the provided fake account API.
- Be simple and concise.
- Have tests that run from `docker-compose up` - our reviewers will run `docker-compose up` to assess if your tests pass.
>> I included my app in this file and passed the tests I implemented by docker-compose up.


Transport important
Should not use default transport becuase caches connections can have unnecessary connections


Production Requirements

Reliable microservices

Possible to different urls

application constiner size

Implement client library





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