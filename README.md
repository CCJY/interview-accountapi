


https://github.com/form3tech-oss/interview-accountapi/

https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts


https://spec.openapis.org/oas/v3.1.0#data-types
https://swagger.io/docs/specification/data-models/data-types/


https://github.com/go-playground/validator



Transport important
Should not use default transport becuase caches connections can have unnecessary connections


Production Requirements

Reliable microservices

Possible to different urls

Retryable
on microservice environments, it has to retry a request when the reuqest has diverse issues.
we have to consider unexpected situations because each service has different url and it can have network issues on microserivce environment.


Simple and light

On form3 tech stack
https://stackshare.io/form3/main


https://filia-aleks.medium.com/aws-lambda-battle-2021-performance-comparison-for-all-languages-c1b441005fd1
Warm start
Cold start

application constiner size

Implement client library

Use :
Opt patterns
Singleton patterns
go generic programming
go has no support method generic, its limit


Should I use singleton for http.client
https://stuartleeks.com/posts/connection-re-use-in-golang-with-http-client/
https://github.com/usbarmory/tamago-go/blob/c117e5d62adf00b99dc5cb9e7e0d3105d87fb09d/src/net/http/transport.go#L63-L66

By default, Transport caches connections for future re-use. This may leave many open connections when accessing many hosts. This behavior can be managed using Transport’s CloseIdleConnections method and the MaxIdleConnsPerHost and DisableKeepAlives fields.

if create new transport per one transaction, may be occured error for maximum socket connections


// Examples:

NewRequestContext(
    &RequestContext[account_types.CreateAccountWithResponse]{
        Method:        http.MethodPost,
        BaseUrl:       baseUrl,
        OperationPath: account_types.OperationPathCreateAccount,
        Header: map[string][]string{
            "Content-Type": {"application/json"},
        },
        Body: reqData1,
    }).Do()



Extra implement 
Validation
Generate Code passed tests by generate tools like openapi, swagger

choose oapi-codegen, but it has no suppport validation, thus I made custom client templates and used validator of google-playground
