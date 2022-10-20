package docs

//go:generate oapi-codegen -old-config-style -package=client -generate types,client -templates=. -o ../../../pkg/client/client.go spec.yaml
