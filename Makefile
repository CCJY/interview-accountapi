
## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	CGO_ENABLED=0 go test -v ./...