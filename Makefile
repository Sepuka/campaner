CONFIG_PATH=config/config.yml.dist
PROGRAM_NAME=campaner

init:
	dep ensure -v
	cp -n $(CONFIG_PATH) config/dev.yml

build:
	go build -o $(PROGRAM_NAME)