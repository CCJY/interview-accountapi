Feature: Get All Account with Filter for Form3 API
    1. should get all account filtered

    Background:
        Given ID generated

    # this scenario for tests
    Scenario: should create accounts
        Then I call the method RandomCreateAccount 10

    Scenario: should get all account
        When I call the method GetAllAccount with params
            """
            {
                "page": {
                    "number": 0,
                    "size": 1
                },
                "filters": [
                    {
                        "key": "bank_id",
                        "value": "400300"
                    }
                ]
            }
            """
        Then the response code should be 200




