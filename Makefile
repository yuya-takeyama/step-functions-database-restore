.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete         cmd/delete/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wait-deleted   cmd/wait-deleted/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/restore        cmd/restore/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wait-available cmd/wait-available/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/modify         cmd/modify/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
