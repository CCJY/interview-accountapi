Feature: Create Account for Form3 API
    Create Account -> Form3 API
    1. should create a account
    2. does not allow to create a existing account resource
    3. does not allow to create an account resource by invalid data

    Background:
        Given ID generated

    Scenario: should create a account
        When I call the method CreateAccount with params
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
        Then the response code should be 201
        And the response should match json:
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

    Scenario: does not allow to create a existing account resource
        When I call the method CreateAccount with params
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
        Then the response code should be 409
        And the response should match json:
            """
            {
                "error_message": "Account cannot be created as it violates a duplicate constraint"
            }
            """

    Scenario: does not allow to create an account resource by invalid data
        When I call the method CreateAccount with params
            """
            {
                "data": {
                    "id": "nil",
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
        Then the response code should be 400


