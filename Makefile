build: build-bootstrap get-ngrok

build-bootstrap:
	GOARCH=amd64 GOOS=linux go build -o artifacts/bootstrap ./src
.PHONY: build-function

get-ngrok:
	curl -L https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip | tar -xvf - -C artifacts/
.PHONY: build-bootstrap

TEMPLATE_FILE := template.yaml
OUTPUT_FILE := serverless-output.yaml
STACK_NAME := interact-custom-runtime

init:
	aws s3 mb "s3://$(SAM_PACKAGE_BUCKET)"

deploy:
	aws cloudformation package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(SAM_PACKAGE_BUCKET) \
		--output-template-file $(OUTPUT_FILE)
	aws cloudformation deploy \
		--template-file $(OUTPUT_FILE) \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides NgrokAuthToken=$(NGROK_AUTH_TOKEN)

start:
	curl `aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--query "Stacks[0].Outputs[0].OutputValue" \
		 --output text`
