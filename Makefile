
APP_NAME=ds-server
DEBUG_APP_NAME=ds-server-debug

.PHONY: clean
clean:
	@rm -f $(APP_NAME)

.PHONY: build
build: clean
	@go build -o $(APP_NAME)

.PHONY: run
run: build
	@./$(SERVER_DIR)/$(APP_NAME)

.PHONY: debug
debug: clean
	@godebug build -o $(DEBUG_APP_NAME)
	@./$(DEBUG_APP_NAME)
