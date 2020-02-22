CONFIG_PATH=config/config.yml.dist
PROGRAM_NAME=campaner

init:
	dep ensure -v
	cp -n $(CONFIG_PATH) config/dev.yml

dependencies:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure

build:
	go build -o $(PROGRAM_NAME)

build_rpi:
	GOOS=linux GOARCH=arm GOARM=7 \
	go build -o $(PROGRAM_NAME)

tests:
	go test ./...