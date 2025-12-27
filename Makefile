.PHONY: generate-oapi

GENERATED_OAPI_DIR := internal/api
OAPI_PKG := github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

generate-oapi:
	@command -v oapi-codegen >/dev/null 2>&1 || go install $(OAPI_PKG)
	@mkdir -p $(GENERATED_OAPI_DIR)
	@oapi-codegen -generate types      -package api -o $(GENERATED_OAPI_DIR)/types.gen.go  api/openapi.yaml
	@oapi-codegen -generate chi-server -package api -o $(GENERATED_OAPI_DIR)/server.gen.go api/openapi.yaml

.PHONY: up down schema sqlc gen ci-check

up:
	docker compose up -d db

down:
	docker compose down -v

sqlc:
	sqlc generate

ci-check:
	# в CI после gen проверяем, что ничего не изменилось
	git diff --exit-code

