Feature: Delete Account
    Create Account

    Scenario: should create a account
        When I call the method CreateAccount with params
            """
            {
                "data": {
                    "id": "6052482c-ad56-4268-a046-23fc0eed6f51",
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

    Scenario: does not allow to delete a account
        When I call the method DeleteAccount with params "6052482c-ad56-4268-a046-23fc0eed6f51" "1"
        Then the response code should be 409
        And the response should match json:
            """
            {
                "error_message": "invalid version"
            }
            """

    Scenario: does not allow to delete a account
        When I call the method DeleteAccount with params "eoifaefoaief" "0"
        Then the response code should be 400
        And the response should match json:
            """
            {
                "error_message": "id is not a valid uuid"
            }
            """

    Scenario: does not allow to delete a account
        When I call the method DeleteAccount with params "" "0"
        Then the response code should be 404
        And the response should match json:
            """
            {}
            """

    Scenario: should delete a account
        When I call the method DeleteAccount with params "6052482c-ad56-4268-a046-23fc0eed6f51" "0"
        Then the response code should be 204
        And the response should match json:
            """
            {}
            """
    Scenario: does not allow to delete an account if not found the account
        When I call the method DeleteAccount with params "6052482c-ad56-4268-a046-23fc0eed6f51" "0"
        Then the response code should be 404
        And the response should match json:
            """
            {}
            """
