default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

gen-api:
	wget -O schema.yml https://raw.githubusercontent.com/BeryJu/gravity/main/schema.yml
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v6.0.0 generate \
		-i /local/schema.yml \
		-g go \
		-o /local/api \
		-c /local/api/config.yaml
	go get
	go fmt .
	go mod tidy
	rm -f .travis.yml
