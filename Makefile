# Include variables from the .envrc file
include .envrc

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/auth application
.PHONY: run/auth
run/auth:
	go run ./cmd/auth -db-connection-string=${MONGODB_URI}
