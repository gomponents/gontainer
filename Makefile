clean:
	rm -rf pkg/template/tmpl_*.go
	rm -rf cmd/license_const.go
	rm -rf app.bin

gogenerate: clean
	go generate ./pkg/... ./cmd/...

tests-unit: gogenerate
	go test -coverprofile=coverage.out ./cmd/... ./pkg/...

lint: gogenerate
	golangci-lint run

code-coverage:
	go tool cover -func=coverage.out

build: export DATETIME = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
build: export GITHASH = $(shell git rev-parse HEAD)
build: export VERSION = dev-$(shell git rev-parse --abbrev-ref HEAD)
build: export DIRTY_SUFFIX = $(shell git diff --quiet || echo '-dirty')
build: clean gogenerate
	go build -v -ldflags="-X 'main.date=${DATETIME}' -X 'main.commit=${GITHASH}${DIRTY_SUFFIX}' -X 'main.version=${VERSION}'" -o app.bin main.go

mv-to-bin-dir:
	mv app.bin /usr/local/bin/gontainer

globally: build mv-to-bin-dir

update-helpers: export HELPERS_PATH = github.com/gomponents/gontainer-helpers
update-helpers:
	go get -u ${HELPERS_PATH}
	cd examples/env && go get -u ${HELPERS_PATH}
	cd examples/global-var && go get -u ${HELPERS_PATH}
	cd examples/library && go get -u ${HELPERS_PATH}
	cd examples/decorators && go get -u ${HELPERS_PATH}

run-example-library:
	cd examples/library && go generate && go run main.go

run-example-env: build
	./app.bin build -i examples/env/gontainer.yml -o examples/env/container.go
	cd examples/env && PERSON_NAME="Harry Potter" PERSON_AGE="13" go run .

run-example-circular-deps: build
	./app.bin build -i examples/circular-deps/gontainer.yml -o /dev/null

run-example-params:
	cd examples/params && gontainer dump-params -i gontainer.yml

run-example-global-var: build
	./app.bin build -i examples/global-var/gontainer.yml -o examples/global-var/container.go
	cd examples/global-var && go run .

run-example-decorators: build
	./app.bin build -i examples/decorators/container/gontainer.yml -o examples/decorators/container/container.go
	cd examples/decorators && go run .

tests: tests-unit lint

.DEFAULT_GOAL := build
