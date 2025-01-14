.PHONY: import testacc testclient test

TIMEOUT ?= 40m
ifdef TEST
    TEST := ./... -run $(TEST)
else
    TEST := ./...
endif

ifdef TF_LOG
    TF_LOG := TF_LOG=$(TF_LOG)
endif

build:
	go build -o terraform-provider-mikrotik

clean:
	rm dist/*

plan: build
	terraform init
	terraform plan

apply:
	terraform apply

test: testclient testacc

testclient:
	cd client; go test $(TEST) -race -v -count 1

testacc:
	TF_ACC=1 $(TF_LOG) go test $(TEST) -v -count 1 -timeout $(TIMEOUT)
