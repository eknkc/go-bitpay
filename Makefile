build:
	rm -rf build/
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s" -o build/gobitpay-darwin-x64
	GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o build/gobitpay-linux-x64
	tar cvzf build/gobitpay-darwin-x64.tar.gz build/gobitpay-darwin-x64
	tar cvzf build/gobitpay-linux-x64.tar.gz build/gobitpay-linux-x64
	rm build/gobitpay-darwin-x64 build/gobitpay-linux-x64

.PHONY: build
