# Define variables for Go and Go-related tools
GO = go

# Run the Go application
serve:
	$(GO) run cmd/server/main.go

migrate :
	$(GO) run cmd/migration/main.go