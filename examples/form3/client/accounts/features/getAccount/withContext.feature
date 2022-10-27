Feature: Client Context
    Timeout specifies a time limit for requests made by this Client.

    Scenario: should get error "deadline exceeded" when client with context has time limit
        Given Timeout 100 milliseconds
        Given MockServer has a response delay time for 150 milliseconds
        When I call the method GetAccountWithContext with params "a6c40f81-90d7-4d17-bbef-b9f48fc80acb"
        Then the response should contain error for "deadline exceed"

