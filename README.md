# fortyfour-backend

Development mode
`go run cmd/api/main.go`

Build and run
`go build -o bin/server cmd/api/main.go`
`./bin/server`

Run tests with:
`go test -v ./internal/services/...`
`go test -v ./internal/handlers/...`
`go test -v ./internal/utils/...`

Or run all tests:
`go test -v ./...`

With coverage:
`go test -v -coverprofile=coverage.out ./...`
`go tool cover -html=coverage.out`

Run specific test
`go test -v -run TestAuthService_Register_Success ./internal/services/`