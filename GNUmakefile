default: testacc

.PHONY: gen
gen:
	go generate ./... -v

# Run acceptance tests
.PHONY: testacc
testacc: gen
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
