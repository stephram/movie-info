.PHONY: build

.ONESHELL:
lint:
	aws cloudformation validate-template --template-body file://template.yaml

build:
	sam build

deploy: build
	sam deploy --stack-name movie-info --s3-bucket aws-sam-cli-managed-default-samclisourcebucket-1cu3rvskuyri4

delete:
	aws cloudformation delete-stack --stack-name movie-info --region ap-southeast-2