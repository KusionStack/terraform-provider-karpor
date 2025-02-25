.PHONY: build install test testacc doc lint

build:
	@go build -o terraform-provider-karpor

install: build
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/KusionStack/karpor/0.1.0/darwin_amd64
	@mv terraform-provider-karpor ~/.terraform.d/plugins/registry.terraform.io/KusionStack/karpor/0.1.0/darwin_amd64

lint:  
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint@latest ..."; $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && echo -e "Installation complete!\n")
	golangci-lint run ././... --fast --verbose --print-resources-usage

test:
	@go test -v ./...

doc: 
	@cd tools && go generate ./...

testacc:
	@TF_ACC=1 go test -v ./...
