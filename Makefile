APP := ghn

build:
	go build -o $(APP) .

run: build
	./$(APP)
