vendor:
	go mod vendor
codegen: vendor
	bash vendor/k8s.io/code-generator/generate-groups.sh "deepcopy,client,informer,lister" crd/generated crd azurekeyvaultsecret:v1


.PHONY: codegen vendor
