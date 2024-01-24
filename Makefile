test:
	go test ./...

mock:
	mockgen -source=internal/services/account.go -destination=mocks/services/account.go
	mockgen -source=internal/services/transaction.go -destination=mocks/services/transaction.go
	mockgen -source=internal/services/user.go -destination=mocks/services/user.go
	mockgen -source=internal/usecases/account.go -destination=mocks/usecases/account.go
	mockgen -source=internal/usecases/transaction.go -destination=mocks/usecases/transaction.go
	mockgen -source=internal/usecases/user.go -destination=mocks/usecases/user.go

dep:
	go mod tidy
	go mod vendor

init-db:
	migrate create -ext sql -dir migrations -seq mock_data