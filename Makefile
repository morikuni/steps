.PHONY: dep
dep:
	go mod donwload

.PHONY: update
update:
	go get -u ./...

.PHONY: test
test:
	go test -race -count 10 ./...

.PHONY: coverage
coverage:
	go test -cover -coverprofile=coverage.out ./...

.PHONY: open-coverage
open-coverage: coverage
	go tool cover -html=coverage.out

