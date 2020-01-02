default: build

build: fmtcheck
	go build

clean:
	rm -f azure-keyvault-external-secret

codegen: dependencies
	bash vendor/k8s.io/code-generator/generate-groups.sh "deepcopy,client,informer,lister" github.com/brennerm/azure-keyvault-external-secret/crd/generated github.com/brennerm/azure-keyvault-external-secret/crd azurekeyvaultsecret:v1 --go-header-file ./gen/boilerplate.go.txt

fmt:
	  find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

fmtcheck:
	gofmt_files=$(find . -name '*.go' | grep -v vendor | xargs gofmt -l)
	if [ -n "${gofmt_files}" ]; then exit 1; fi

dependencies:
	go mod vendor; go mod tidy

docker-image:
	DOCKER_BUILDKIT="0" docker build -t brennerm/azure-keyvault-external-secret:latest .

.PHONY: build clean codegen docker-image dependencies fmt fmtcheck
