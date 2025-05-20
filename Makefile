APP := ghn
DESTINATION_PATH ?= /usr/local/bin/$(APP)

build:
	go build -o $(APP) .

run: build
	./$(APP)

scan:
	trufflehog git file://.

install: build
	sudo mv ./$(APP) $(DESTINATION_PATH)
