
APP_NAME=duckstatic
DEBUG_APP_NAME=duckstatic-debug

clean:
	@rm -f $(APP_NAME)

build: clean
	@go build -o $(APP_NAME)

run: build
	@./$(SERVER_DIR)/$(APP_NAME)

debug: clean
	@godebug build -o $(DEBUG_APP_NAME)
	@./$(DEBUG_APP_NAME)

.PHONY: clean build run debug
