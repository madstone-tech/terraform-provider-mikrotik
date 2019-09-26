.PHONY: import testacc

TEST ?= ./...

build:
	GO111MODULE=on go build -o terraform-provider-mikrotik

plan: build
	terraform init
	terraform plan

apply:
	terraform apply

test:
	go test $(TEST) -v

testacc:
	TF_ACC=1 go test $(TEST) -v