PACKAGE_NAME:='yumster'
BUILT_ON:=$(shell date)
COMMIT_HASH:=$(shell git log -n 1 --pretty=format:"%H")
PACKAGES:=$(shell go list ./... | grep -v /vendor/)
LDFLAGS:='-s -w -X "main.builtOn=$(BUILT_ON)" -X "main.commitHash=$(COMMIT_HASH)"'

all: build docker

test:
	go test -cover -v $(PACKAGES)

update-deps:
	go get -u ./...
	go mod tidy

gofmt:
	go fmt ./...

run: config
	go run -ldflags $(LDFLAGS) `find . | grep -v 'test\|vendor\|repo' | grep \.go`

build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -a -o $(PACKAGE_NAME) .

clean:
	rm -rf yumster* coverage.out coverage-all.out repodata *.rpm *.sqlite

config:
	printf "upload_dir: .\ndev_mode: true" > yumster.yml

docker:
	printf "upload_dir: /repo\n" > yumster.yml
	docker build -t yumster:latest .

# Run just API without NGINX
drun:
	docker run -d -p 8080:8080 --name yumster yumster:latest

compose:
	docker-compose up
