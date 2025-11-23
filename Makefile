update-schema:
	oapi-codegen -package api -generate types,server,spec openapi.yaml > internal/api/server.gen.go

update-sql:
	sqlc generate

lint:
	golangci-lint run ./...
