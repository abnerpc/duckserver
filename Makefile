
SERVER_NAME=duckserver
DEBUG_SERVER_NAME=duckserver_debug
SERVER_DIR=server

clean-server:
	@rm -f $(SERVER_DIR)/$(SERVER_NAME)*

build-server: clean-server
	@go build ./$(SERVER_DIR) -o $(SERVER_NAME)

run-server: build-server
	@./$(SERVER_DIR)/$(SERVER_NAME)

debug-server: clean-server
	@godebug build ./$(SERVER_DIR) -o $(DEBUG_SERVER_NAME)
	@./$(DEBUG_SERVER_NAME)
