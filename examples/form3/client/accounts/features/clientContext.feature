Feature: Client Context
    Timeout specifies a time limit for requests made by this Client.

    Background:
        Given ID generated

    Scenario: should get error "deadline exceeded" when client with context has time limit
        Given Timeout 100 milliseconds
        Given MockServer has a response delay time for 150 milliseconds
        When I call the method CreateAccountWithContext with params
            """
            {
                "data": {
                    "id": "{{.Id}}",
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

    Scenario: should get error "deadline exceeded" when client with context has time limit
        Given Timeout 100 milliseconds
        Given MockServer has a response delay time for 150 milliseconds
        When I call the method GetAccountWithContext with params "a6c40f81-90d7-4d17-bbef-b9f48fc80acb"
        Then the response should contain error for "deadline exceed"

    Scenario: should get error "deadline exceeded" when client with context has time limit
        Given Timeout 100 milliseconds
        Given MockServer has a response delay time for 150 milliseconds
        When I call the method DeleteAccountWithContext with params "a6c40f81-90d7-4d17-bbef-b9f48fc80acb" "0"
        Then the response should contain error for "deadline exceed"
