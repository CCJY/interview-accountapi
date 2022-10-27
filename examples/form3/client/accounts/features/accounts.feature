Feature: Form3 account client
    Hello

    Scenario: should create a account
        When I call the method CreateAccount with params
            """
            {
                "data": {
                    "id": "0d209d7f-d07a-4542-947f-5885fddddae5",
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
        And the response should match json:
            """
            {
                "data": {
                    "id": "0d209d7f-d07a-4542-947f-5885fddddae5",
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

    Scenario: does not allow to create a account existing
        When I call the method CreateAccount with params
            """
            {
                "data": {
                    "id": "0d209d7f-d07a-4542-947f-5885fddddae5",
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
        Then the response code should be 409
        And the response should match json:
            """
            {
                "error_message": "Account cannot be created as it violates a duplicate constraint"
            }
            """

    Scenario: should get a account
        When I call the method GetAccount with params "0d209d7f-d07a-4542-947f-5885fddddae5"
        Then the response code should be 200
        And the response should match json:
            """
            {
                "data": {
                    "id": "0d209d7f-d07a-4542-947f-5885fddddae5",
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

    Scenario: should delete a account
        When I call the method DeleteAccount with params "0d209d7f-d07a-4542-947f-5885fddddae5" "0"
        Then the response code should be 204
        And the response should match json:
            """
            {}
            """

    Scenario: does not allow to delete an account if not found the account
        When I call the method DeleteAccount with params "0d209d7f-d07a-4542-947f-5885fddddae5" "0"
        Then the response code should be 404
        And the response should match json:
            """
            {}
            """

# Scenario: should be deadline exceeded by context
#     When I call the method "DELETE" request to "/v1/organisation/accounts/{account_id}?version={version}" with params "0d209d7f-d07a-4542-947f-5885fddddae5" "0"
#     Then the response code should be 204
#     And the response should match json:
#         """
#         """

# Scenario: should be deadline exceeded by timeout
#     When I call the method "DELETE" request to "/v1/organisation/accounts/{account_id}?version={version}" with params "0d209d7f-d07a-4542-947f-5885fddddae5" "0"
#     Then the response code should be 204
#     And the response should match json:
#         """
#         """




