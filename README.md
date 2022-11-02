# Junyoung Choi

## I am New to GO :smile:
The project started on October 18, 2022.

I started Go for the first time. It was interesting because it was a little different from other languages that I used before. 

For more details, I have used the latest version of Golang (1.19) and tried to make it general for other projects. The sources in the pkg/client are general library, and the Form3 client in the examples is the result of using them. When I made this project, as I am starting Go for the first time, I decided to try generic programming supported from 1.18. I tried to make the request and response data generic. In addition, I included my app in `docker-compose.yml` and passed the tests I implemented by docker-compose up.


# Result

If you're just starting out, the form3-fake-accountapi is a good resource. It definitely seems to help when studying Golang. This project took about 2 weeks. It was a time to study Golang and get to know various patterns used in Golang. It's my first time using Golang, so I've been trying to do a lot of things with my own greed. Actually, my future plan is to template the client library using something like oapi-codegen, create a logging interface, and implement how to do tracing in the cloud environment, but unfortunately, I don't know when I will do it. When I went to the study branch, I made an example template using resty client and oapi-codegen. oapi-codegen has its own template, but I have customized the template. And I also made a very simple Client library in the study branch. Hope this helps someone.

Unfortunately, I couldn't go to the next step.
After receiving the feedback, I summarized the reasons I thought.
1. The test cases are unorganized and difficult to read. And the lack of explanation seems to have been a big factor. I thought BDD would be enough to explain, but that was just my opinion because I thought this was documentation enough. I guess I didn't provide a sufficient explanation. And above all, it seems that the test cases were not sophisticated.
2. Implementing a generic client library seems complicated. It seems to have been negatively impacted by undesired features. In particular, the retry feature didn't seem appropriate for them to test. They want it to be simple. However, it is difficult to create a simple test case based on a production environment. I think the submission guidelines are a bit confusing. I should have contacted them a lot and asked more about the author's intentions.

In fact, I wanted to receive more feedback, so I opened it. How can I communicate more clearly with my team members using TDD or BDD? And how can these implementations be trusted?


# Summary

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


<details><summary>Example Use</summary>

<p>

`DELETE` and `GET` account can also be used in the same way.

- If you don't need an option, you can use it.
```go
got, err := client.CreateAccount(&account)
if err != nil {
    ...
}
fmt.Println(got.ContextData.Data.Id)
```

- This can be used if a retry is required. Additionally, you can set a retry when you create the Client.
```go
got, err = client.NewCreateAccountRequest(&reqData).
            WithRetry(
		        client.WithRetryPolicyExpoFullyBackOff(100, 300, 3),
            ).Do()
```

- This is used when using context.
```go
got, err = client.CreateAccountWithContext(ctx, &account)
```

- This is used when you want to manipulate data.
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

- This is used when you want to get accounts filtered.
```go
got, err = client.GetAllAccount(
            types.WithPage(0, 1),
            types.WithFilter("country", "GB"))
```

</p>
</details>

&nbsp;

## Tests
The client used fake account API, which is provided as a Docker container in the file `docker-compose.yaml` for operations HTTP Methods `CREATE`, `DELETE`, and `GET`, and Mockserver used it for context, timeout, and retry.
### Used Packages
- testing package for TDD
- [godog](https://github.com/cucumber/godog) for BDD
- [lo](https://github.com/samber/lo) 


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

### Result
![TEST](./test_report_img.png)
<br>

# Form3 Take Home Exercise

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
