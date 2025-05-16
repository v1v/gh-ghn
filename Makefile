APP := ghn

build:
	go build -o $(APP) .

run: build
	./$(APP)

scan:
	trufflehog git file://.
