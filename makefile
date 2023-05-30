BINARY_NAME=globalssh
 
init:
	go mod init ${BINARY_NAME}


tools:
	go install github.com/mgechev/revive@latest
	go install golang.org/x/tools/gopls@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

deps:
	go mod tidy
	go get github.com/redis/go-redis/v9
	go get github.com/creack/pty
	go get github.com/mattn/go-isatty
	go get github.com/json-iterator/go


build:
	go build -o ${BINARY_NAME} 
	gopls check 
	revive ./... 
	govulncheck ./...


 
test:
	go test -v globalssh/net
	gopls check
	revive

viewcpu:
	go tool pprof cpu.prof

viewmem:
	go tool pprof mem.prof


prof:
	go test ${BINARY_NAME}/server -cpuprofile cpu.prof -memprofile mem.prof
	go tool pprof -raw -output=cpu.txt cpu.prof
	./helpers/stackcollapse-go.pl cpu.txt | ./helpers/flamegraph.pl > cpu.svg
mem:
	go build -gcflags=-m ./...

runc:
	go build -o ${BINARY_NAME}
	./${BINARY_NAME} client
runs:
	go build -o ${BINARY_NAME}
	./${BINARY_NAME} server

release:
	./helpers/build_platforms.sh
	git add .
	git commit . -m "release"
	git push
clean:
	go clean
	rm ${BINARY_NAME}

install: 
	go build -o ${BINARY_NAME}
	echo "this part will only work on unix"
	sudo mv ${BINARY_NAME} /usr/local/bin
	go build -o ${BINARY_NAME}
