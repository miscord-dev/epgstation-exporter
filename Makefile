api/epgstation/api.gen.go: api/epgstation/schema.json
	oapi-codegen -package epgstation api/epgstation/schema.json > api/epgstation/api.gen.go

.PHONY: go-generate
go-generate:
	go generate ./...

.PHONY: generate
generate: api/epgstation/api.gen.go go-generate
