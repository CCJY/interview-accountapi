Feature: Client Retry Tests
    1. before retry a request, it should get error about "deadline exceed"
    2. after failing twice, it succeeds at the end

    unknown hosts or server status 5xx
    Retry is ignored when context is applied.
    Scenario: before retry a request, it should get error about "deadline exceed"
        Given Context of client has time limt for 100 ms
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
        Then the response should contain error for "deadline exceed"

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