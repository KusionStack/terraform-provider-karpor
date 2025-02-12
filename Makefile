.PHONY: build install test testacc

build:
	@go build -o terraform-provider-karpor

install: build
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/KusionStack/karpor/0.1.0/darwin_amd64
	@mv terraform-provider-karpor ~/.terraform.d/plugins/registry.terraform.io/KusionStack/karpor/0.1.0/darwin_amd64

test:
	@go test -v ./...

testacc:
	@TF_ACC=1 go test -v ./...
