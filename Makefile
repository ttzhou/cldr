.PHONY: fmt-check
fmt-check:
	@.tools/golangci-lint fmt -d

.PHONY: fmt
fmt:
	@.tools/golangci-lint fmt

.PHONY: lint
lint:
	@.tools/golangci-lint run -v

.PHONY: lint
tidy:
	@go get -u -t . && go mod tidy

.PHONY: pre-commit
pre-commit: fmt

.PHONY: pre-push
pre-push: pre-commit lint
	# validate codecov config
	@curl -s --fail-with-body --data-binary @.codecov.yml https://codecov.io/validate

.PHONY: lint-fix
lint-fix:
	@.tools/golangci-lint run -v --fix

.PHONY: test
test:
	@go test -v -coverpkg=./num/... ./test/... 

.PHONY: gen-test-cover
gen-test-cover: # not intended for direct use
	@go test -v -coverpkg=./num/... -coverprofile=cover.out ./test/... 

.PHONY: test-cover-report-cli
test-cover-report-cli: gen-test-cover
	@go tool cover -func=cover.out

.PHONY: test-cover-report-browser
test-cover-report-browser: gen-test-cover
	@go tool cover -html=cover.out

.PHONY: clean
clean:
	@rm cover.out

.PHONY: gen-locales
gen-locales:
	@go generate

