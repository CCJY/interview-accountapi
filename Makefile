## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	sleep 7
	CGO_ENABLED=0 go test ./...