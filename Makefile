.PHONY: build

S3_BUCKET ?= aws-sam-cli-managed-default-samclisourcebucket-1cu3rvskuyri4
STACK_NAME ?= movie-info
TABLE_NAME ?= MoviesTable

.ONESHELL:
setup:
	go get -u github.com/aws/aws-sdk-go/...
	command -v golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.33.0
	command -v mockery || go get github.com/vektra/mockery/v2/.../

build-mocks:
	#mockery --all --dir ./internal
	#mockery --all -srcpkg github.com/aws/aws-sdk-go/aws
	mockery --all --srcpkg github.com/aws/aws-sdk-go/service/dynamodb
	mockery --name Repository --dir ./internal/repository --recursive

lint:
	aws cloudformation validate-template --template-body file://template.yaml
	golangci-lint run ./internal/... ./lambdas/...

test: lint
	go test ./internal/... ./lambdas/... -cover -coverprofile=coverage.out
	@echo "Tests complete"
	go tool cover -html=coverage.out -o ./coverage.html
	@echo "Coverage report written to coverage.html"

build: test
	sam build

package: build
	sam package --s3-bucket $(S3_BUCKET)

deploy: build
	sam deploy --stack-name movie-info --s3-bucket aws-sam-cli-managed-default-samclisourcebucket-1cu3rvskuyri4

delete:
	aws cloudformation delete-stack --stack-name movie-info --region ap-southeast-2

.ONESHELL:
describe-stack:
	aws cloudformation describe-stacks --stack-name $(STACK_NAME)

.ONESHELL:
describe-events:
	aws cloudformation describe-stack-events --stack-name $(STACK_NAME)

.ONESHELL:
scan-table:
	@echo "Scan Table: $(TABLE_NME)"
	aws dynamodb scan --table-name $(TABLE_NAME)