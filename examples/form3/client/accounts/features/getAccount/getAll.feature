Feature: Get All Account for Form3 API
    1. should get all account

    # this scenario for tests
    Scenario: should create a account
        When I call the method CreateAccount with params
            """
            {
                "data": {
                    "id": "a6c40f81-90d7-4d17-bbef-b9f48fc80acb",
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
        Then the response code should be 201


    Scenario: should get all account
        When I call the method GetAllAccount
        Then the response code should be 200

    # this scenario for tests
    Scenario: should delete a account
        When I call the method DeleteAccount with params "a6c40f81-90d7-4d17-bbef-b9f48fc80acb" "0"
        Then the response code should be 204
        And the response should match json:
            """
            {}
            """


